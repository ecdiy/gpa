package gpa

import (
	"github.com/cihub/seelog"
	"reflect"
	"strings"
)

/**
返回影响记录的行数
 */
func (dao Impl) Exec(runSql string, param ... interface{}) (int64, error) {
	return dao.ExecInt64(runSql, param)
}

func (dao *Impl) ExecInt64(runSql string, p []interface{}) (int64, error) {
	row, er := dao.conn.Exec(runSql, p...)
	if er == nil {
		ra, _ := row.RowsAffected()
		return ra, nil
	} else {
		seelog.Error("SQL执行失败:", runSql, er)
		return -1, er
	}
}

func (dao *Impl) Save(model interface{}) (int64, error) {
	toe := reflect.TypeOf(model).Elem()
	voe := reflect.ValueOf(model).Elem()
	e, err := dao.exist(toe, voe)
	if err != nil {
		return -1, err
	} else {
		if e == 1 {
			return dao.update(toe, voe)
		} else {
			return dao.insert(toe, voe)
		}
	}
}

/**
返回 自增ID
 */
func (dao Impl) Insert(s string, param ... interface{}) (int64, error) {
	row, err := dao.conn.Exec(s, param...)
	if err == nil {
		return row.LastInsertId()
	} else {
		seelog.Error("insert对象失败,", s, param, err)
		return -1, err
	}
}

func (dao *Impl) exist(toe reflect.Type, voe reflect.Value) (int64, error) {
	n := toe.NumField()
	s := "select 1 from " + toe.Name() + " where "
	var param []interface{}
	for i := 0; i < n; i++ {
		fx := toe.Field(i)
		tag := string(fx.Tag)
		col := strings.ToLower(fx.Name[0:1]) + fx.Name[1:]
		vF := voe.FieldByName(fx.Name).String()
		if len(vF) >= 1 {
			if strings.Index(tag, PrimaryId) >= 0 {
				s += col + "=? and "
				param = append(param, vF)
			}
		}
	}
	s = s[0:len(s)-4]
	//fmt.Println(s)
	rows, err := dao.conn.Query(s, param...)
	defer rows.Close()
	if err != nil {
		seelog.Error("SQL出错", s, ":", err)
		return -1, err
	}
	if rows.Next() {
		var res int64
		rows.Scan(&res)
		return res, nil
	} else {
		return 0, nil
	}
}

func (dao *Impl) insert(toe reflect.Type, voe reflect.Value) (int64, error) {
	n := toe.NumField()
	s := "insert into " + toe.Name() + "("
	cols := ""
	var param []interface{}
	auto := ""
	for i := 0; i < n; i++ {
		fx := toe.Field(i)
		tag := string(fx.Tag)
		vF := voe.FieldByName(fx.Name).String()
		if len(vF) == 0 {
			continue
		}
		c := strings.ToLower(fx.Name[0:1]) + fx.Name[1:]
		if strings.Index(tag, AutoIncrement) < 0 {
			s += c + ","
			cols += "?,"
			param = append(param, voe.FieldByName(fx.Name).Interface())
		} else {
			auto = fx.Name
		}
	}

	s = s[0:len(s)-1] + ")values(" + cols[0:len(cols)-1] + ")"
	row, err := dao.conn.Exec(s, param...)
	if err == nil {
		ria, _ := row.RowsAffected()
		if ria > 0 {
			rii, _ := row.LastInsertId()
			if rii > 0 && len(auto) > 0 {
				voe.FieldByName(auto).SetInt(rii)
			}
		}
		return ria, err
	} else {
		seelog.Error("insert对象失败,", s, param, err)
		return -1, err
	}
}

func (dao *Impl) update(toe reflect.Type, voe reflect.Value) (int64, error) {
	n := toe.NumField()
	s := "update " + toe.Name() + " set "
	pri := ""
	var param, priParam []interface{}
	for i := 0; i < n; i++ {
		fx := toe.Field(i)
		tag := string(fx.Tag)
		col := strings.ToLower(fx.Name[0:1]) + fx.Name[1:]
		vF := voe.FieldByName(fx.Name).String()
		if len(vF) >= 1 {
			if strings.Index(tag, PrimaryId) >= 0 || strings.Index(tag, AutoIncrement) >= 0 {
				pri += col + "=? and "
				priParam = append(priParam, vF)
			} else {
				s += col + "=?,"
				param = append(param, vF)
			}
		}
	}
	s = s[0:len(s)-1] + " where " + pri[0:len(pri)-4]
	//fmt.Println(s)
	for _, v := range priParam {
		param = append(param, v)
	}
	row, er := dao.conn.Exec(s, param...)
	if er == nil {
		raf, _ := row.RowsAffected()
		return raf, er
	} else {
		seelog.Error("update对象失败,", er)
		return -1, er
	}
}
