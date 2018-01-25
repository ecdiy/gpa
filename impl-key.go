package gpa

func (dao *Impl) Get(key string) (string, bool, error) {
	rows, _ := dao.GetRows(`select V from KeyVal where K=?`, key)
	defer rows.Close()
	cols, _ := rows.Columns()
	return dao.QueryStringBool(rows, cols)
}

func (dao *Impl) Set(k string, v interface{}) (int64, error) {
	v, b, _ := dao.Get(k)
	if b {
		return dao.Exec(`update KeyVal set V=? where K=?`, v, k)
	} else {
		return dao.Exec("insert into KeyVal(K,V) values(?,?)", k, v)
	}
}
