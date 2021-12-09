package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)

var queryCols = []vtab.Column{
	{Name: "column_names", Type: "JSON", OrderBy: vtab.NONE},
	{Name: "column_types", Type: "JSON", OrderBy: vtab.NONE},
	{Name: "contents", Type: "JSON", OrderBy: vtab.NONE},

	{Name: "database_name", Type: "TEXT", Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, OmitCheck: true}}, OrderBy: vtab.NONE},
	{Name: "query", Type: "TEXT", Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, OmitCheck: true}}, OrderBy: vtab.NONE},
}

func newQueryTableFunc(databases dbmap) sqlite.Module {
	return vtab.NewTableFunc("dblite_query", queryCols, func(constraints []*vtab.Constraint, order []*sqlite.OrderBy) (vtab.Iterator, error) {
		var name, query string

		for _, constraint := range constraints {
			if constraint.Op == sqlite.INDEX_CONSTRAINT_EQ {
				switch queryCols[constraint.ColIndex].Name {
				case "database_name":
					name = constraint.Value.Text()
				case "query":
					query = constraint.Value.Text()
				}
			}
		}

		if name == "" {
			return nil, errors.New("must supply a database name")
		}

		if query == "" {
			return nil, errors.New("must supply a query")
		}

		if db, ok := databases[name]; !ok {
			return nil, fmt.Errorf("database with name %s has not been opened", name)
		} else {
			if rows, err := db.QueryContext(context.TODO(), query); err != nil {
				return nil, fmt.Errorf("could not execute query: %v", err)
			} else {
				return &queryResultsIter{rows: rows}, nil
			}
		}
	})
}

type queryResultsIter struct {
	rows *sql.Rows
}

func (i *queryResultsIter) Column(ctx vtab.Context, c int) error {
	columnTypes, err := i.rows.ColumnTypes()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columnTypes))
	for i := range values {
		values[i] = new(interface{})
	}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	err = i.rows.Scan(values...)
	if err != nil {
		return err
	}

	columnNames := make([]string, len(columnTypes))
	columnTypeNames := make([]string, len(columnTypes))

	dest := make(map[string]interface{})
	for i, column := range columnTypes {
		columnNames[i] = column.Name()
		columnTypeNames[i] = column.DatabaseTypeName()
		dest[column.Name()] = *(values[i].(*interface{}))
	}

	if err := enc.Encode(dest); err != nil {
		return err
	}

	switch queryCols[c].Name {
	case "column_names":
		if b, err := json.Marshal(columnNames); err != nil {
			return err
		} else {
			ctx.ResultText(string(b))
		}
	case "column_types":
		if b, err := json.Marshal(columnTypeNames); err != nil {
			return err
		} else {
			ctx.ResultText(string(b))
		}
	case "contents":
		ctx.ResultText(b.String())
	}
	return nil
}

func (i *queryResultsIter) Next() (vtab.Row, error) {
	if i.rows.Next() {
		return i, nil
	} else {
		if err := i.rows.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
}
