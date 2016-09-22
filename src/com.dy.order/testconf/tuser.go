package testconf

import (
	"com.dy.order/common"
	"fmt"
	"time"
)

type tUser struct {
	uid   int64
	uname string
}

func TestUser() *tUser {
	return &tUser{}
}

func (u *tUser) Id(id int64) *tUser {
	u.uid = id
	return u
}

func (u *tUser) Name(name string) *tUser {
	u.uname = name
	return u
}

func (u *tUser) Save() {
	db := common.DbConn()
	if u.uname == "" {
		u.uname = fmt.Sprintf("test%v", time.Now().UnixNano())
	}

	_, err := db.Exec(`delete from UCS_USR where tid=? or usr=?`, u.uid, u.uname)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`insert into UCS_USR(tid,usr,pwd,status,time) values(?,?,0,0,0)`, u.uid, u.uname)
	if err != nil {
		panic(err)
	}

}
