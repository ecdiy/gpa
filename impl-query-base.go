package gpa

import (
	"database/sql"
)

func ( *Impl) QueryInt64Bool(rows *sql.Rows, cols []string) (int64, bool, error) {
	if rows.Next() {
		var r int64
		rows.Scan(&r)
		return r, true, nil
	}
	return 0, false, nil
}

func ( *Impl) QueryIntBool(rows *sql.Rows, cols []string) (int, bool, error) {
	if rows.Next() {
		var r int
		rows.Scan(&r)
		return r, true, nil
	}
	return 0, false, nil
}

func ( *Impl) QueryStringBool(rows *sql.Rows, cols []string) (string, bool, error) {
	if rows.Next() {
		var r string
		rows.Scan(&r)
		return r, true, nil
	}
	return "", false, nil
}

//---------------------------

