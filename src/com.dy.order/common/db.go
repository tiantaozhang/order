package common

import (
	"com.dy.tool/dbMgr"
	"database/sql"
)

const (
	MaxOpenConns = 16
)

//var db *sql.DB

func Init(driver string, db string) error {
	dbMgr.Init(driver, db, "order")
	dbMgr.SetMaxOpenConn("order", MaxOpenConns)
	return CheckDb(DbConn())
}

func DbConn() *sql.DB {
	return dbMgr.DbConn("order")
	//return db
}

// func init() {
// 	db, _ = ConnectTestDB("mydb")
// }

// func ConnectTestDB(dbname string) (db *sql.DB, err error) {

// 	//db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mydb?charset=utf8")
// 	mydbname := "root:123456@tcp(127.0.0.1:3306)/" + dbname + "?charset=utf8"
// 	db, err = sql.Open("mysql", mydbname)
// 	if !checkErr(err) {
// 		return
// 	}
// 	db.SetMaxOpenConns(100)
// 	db.SetMaxIdleConns(50)
// 	err = db.Ping()
// 	if !checkErr(err) {
// 		return
// 	}

// 	return
// }

// func Close(db *sql.DB) (err error) {

// 	err = db.Close()
// 	if !checkErr(err) {
// 		return
// 	}
// 	return nil
// }

// func checkErr(err error) bool {
// 	if err != nil {

// 		return false
// 	}
// 	return true
// }
