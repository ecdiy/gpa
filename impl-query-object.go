package gpa

import (
	"database/sql"
	"reflect"
	"strings"
)

func (dao *Impl) QueryObjectBool(rows *sql.Rows, cols []string, resultType reflect.Type) []reflect.Value {
	v := reflect.New(resultType).Elem()
	numF := resultType.NumField()
	if rows.Next() {
		oneRow := make([]interface{}, len(cols))
		for j, c := range cols {
			cLow := strings.ToLower(c)
			for i := 0; i < numF; i++ {
				tf := resultType.Field(i)
				if cLow == strings.ToLower(strings.ToLower(tf.Name)) {
					oneRow[j] = v.FieldByName(tf.Name).Addr().Interface()
					break
				}
			}
		}
		err := rows.Scan(oneRow...)
		if err != nil {
			return []reflect.Value{nilVf, reflect.ValueOf(false), reflect.ValueOf(err)}
		}
	}
	return []reflect.Value{v, reflect.ValueOf(true), nilVf}
}
