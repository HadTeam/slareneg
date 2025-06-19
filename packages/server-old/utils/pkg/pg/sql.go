package pg

import (
	"database/sql"
	"strings"
)

func SqlQueryExist(sql string, arg ...any) bool {
	r, err := db.Query(strings.Replace(sql, "SELECT *", "SELECT count(*)", 1), arg...)
	if err != nil {
		logger.Warn("Statements error (existing query): ", err)
		logger.Trace("sql: ", sql, "arg: ", arg)
		return false
	}
	var count int
	r.Next()
	err = r.Scan(&count)
	if err != nil {
		count = 0
	}
	_ = r.Close()
	return count != 0
}

func SqlQuery(sql string, arg ...any) *sql.Rows { // WARN: Remember to close the result set
	r, err := db.Query(sql, arg...)
	if err != nil {
		logger.Warn("Statements error (query): ", err)
		logger.Trace("sql: ", sql, "arg: ", arg)
	}
	return r
}

func SqlExec(sql string, arg ...any) sql.Result {
	r, err := db.Exec(sql, arg...)
	if err != nil {
		logger.Warn("Statements error (exec): ", err, ' ', sql, ' ', arg)
		logger.Trace("sql: ", sql, "arg: ", arg)
	}
	return r
}
