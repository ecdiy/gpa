package gpa

import (
	"database/sql"
	"reflect"
)

func (*Impl) QueryArrayMapStringInterface(rows *sql.Rows, cols []string) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		arr := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			var inf string
			arr[i] = &inf
		}
		rows.Scan(arr...)
		res := make(map[string]interface{})
		for i := 0; i < len(cols); i++ {
			res[cols[i]] = arr[i]
		}
		result = append(result, res)
	}
	return result, nil
}

func (*Impl) QueryArrayMapStringString(rows *sql.Rows, cols []string) ([]map[string]string, error) {
	result := make([]map[string]string, 0)
	for rows.Next() {
		arr := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			var inf string
			arr[i] = &inf
		}
		rows.Scan(arr...)
		res := make(map[string]string)
		for i := 0; i < len(cols); i++ {
			res[cols[i]] = reflect.ValueOf(arr[i]).Elem().Interface().(string)
		}
		result = append(result, res)
	}
	return result, nil
}
