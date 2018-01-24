package gpa

import (
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

	Get(key string) (string, bool, error)
	Set(key string, val interface{}) (int64, error)
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
