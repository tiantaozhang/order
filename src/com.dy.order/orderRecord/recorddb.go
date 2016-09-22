package orderRecord

import (
	"com.dy.order/common"
	"github.com/Centny/gwf/log"
)

type orderRecord struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Money      float64 `json:"money"`
	Uid        int64   `json:"uid"`
	PayType    string  `json:"pay_type"`
	TargetId   string  `json:"target_id"`
	TargetType string  `json:"target_type"`
	Ono        string  `json:"ono"`
	Status     string  `json:"status"`
}

func AddOrderRecord(name, rType string, money float64, uid int64, payType, targetId, targetType, ono, status string) error {
	db := common.DbConn()
	tx, _ := db.Begin()
	_sql := `insert into ods_record(name,type,money,uid,pay_type,target_id,target_type,ono,status) value(?,?,?,?,?,?,?,?,?)`
	_, err := tx.Exec(_sql, name, rType, money, uid, payType, targetId, targetType, ono, status)
	if err != nil {
		log.E("AddOrderRecord error %v", err.Error())
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.E("AddOrderRecord commit error %v", err.Error())
		tx.Rollback()
		return err
	}
	return nil
}

func GetOrderRecord(uid int64) ([]orderRecord, error) {
	var rs []orderRecord
	db := common.DbConn()
	_sql := "select name,type,money,uid,pay_type,target_id,target_type,ono,status from ods_record where uid = ?"
	row, err := db.Query(_sql, uid)
	if err != nil {
		return rs, err
	}
	for row.Next() {
		var name, types, payType, targetId, targetType, onum, status string
		var usrid int64
		var money float64
		tmpRes := orderRecord{}
		row.Scan(&name, &types, &money, &usrid, &payType, &targetId, &targetType, &onum, &status)
		tmpRes.Name = name
		tmpRes.Type = types
		tmpRes.Money = money
		tmpRes.Uid = usrid
		tmpRes.PayType = payType
		tmpRes.TargetId = targetId
		tmpRes.TargetType = targetType
		tmpRes.Ono = onum
		tmpRes.Status = status
		rs = append(rs, tmpRes)
	}
	return rs, nil
}

func CheckUsrPaidItem(uid, pid int64) (bool, error) {
	_sql := "select count(a.tid) from ods_order a left join ods_order_item b on b.oid=a.tid where b.p_id=? and a.buyer=? and a.status in ('PAID','DONE')"
	db := common.DbConn()
	var count int64
	err := db.QueryRow(_sql, pid, uid).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
