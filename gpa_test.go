package gpa

import (
	"testing"
	"time"
	"fmt"
)

type SysRole struct {
	Id       int64 `@Id,AutoIncrement`
	RoleName string
	AppId    int
	Creator  int64
	CreateAt time.Time
}

var (
	Sql = &SqlAction{}
)

type SqlAction struct {
	SysRoleDel func(roleId int64, roleId2 int64) (int64, error) `delete from SysRole where id=? and creator!=0 and 0=(SELECT count(*) from SysUserRole where roleId=?)`
	AllAppId   func() ([]int, error)                            `select id from App  `
	IntArray2  func(appId int) ([]int, error)                   `select Id,AppEnable from App where Id=?`

	//FindRole1 func(roleId int64) ([]SysRole, error) `select id, createAt  from SysRole where id=?`
	FindRole2 func() (SysRole, error) `select id, createAt  from SysRole where id=3`
	//FindRoleMap      func() (map[string]string, error)     `select id, createAt  from SysRole where id=3`
	//FindRoleMapArray func() ([]map[string]string, error)   `select id, createAt  from SysRole where id=3`

	//UserLogin func(password, username, appId string) []string `select id,username,if(invalid=0,if(md5(CONCAT(username,',',?))=password,0,2),1) state from SysUser where username=? and appId=?`
}

func Test_Gpa(t *testing.T) {
	sqlAction := &SqlAction{}
	GetGpa("base-sys-user", sqlAction)
	v, e := sqlAction.IntArray2(1)
	fmt.Println(v, e, len(v))
	if e != nil || len(v) != 2 {
		t.Error("~~~~")
	}
	ids, e := sqlAction.AllAppId()
	fmt.Println("AllAppId:", ids)
}

//func initTest() {
//	dao = InitGPA(Sql)
//}
//func Test_Sql(t *testing.T) {
//	initTest()
//
//	sr := &SysRole{RoleName: "test", AppId: 1, CreateAt: time.Now()}
//	dao.Save(sr)
//	fmt.Print(sr)
//
//	//v, _ := Sql.FindRole2()
//	////fmt.Println("~~~~~~~~~~", v)
//	//v1, e1 := Sql.FindRole1()
//	//fmt.Println("~~~~~v1~~~~~", e1, v1)
//	//vm, _ := Sql.FindRoleMap()
//	//fmt.Println("~~~~~vm~~~~~", vm)
//	//vma, _ := Sql.FindRoleMapArray()
//	//fmt.Println("~~~~~vma~~~~~", vma)
//	//Sql.SysRoleDel(3, 3)
//	//fmt.Println(Sql.FindNoSupRoleAppId())
//	//fmt.Println(Sql.IntArray2())
//}

type Animal interface {
	Speak() string
}

type Cat struct{}

func (c Cat) Speak() string {
	return "cat"
}

type Dog struct{}

func (d *Dog) Speak() string {
	return "dog"
}

func Test_BSql(t *testing.T) {
	animals := []Animal{Cat{}, &Dog{}}
	for _, animal := range animals {
		fmt.Println(animal.Speak())
	}
}
