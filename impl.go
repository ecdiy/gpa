package gpa

import (
	"database/sql"
	"reflect"
	"strings"
	"github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

func GetImpl(driverName, dbUri string, models ... interface{}) *Impl {
	impl := &Impl{driver: driverName, dsn: dbUri}
	var err error
	impl.conn, err = sql.Open(impl.driver, impl.dsn)
	if err != nil {
		panic("数据库连接错误")
	} else {
		impl.conn.SetMaxOpenConns(5)
		//	dao.conn.SetMaxIdleConns(0)
		impl.conn.Ping()
	}
	for _, d := range models {
		impl.setMethodImpl(d)
	}
	return impl
}

func GetGpa(dbName, dbUri string, models ... interface{}) Gpa {
	return GetImpl(dbName, dbUri, models ...)
}

func getSqlByMethod(ft reflect.StructField) string {
	name := ft.Name
	if strings.Index(name, "FindBy") == 0 {
		ty := ft.Type.String()

		lk := strings.LastIndex(ty, "(")
		if lk > 0 {
			ty = ty[lk+1:]
		}
		d := strings.Index(ty, ".")
		x := strings.Index(ty, ",")
		if d > 0 && x > d {
			tb := ty[d+1:x ]
			rep := strings.Replace(name[6:], "And", "=? And ", -1)
			return "select * from " + tb + " where " + rep + "=?"
		} else {
			seelog.Error("错误的命令格式:" + name + "," + ty)
		}
	}
	return ""
}

func (impl *Impl) setMethodImpl(di interface{}) {
	toe := reflect.TypeOf(di).Elem()
	voe := reflect.ValueOf(di).Elem()
	implVO := reflect.ValueOf(&impl).Elem()
	for i := 0; i < voe.NumField(); i++ {
		ft := toe.Field(i)
		runSql := strings.TrimSpace(string(ft.Tag))
		if len(runSql) < 1 {
			runSql = getSqlByMethod(ft)
			if len(runSql) < 1 {
				seelog.Error("方法定义错误,没有设置RunSql:", ft.Name, ";", ft.Type.String())
				continue
			} else {
				seelog.Info("gen sql:" + runSql)
			}
		}
		if strings.Index(ft.Type.String(), "(") < 0 {
			seelog.Error("不是函数" + ft.Name)
			continue
		}
		rt := ft.Type.String()
		rt = strings.TrimSpace(rt[strings.Index(rt, ")")+1:])

		fv := voe.Field(i)
		methodName := ""

		rts := strings.Replace(rt[1:len(rt)-1], "[]", "array_", -1)
		rts = strings.Replace(rts, ", ", "_", -1)
		rts = strings.Replace(rts, "_error", "", -1)
		rts = strings.Replace(rts, " {}", "", -1)
		rts = strings.Replace(rts, "[", "_", -1)
		rts = strings.Replace(rts, "]", "_", -1)
		rtArray := strings.Split(rts, "_")
		obj := false
		for _, r := range rtArray {
			if len(r) > 1 {
				if strings.Index(r, ".") > 0 {
					methodName += "Object"
					obj = true
				} else {
					methodName += strings.ToUpper(r[0:1]) + r[1:]
				}
			}
		}
		if obj {
			fv.Set(reflect.MakeFunc(fv.Type(), func(in []reflect.Value) ([]reflect.Value) {
				v := vti(in)
				rows, err := impl.conn.Query(runSql, v...)
				defer rows.Close()
				if err != nil {
					seelog.Error("调用SQL出错了\n\t", impl.dsn, "\n\t", runSql, vti(in), "\n\t", err)
					seelog.Flush()
					return []reflect.Value{nilVf, nilVf}
				} else {
					cols, _ := rows.Columns()
					return impl.QueryObjectBool(rows, cols, ft.Type.Out(0))
				}
			}))
		} else {
			lowSql := strings.ToLower(runSql)[0:6]
			if lowSql == "insert" || lowSql == "update" || lowSql == "delete" {
				methodName = "Exec" + methodName
			} else {
				methodName = "Query" + methodName
			}
			implM, b := implVO.Type().MethodByName(methodName)
			if b {
				fv.Set(reflect.MakeFunc(fv.Type(), func(in []reflect.Value) ([]reflect.Value) {
					defer func() {
						if err := recover(); err != nil {
							seelog.Error("methodName=", methodName, ";runSql=", runSql, err)
						}
					}()
					if methodName[0:5] == "Query" {
						rows, err := impl.conn.Query(runSql, vti(in)...)
						defer rows.Close()
						if err != nil {
							seelog.Error("调用SQL出错了", runSql, err)
							return []reflect.Value{nilVf, reflect.ValueOf(err)}
						} else {
							cols, _ := rows.Columns()
							params := make([]reflect.Value, 2)
							params[0] = reflect.ValueOf(rows)
							params[1] = reflect.ValueOf(cols)
							return implVO.Method(implM.Index).Call(params)
						}
					} else {
						params := make([]reflect.Value, 2)
						params[0] = reflect.ValueOf(runSql)
						params[1] = reflect.ValueOf(vti(in))
						return implVO.Method(implM.Index).Call(params)
					}
				}))
			} else {
				msg := "方法没有现实:\nfunc (impl *Impl) " + methodName + "(rows sql.Rows, cols []string) " + rt + "{\n\t\n}\n;" + ft.Name + ";sql=" + runSql
				seelog.Error(msg)
				panic(msg)
			}
		}
	}
}

func (impl *Impl) GetRows(runSql string, param ... interface{}) (*sql.Rows, error) {
	return impl.conn.Query(runSql, param...)
}
