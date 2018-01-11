package gpa

import (
	"strings"
	"reflect"
	"database/sql"
)

const (
	PrimaryId     = "@Id"
	AutoIncrement = "AutoIncrement"
)

var nilVf = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())

type Gpa interface {
	Save(model interface{}) (int64, error)
	Insert(s string, param ... interface{}) (int64, error)
	Exec(s string, param ... interface{}) (int64, error)
}

type Impl struct {
	driver, dsn string
	conn        *sql.DB
}

func vti(in []reflect.Value) []interface{} {
	p := make([]interface{}, len(in))
	for idx, pin := range in {
		p[idx] = pin.Interface()
	}
	return p
}

func FmtGroupConcatArray(list []map[string]interface{}, colName string) []map[string]interface{} {
	for _, v := range list {
		sub, b := v[colName].(string)
		if b {
			subs := strings.Split(sub, ",")
			subA := make([][]string, len(subs))
			for i, sb := range subs {
				subA[i] = strings.Split(sb, "#")
			}
			v[colName] = subA
		}
	}
	return list
}

func FmtGroupConcat(list []map[string]interface{}, colName string) []map[string]interface{} {
	for _, v := range list {
		sub, b := v[colName].(string)
		if b {
			v[colName] = strings.Split(sub, ",")
		}
	}
	return list
}
