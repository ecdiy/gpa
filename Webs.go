package gpa

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"encoding/json"
)

func ginMap(c *gin.Context) (map[string]interface{}, error) {
	row, b := c.GetRawData()
	if b == nil {
		var data map[string]interface{}
		je := json.Unmarshal(row, &data)
		if je != nil {
			seelog.Error("RawData JSON error:", row, je)
		}
		return data, je
	}
	return nil, b
}

func (me *Gpa) MapInsert(table string, data map[string]interface{}, userId int64) (int64, error) {
	var vs []interface{}
	sql := "insert into " + table + "("
	s2 := ""
	for k, v := range data {
		fst := k[0:1]
		if fst >= "A" && fst <= "Z" {
			sql += k + ","
			s2 += "?,"
			vs = append(vs, v)
		}
	}
	vs = append(vs, userId)
	sql = sql + "UserId,CreateAt)values(" + s2 + "?,now())"
	return me.Insert(sql, vs...)
}

func (me *Gpa) WebInsert(c *gin.Context, verify func(c *gin.Context) (bool, int64)) {
	auth, uId := verify(c)
	if auth {
		data, ed := ginMap(c)
		if ed == nil {
			id, e := me.MapInsert(data["table"].(string), data, uId)
			if e == nil {
				data["Id"] = id
			}
		} else {
			seelog.Error("GinInsert.数据转换错误，数据不合法", ed)
			c.AbortWithStatus(500)
		}
	} else {
		c.AbortWithStatus(401)
	}
}

func (me *Gpa) WebDeleteById(c *gin.Context, verify func(c *gin.Context) (bool, int64)) {
	auth, uId := verify(c)
	if auth {
		m, _ := ginMap(c)
		sql := "delete from " + m["table"].(string) + " where Id=? and UserId=?"
		row, e := me.Exec(sql, m["Id"], uId)
		if e == nil {
			m["rowsAffected"] = row
			c.JSON(200, m)
		}
	} else {
		c.AbortWithStatus(401)
	}
}

func (me *Gpa) MapUpdate(data map[string]interface{}, userId int64) (int64, error) {
	sql := "update " + data["table"].(string) + " set "
	var vs []interface{}
	for k, v := range data {
		fst := k[0:1]
		if k != "Id" && fst >= "A" && fst <= "Z" {
			sql += k + "=?,"
			vs = append(vs, v)
		}
	}
	sql = sql[0:len(sql)-1] + ",ModifyAt=now() where Id=? and UserId=?"
	vs = append(vs, data["Id"])
	vs = append(vs, userId)
	row, ee := me.Exec(sql, vs...)
	if ee != nil {
		seelog.Error("WebUpdateById.执行SQL失败。", sql, vs, ee)
	}
	return row, ee
}

func (me *Gpa) WebUpdateById(c *gin.Context, verify func(c *gin.Context) (bool, int64)) {
	auth, uId := verify(c)
	if auth {
		data, ed := ginMap(c)
		if ed == nil {
			row, ee := me.MapUpdate(data, uId)
			if ee != nil {
				c.AbortWithStatus(500)
			} else {
				res := make(map[string]interface{})
				res["Id"] = data["Id"]
				res["rowsAffected"] = row
				c.JSON(200, data)
			}
		} else {
			seelog.Error("WebUpdateById.数据转换错误，数据不合法", ed)
			c.AbortWithStatus(500)
		}
	} else {
		c.AbortWithStatus(401)
	}
}
