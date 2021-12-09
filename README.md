# dblite

SQLite extension for accessing other SQL databases, in SQLite.
Similar to how [Postgres Foreign Data Wrappers](https://wiki.postgresql.org/wiki/Foreign_data_wrappers) enable access to other databases in PostgreSQL, this project is a way to access "foreign" SQL databases in SQLite.
It compiles to a [run-time loadable extension](https://www.sqlite.org/loadext.html), using [`riyaz-ali/sqlite`](https://github.com/riyaz-ali/sqlite).

It uses `database/sql`, so should be compatible with any database that [implements a driver](https://github.com/golang/go/wiki/SQLDrivers).

## (Currently) Supported Databases

- PostgreSQL (`postgres`)
- MySQL (`mysql`)

## Available Functions

### `dblite_open`

Scalar function that registers a database by name for use in later queries.

Params:
1. `TEXT` driver name (see supported databases above)
2. `TEXT` give a name to the database connection, referenced in subsequent queries
3. `TEXT` database connection string

```sql
SELECT dblite_open('postgres', 'mydb', 'postgres://...')
```

### `dblite_ping`

Scalar function that checks if an opened database can be connected to.

Params:
1. `TEXT` database name (string used as second param to `db_open`)

```sql
SELECT dblite_ping('mydb')
```

### `dblite_exec`

Scalar function that executes a SQL statement without any results.

Params:
1. `TEXT` database name
2. `TEXT` SQL to execute

```sql
SELECT dblite_exec('mydb', 'DROP TABLE ...')
```

### `dblite_query`

Table-valued function that executes a SQL statement and returns the results.

Params:
1. `TEXT` database name
2. `TEXT` SQL to execute

```sql
SELECT * FROM dblite_query('mydb', 'SELECT * FROM ...')
```

### `dblite_close`

Scalar function that closes a database and deregisters it.

Params:
1. `TEXT` database name

```sql
SELECT dblite_close('mydb')
```
