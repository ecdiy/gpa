package gpa

import (
	"database/sql"
	"reflect"
)

func ( *Impl) QueryArrayArrayString(rows *sql.Rows, cols []string) ([][]string, error) {
	var result [][]string
	colLen := len(cols)
	for rows.Next() {
		arr := make([]interface{}, colLen)
		for i := 0; i < colLen; i++ {
			var inf string
			arr[i] = &inf
		}
		rows.Scan(arr...)
		res := make([]string, colLen)
		for i := 0; i < colLen; i++ {
			res[i] = reflect.ValueOf(arr[i]).Elem().Interface().(string)
		}
		result = append(result, res)
	}
	return result, nil
}
