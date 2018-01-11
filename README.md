 
## SQL 操作
```$xslt
 
 一条记录的多列模式
	UserLogin func(password, username, appId string) []string `select id,username,if(invalid=0,if(md5(CONCAT(username,',',?))=password,0,2),1) state from SysUser where username=? and appId=?`
 
 一列int 的多行
	FindNoSupRoleAppId func() []int `select id from App where id not in (select appId from SysRole where creator=0)`
	
 一列string 的多行
	FindAllDayGetList func(time int64) []string `select sId from Stock where dayGetAt<?`

对象查询    ？

对象数组查询 ？？

 insert update delete
	SysRoleDel func(roleId int64, roleId2 int64) int64 `delete from SysRole where id=? and grantId!=0 and 0=(SELECT count(*) from SysUserRole where roleId=?)`

```

##   DAO
```

Save
```