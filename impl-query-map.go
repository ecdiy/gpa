package gpa

import (
	"database/sql"
	"reflect"
)

func (impl *Impl) QueryMapStringInterfaceBool(rows sql.Rows, cols []string) (map[string]interface{}, bool, error) {
	if rows.Next() {
		arr := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			var inf interface{}
			arr[i] = &inf
		}
		rows.Scan(arr...)
		res := make(map[string]interface{})
		for i := 0; i < len(cols); i++ {
			res[cols[i]] = reflect.ValueOf(arr[i]).Elem().Interface()
		}
		return res, true, nil
	}
	return nil, false, nil
}

func (impl *Impl) QueryMapStringStringBool(rows sql.Rows, cols []string) (map[string]string, bool, error) {
	if rows.Next() {
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
		return res, true, nil
	}
	return nil, false, nil
}
