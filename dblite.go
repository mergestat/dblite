package main

import (
	"database/sql"
	"fmt"

	"go.riyazali.net/sqlite"
)

type dbmap map[string]*sql.DB

type open struct {
	databases dbmap
}

func newOpenFunc(databases dbmap) sqlite.Function {
	return &open{databases: databases}
}

func (o *open) Args() int           { return 3 }
func (o *open) Deterministic() bool { return false }
func (o *open) Apply(context *sqlite.Context, values ...sqlite.Value) {
	driver := values[0].Text()
	name := values[1].Text()
	connectionString := values[2].Text()

	db, err := sql.Open(driver, connectionString)
	if err != nil {
		context.ResultError(err)
		return
	}

	if _, ok := o.databases[name]; ok {
		context.ResultError(fmt.Errorf("database with name: %s already exists, please call dblite_close('%s') before opening a DB with the same name", name, name))
		return
	}

	databases[name] = db
	context.ResultInt(1)
}

type close struct {
	databases dbmap
}

func newCloseFunc(databases dbmap) sqlite.Function {
	return &close{databases: databases}
}

func (c *close) Args() int           { return 1 }
func (c *close) Deterministic() bool { return false }
func (c *close) Apply(context *sqlite.Context, values ...sqlite.Value) {
	name := values[0].Text()

	if db, ok := c.databases[name]; ok {
		if err := db.Close(); err != nil {
			context.ResultError(err)
			return
		}
		delete(c.databases, name)
		context.ResultInt(1)
	} else {
		context.ResultInt(0)
		return
	}
}

type ping struct {
	databases dbmap
}

func newPingFunc(databases dbmap) sqlite.Function {
	return &ping{databases: databases}
}

func (p *ping) Args() int           { return 1 }
func (p *ping) Deterministic() bool { return false }
func (p *ping) Apply(context *sqlite.Context, values ...sqlite.Value) {
	name := values[0].Text()

	if db, ok := p.databases[name]; ok {
		if err := db.Ping(); err != nil {
			context.ResultError(err)
			return
		} else {
			context.ResultInt(1)
		}
	} else {
		context.ResultError(fmt.Errorf("unknown database with name: %s", name))
		return
	}
}

type exec struct {
	databases dbmap
}

func newExecFunc(databases dbmap) sqlite.Function {
	return &exec{databases: databases}
}

func (e *exec) Args() int           { return 2 }
func (e *exec) Deterministic() bool { return false }
func (e *exec) Apply(context *sqlite.Context, values ...sqlite.Value) {
	name := values[0].Text()
	query := values[1].Text()

	if db, ok := e.databases[name]; ok {
		if _, err := db.Exec(query); err != nil {
			context.ResultError(err)
			return
		} else {
			context.ResultInt(1)
		}
	} else {
		context.ResultError(fmt.Errorf("unknown database with name: %s", name))
		return
	}
}
