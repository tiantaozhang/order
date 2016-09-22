package common

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"math/rand"
	"org.cny.uas/usr"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

func UidFromHS(hs *routing.HTTPSession) int64 {
	return hs.Kvs["USR"].(*usr.Usr).Tid
}

func NowTimeInt64() int64 {
	return time.Now().UnixNano()
}

func NowTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type Int64Slice []int64

func (a Int64Slice) Join(sep string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return fmt.Sprintf("%v", a[0])
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(fmt.Sprintf("%v", a[i]))
	}

	b := make([]byte, n)
	bp := copy(b, fmt.Sprintf("%v", a[0]))
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], fmt.Sprintf("%v", s))
	}
	return string(b)
}

type User struct {
	Id   int64
	Name string
}

func GetUserById(id int64) (*User, error) {
	db := DbConn()
	var name string
	err := db.QueryRow(`select usr from UCS_USR where tid=?`, id).Scan(&name)
	if err != nil {
		return nil, errors.New("查找用户失败")
	}
	return &User{Id: id, Name: name}, nil
}

func GetUserNameById(id int64) (string, error) {
	db := DbConn()
	var name string
	err := db.QueryRow(`select usr from UCS_USR where tid=?`, id).Scan(&name)
	if err != nil {
		return "", errors.New("查找用户失败")
	}
	return name, nil
}

func RandInt(min int, max int) int {
	if max-min <= 0 {
		return min
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

const (
	//log W second
	FUNC_RUN_TIME = 5.0
)

//log func run time
func LogRunTime(arg_funcName string, arg_sTime time.Time) {
	rTime := time.Now().Sub(arg_sTime)
	if rTime.Seconds() > FUNC_RUN_TIME {
		log.W("func %v run time is %v>%v,", arg_funcName, rTime, FUNC_RUN_TIME)
	}
	if rTime.Seconds() > 0.5 {
		log.D("func %v run time >0.5s is %v", arg_funcName, rTime)
	}
}

//time string to time
func StringToTime(arg_time string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation(TimeFormat, arg_time, loc)
	if nil != err {
		log.E("%v%v", err, t)
		return t, err
	}
	return t, nil
}

//check appId is exist
func CheckAppId(db *sql.DB, arg_appId string) (int64, error) {
	appIdCount := int64(0)
	checkAppIdSql := "SELECT COUNT(TID) FROM ods_app WHERE APP_ID=? AND STATUS='NORMAL'"
	log.D("checkAppId sql is :%v :%v ", checkAppIdSql, arg_appId)
	err := db.QueryRow(checkAppIdSql, arg_appId).Scan(&appIdCount)
	if nil != err {
		return appIdCount, err
	}
	if appIdCount == 0 {
		return appIdCount, errors.New("来源系统不存在！")
	}
	return appIdCount, nil
}
