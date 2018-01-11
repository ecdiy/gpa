package gpa

import (
	"database/sql"
)

func (dao *Impl) QueryArrayInt(rows *sql.Rows, cols []string) ([]int, error) {
	colLen := len(cols)
	if colLen == 1 {
		var list []int
		for rows.Next() {
			var r int
			rows.Scan(&r)
			list = append(list, r)
		}
		return list, nil
	} else {
		for rows.Next() {
			arr := make([]interface{}, colLen)
			res := make([]int, colLen)
			for i := 0; i < colLen; i++ {
				arr[i] = &res[i]
			}
			return res, nil
		}
	}
	return []int{}, nil
}

func (impl *Impl) QueryArrayString(rows *sql.Rows, cols []string) ([]string, error) {
	colLen := len(cols)
	if colLen == 1 {
		var list []string
		for rows.Next() {
			var r string
			rows.Scan(&r)
			list = append(list, r)
		}
		return list, nil
	} else {
		if rows.Next() {
			arr := make([]interface{}, colLen)
			res := make([]string, colLen)
			for i := 0; i < colLen; i++ {
				arr[i] = &res[i]
			}
			rows.Scan(arr...)
			return res, nil
		}
	}
	return nil, nil
}
