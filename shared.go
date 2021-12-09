package main

import (
	"go.riyazali.net/sqlite"
)

var (
	databases = make(dbmap)
)

func init() {
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateFunction("dblite_open", newOpenFunc(databases)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		if err := api.CreateFunction("dblite_close", newCloseFunc(databases)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		if err := api.CreateFunction("dblite_ping", newPingFunc(databases)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		if err := api.CreateFunction("dblite_exec", newExecFunc(databases)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		if err := api.CreateModule("dblite_query", newQueryTableFunc(databases), sqlite.EponymousOnly(true)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		return sqlite.SQLITE_OK, nil
	})
}

func main() {}
