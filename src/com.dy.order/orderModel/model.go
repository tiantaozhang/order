package orderModel

import (
	"com.dy.order/common"
	//"com.dy.order/conf"
	"database/sql"
	//"encoding/json"
	"errors"
	//"fmt"
	"fmt"
	//"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
	"math/rand"
	"org.cny.uap/sync"
	"time"
	//"io/ioutil"
	//"net/http"
)

func NewOrderNo() string {
	return fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), RandInt(10000, 99999))
}

func RandInt(min int, max int) int {
	if max-min <= 0 {
		return min
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func CheckIsExist(db *sql.DB, ono string, p_id int64, schema string) (bl bool, err error) {
	//db := common.DbConn()
	IdCount := int64(0)
	//tid := int64(0)
	if p_id == 0 {
		checkIdSql := "select count(*) from " + schema + " where ono=?"
		//checkIdSql := "select tid,count(*) from " + schema + " where ono=?"
		if err = db.QueryRow(checkIdSql, ono).Scan(&IdCount); err != nil {
			return false, err
		}
	} else {
		checkIdSql := "select count(*) from " + schema + " where ono=? and p_id=?"
		if err = db.QueryRow(checkIdSql, ono, p_id).Scan(&IdCount); err != nil {
			return false, err
		}
	}
	// if err != nil {
	// 	return false, err
	// }
	if IdCount == 0 {
		return false, nil
	}
	return true, nil
}

func InsertRecord(tx *sql.Tx, name string, Type string, money float64, uid int64, pay_type string, target_id int64, target_type string, ono string, status string) (bool, error) {
	_sql := `insert into ods_record(name,type,money,uid,pay_type,target_id,target_type,ono,status) value(?,?,?,?,?,?,?,?,?)`
	//buyer
	_, err := tx.Exec(_sql, name, Type, money, uid, pay_type, target_id, target_type, ono, status)
	if err != nil {
		log.E("Add ods_record buyer error %v", err.Error())
		tx.Rollback()
		return false, err
	}
	//seller
	return true, nil
}

func InsertOrderItem(tx *sql.Tx, js CommonRemoteReqStruct, ono string) error {

	for i := 0; i < len(js.OrderItem); i++ {
		_sql := `insert into ods_order_item(ono,oid,p_name,p_id,p_type,p_img,p_count,p_from,notified,price,type,status) value(?,?,?,?,?,?,?,?,?,?,?,?)`
		_, err := tx.Exec(_sql, ono, js.OrderItem[i].Oid, js.OrderItem[i].P_name, js.OrderItem[i].P_id, js.OrderItem[i].P_type, js.OrderItem[i].P_img, js.OrderItem[i].P_count, js.OrderItem[i].P_from, js.OrderItem[i].Notified, js.OrderItem[i].Price, js.Type, js.OrderItem[i].Status)
		if err != nil {
			log.E("Add OrderItem error %v", err.Error())
			tx.Rollback()
			return err
		}
	}
	return nil
}

func InsertOdsOrder(tx *sql.Tx, js CommonRemoteReqStruct, ono string) error {
	_sql := `insert into ods_order(ono,buyer,seller,total_price,type,status,return_url,expand) value(?,?,?,?,?,?,?,?)`
	_, err := tx.Exec(_sql, ono, js.Buyer, js.Seller, js.TotalFee, js.Type, js.Status, js.Return_url, js.Expand)
	if err != nil {
		log.E("Add Order error %v", err.Error())
		tx.Rollback()
		return err
	}
	return nil
}

/*有积分*/
func InsertWithIntegral(tx *sql.Tx, js CommonRemoteReqStruct, ono string) error {
	//uid-->buyer  target_id-->seller
	if bl, err := InsertRecord(tx, js.Subject, "INCOME", float64(js.Integral), js.Buyer, "大洋币", js.Seller, "USER", ono, "NOT_PAY"); bl != true {
		log.E("insertRecord: %v", err.Error())
		tx.Rollback()
		return err
	}
	//uid-->seller  target_id-->buyer
	if bl, err := InsertRecord(tx, js.Subject, "PAY", float64(js.Integral), js.Seller, "大洋币", js.Buyer, "USER", ono, "NOT_PAY"); bl != true {
		log.E("insertRecord: %v", err.Error())
		tx.Rollback()
		return err
	}
	//deduct
	_n_integral := ^js.Integral + 1
	if bl, err := UpdateIntegral(tx, _n_integral, js.Buyer); bl != true {
		log.E("UpdateIntegral: %v", err.Error())
		tx.Rollback()
		return err
	}
	return nil
}

func InsertWithRMB(tx *sql.Tx, js CommonRemoteReqStruct, ono string, _needpay float64, payType string) error {
	//record
	//uid-->buyer  target_id-->seller
	if bl, err := InsertRecord(tx, js.Subject, "INCOME", _needpay, js.Buyer, payType, js.Seller, "USER", ono, "NOT_PAY"); bl != true {
		log.E("insertRecord: %v", err.Error())
		tx.Rollback()
		return err
	}
	//add integral
	//uid-->seller  target_id-->buyer
	if bl, err := InsertRecord(tx, js.Subject, "PAY", _needpay, js.Seller, payType, js.Buyer, "USER", ono, "NOT_PAY"); bl != true {
		log.E("insertRecord: %v", err.Error())
		tx.Rollback()
		return err
	}
	return nil
}

/*
_integral 负数减积分
*/
func UpdateIntegral(tx *sql.Tx, _integral int64, tid int64) (bool, error) {

	_score, err := GetIntegral(common.DbConn(), tid)
	if err != nil {
		tx.Rollback()
		log.E("GetIntegral err: %v", err.Error())
		return false, err
	}
	_updateIntegral := _score + _integral
	_sql := `update uap_attr set integral=? where oid=?`
	_, err = tx.Exec(_sql, _updateIntegral, tid)
	if err != nil {
		log.E("update uap_attr integral %v", err.Error())
		tx.Rollback()
		return false, err
	}
	return true, nil
}

func GetIntegral(db *sql.DB, uid int64) (int64, error) {
	var _score int64 = 0
	_sql := `select integral from uap_attr where oid=?`
	if err := db.QueryRow(_sql, uid).Scan(&_score); err != nil {
		log.E("getIntegral :%v", err.Error())
		return _score, err
	}
	return _score, nil

}

func DetectIntegral(db *sql.DB, integral int64, buyer int64, totalFee float64) error {
	//检测积分
	if integral > 0 {
		if s, err := GetIntegral(db, buyer); err != nil {
			return err
		} else if s < integral {
			return errors.New("积分不足")
		} else if totalFee < float64(integral)/100.0 {
			return errors.New("总额比积分少")
		}
	}
	return nil
}

//integral relate
// func UpdateIntegralDb(db *sql.DB, _integral int64, tid int64) (bool, error) {

// 	_score, err := GetIntegral(db, tid)
// 	if err != nil {
// 		return false, err
// 	}
// 	_updateIntegral := _score + _integral
// 	_sql := `update uap_attr set integral=? where oid=?`
// 	_, err = db.Exec(_sql, _updateIntegral, tid)
// 	if err != nil {
// 		log.E("update uap_attr integral %v", err.Error())
// 		return false, err
// 	}
// 	return true, nil
// }

//----------------积分相关-------------------
//以后会用到，请勿删除
// func ResumeIntegral(ono string) error {
// 	var imoney float64
// 	var buyer int64
// 	db := common.DbConn()
// 	_sql := `select money,o.buyer from ods_record r join ods_order o on r.uid=o.buyer and r.ono=o.ono  where o.ono=? and r.pay_type='大洋币'`
// 	err := db.QueryRow(_sql, ono).Scan(&imoney, &buyer)
// 	if err != nil {
// 		log.E("Query ods_record money error %v", err.Error())
// 		return errors.New("Query ods_record money error")
// 	}
// 	if imoney > 0 {
// 		_n_integral := int64(imoney)
// 		if bl, err := UpdateIntegralDb(db, _n_integral, buyer); bl != true {
// 			log.E("UpdateIntegral: %v", err.Error())
// 			return err
// 		}
// 	}
// 	return nil
// }

func PayTypeAndNeedPay(integral int64, totalFee float64) (float64, float64, int) {
	//integral 假定1=＝1分
	_integral := float64(integral) / 100.0
	_needpay := totalFee - _integral
	_payway := NORMAL //0:正常支付 1:有积分有💰 2:全积分
	totalFee = _needpay
	if 0 == _needpay {
		_payway = INTEGRALONLY
	} else if _needpay > 0 && integral > 0 {
		_payway = BOTH
	}
	return _integral, _needpay, _payway
}

func UpdateRecord(tx *sql.Tx, Type string, status string, uid int64, ono string) (bool, error) {

	_sql := `update ods_record set type=?,status=? where ono=? and uid=?`
	_, err := tx.Exec(_sql, Type, status, ono, uid)
	if err != nil {
		log.E("update ods_record buyer or seller error %v", err.Error())
		tx.Rollback()
		return false, err
	}
	return true, nil
}

func CheckParas(js CommonRemoteReqStruct) error {
	if "" == js.Return_url || "" == js.Expand {
		return errors.New("return_url || expand is nil")
	}
	if js.Buyer <= int64(0) || js.Seller <= int64(0) {
		return errors.New("buyer || seller is nil")
	}
	//add subject
	if "" == js.Subject {
		return errors.New("subject is nil")
	}

	return nil
}

func SynUser(seller, buyer int64) error {
	log.D("seller:%d,buyer:%d", seller, buyer)
	if _, err := sync.SyncUsr(seller, ""); err != nil {
		log.E("sync.SyncUsr seller:%v", err.Error())
		return err
	}
	if _, err := sync.SyncUsr(buyer, ""); err != nil {
		log.E("sync.SyncUsr buyer:%v", err.Error())
		return err
	}
	return nil
}

func CheckPaidcb(akey string) error {
	paid_cb := ""
	db := common.DbConn()
	_sql := `select aval from ods_order_env where akey ='` + akey + `'`
	log.D("%s", _sql)
	if err := db.QueryRow(_sql).Scan(&paid_cb); err != nil {
		log.E("query aval err in ods_order_env")
		return err
	}
	log.D("paid_cb:%s", paid_cb)
	return nil
}

// type TSt struct {
// 	Tid int64 `m2s:"TID"`
// }

// func ListData() error {
// 	ts := []TSt{}
// 	err := dbutil.DbQueryS(common.DbConn(), &ts, "select tid from ods_order where tid>?", 1)
// 	if err != nil {
// 		log.D("Err:%v", err.Error())
// 		return err
// 	}
// 	return nil
// }
