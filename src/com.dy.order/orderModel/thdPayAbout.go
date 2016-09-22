package orderModel

import (
	"com.dy.order/common"
	"errors"
	"github.com/Centny/gwf/log"
	//"com.dy.tool/dbMgr"
	"net/http"
	//"net/url"
	//"com.dy.order/conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

/*
callback
*/
func Callback(ono string) (bl bool, err error) {
	log.D("callback")
	var strurl1 string
	var strurl2 string

	//strurl := fmt.Sprintf("%s%s", `http://`, conf.Rcp_host())
	db := common.DbConn()
	_sql := `select aval from ods_order_env where akey in (select p_from from ods_order_item where ono=?)`
	if err := db.QueryRow(_sql, ono).Scan(&strurl1); err != nil {
		log.E("query aval err in ods_order_env")
		return false, err
	}
	_sql = `select expand from ods_order where ono=?`
	if err := db.QueryRow(_sql, ono).Scan(&strurl2); err != nil {
		log.E("query aval err in ods_order")
		return false, err
	}
	strurl := strurl1 + strurl2
	//strurl := strurl + strurl1 + strurl2
	//Rcp_host

	fmt.Printf("strurl:%s\n", strurl)
	if strurl != "" {
		res, err := http.Get(strurl)
		if err != nil {
			log.E("callback:get err:%s", err.Error())
			return false, err
		}
		defer res.Body.Close()
		got, err := ioutil.ReadAll(res.Body)
		var cb CBStruct
		if err := json.Unmarshal(got, &cb); err != nil {
			log.E("json err:%s", err.Error())
			return false, nil
		}
		log.D("cb.Code:%v,cb.Data:%v", cb.Code, cb.Data)
		if cb.Code == int64(0) {
			return true, nil
		}
		if cb.Code == int64(2) {
			//参与了课程，可认为成功?
			log.D("你已经参与过这门课程")
			return true, nil
		}
		return false, nil
	}
	//如果没返回，就不管
	return true, nil
	//id=40001&token=5e6248a918eb211ab85381c6499adeb8-db481955-9910-4db3-aa6c-b401f3831743
}

func IfAlipaySuccessFail(name string, Type string, money float64, pay_type string, targetid string, target_type string, ono string, status string) {
	//	conn = dbMgr.DbConn()
	log.D("数据库出错，后续处理中")
	go func() {
		timer := time.NewTicker(15 * time.Second)
		i := 0
	breakf:
		for {
			select {
			case <-timer.C:
				if i < 3 {
					log.D("the %d times  ccb fail", i)
					i++
				} else {
					log.D("kill timer ali fail")
					timer.Stop()
					break breakf
				}
				if bl, _ := AlipayPaySuccess(name, Type, money, pay_type, targetid, target_type, ono, status); bl {
					log.D("ccb success in fail")
					timer.Stop()
					break breakf
				}
			}
		}
	}()
}

/*
func:after pay success,but not refund
*/
func AlipayPaySuccess(name string, Type string, money float64, pay_type string, targetid string, target_type string, ono string, status string) (bool, error) {
	callback := false
	defer func() {
		if callback == true {
			if bl, _ := Callback(ono); bl != true {
				log.D("callback=false")
				go func(ono string) {
					log.D("回调出错，后续处理中")
					timer := time.NewTicker(3 * time.Second)
					i := 0
				breakf:
					for {
						select {
						case <-timer.C:
							if i >= 5 {
								timer.Stop()
								log.D("kill timer ali")
								break breakf
							} else {
								log.D("the %d times ccb", i)
								i++
							}
							if bl, _ := Callback(ono); bl {
								log.D("ccb success ")
								timer.Stop()
								break breakf
							}
						}
					}
				}(ono)
			} else {
				log.D("callback=true")
			}
		}
	}()

	db := common.DbConn()
	tx, _ := db.Begin()
	var uid int64
	var target_id int64
	// var imoney float64
	// var buyer int64
	_sql := `select buyer,seller from  ods_order o join ods_record r  where r.ono=o.ono and r.ono =? order by r.tid asc`
	err := tx.QueryRow(_sql, ono).Scan(&uid, &target_id)
	//err1 := tx.QueryRow(_sql, ono).Scan(&target_id)
	if err != nil /*|| err1 != nil*/ {
		log.E("Query ods_record uid ,target_id error %v", err.Error())
		tx.Rollback()
		return false, errors.New("Query record uid,target_id error")
	}
	if uid == 0 || target_id == 0 {
		tx.Rollback()
		return false, errors.New("uid or target not exist")
	}
	//integral
	//

	//buyer-->uid  record
	Type = "PAID"
	sts := "PAID"
	if bl, err := UpdateRecord(tx, Type, sts, uid, ono); bl != true {
		return false, err
	}
	//seller-->uid  record
	Type = "INCOME"
	if bl, err := UpdateRecord(tx, Type, sts, target_id, ono); bl != true {
		return false, err
	}

	_sql = `update ods_order set status=? where ono =?`
	_, err = tx.Exec(_sql, status, ono)
	if err != nil {
		log.E("Add ods_record error %v", err.Error())
		tx.Rollback()
		return false, err
	}
	_sql = `update ods_order set wno=NULL where ono =?`
	_, err = tx.Exec(_sql, ono)
	if err != nil {
		log.E("Add ods_record error %v", err.Error())
		tx.Rollback()
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		log.E("AlipayPaySuccess commit error %v", err.Error())
		tx.Rollback()
		return false, err
	}
	callback = true
	return true, nil

}

func DealAliReturn(r *http.Request) (callbackMsg string) {
	trade_status := r.FormValue("trade_status")
	out_trade_no := r.FormValue("out_trade_no")
	buyer_email := r.FormValue("buyer_email")
	subject := r.FormValue("subject")
	log.D("buyer_email is : %v ", buyer_email)
	log.D("subject is : %v ", subject)
	log.D("trade_status is : %v ", trade_status)
	log.D("out_trade_no is : %v ", out_trade_no)

	var total_fee float64
	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	if trade_status == "TRADE_SUCCESS" {
		//todo : deal the order
		log.D("trade_success in webreturn")
		var rurl string
		_sql := `select return_url from ods_order where ono=?`
		errR := common.DbConn().QueryRow(_sql, out_trade_no).Scan(&rurl)
		if errR != nil {
			log.E("can't find return_url in order")
			callbackMsg = "没有回调地址，请联系客服"
			return
		}
		strUrl := fmt.Sprintf("location.href='%s'", rurl)
		callbackMsg = "<html>"
		callbackMsg += "<head>"
		callbackMsg += "<script type='text/javascript'>"
		callbackMsg += strUrl
		callbackMsg += "</script>"
		callbackMsg += "</head>"
		callbackMsg += "<body>"

		callbackMsg += "</body>"
		callbackMsg += "</html>"
	}

	if trade_status == "TRADE_FINISHED" {

	}
	return
}

func DealAliNotify(r *http.Request, pay_type string, from string) {

	trade_status := r.FormValue("trade_status") //-
	out_trade_no := r.FormValue("out_trade_no") //ono
	buyer_email := r.FormValue("buyer_email")   //-
	subject := r.FormValue("subject")           //name
	sellerId := r.FormValue("seller_id")        //targetid
	buyerId := r.FormValue("buyer_id")          //-
	//pay_type := `ALIPAY`
	target_type := "USER"
	Type := "INCOMING" //PAY|INCOMING|REFUND
	//status := "N"
	status := "PAID" //改N->PAID
	var total_fee float64
	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	log.D("call from is: %v", from)
	log.D("trade_status is : %v ", trade_status)
	log.D("out_trade_no is : %v ", out_trade_no)
	log.D("buyer_email is : %v ", buyer_email)
	log.D("subject is : %v ", subject)
	log.D("total_fee is %v ", total_fee)
	log.D("buyerId is %v", buyerId)

	//判断该笔订单是否在商户网站中已经做过处理
	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
	//如果有做过处理，不执行商户的业务程序

	//注意：
	//该种交易状态只在一种情况下出现——开通了高级即时到账，买家付款成功后。

	if trade_status == "TRADE_SUCCESS" {
		log.D("TRADE_SUCCESS,处理订单中...")
		if bl, _ := AlipayPaySuccess(subject, Type, total_fee, pay_type, sellerId, target_type, out_trade_no, status); bl != true {
			IfAlipaySuccessFail(subject, Type, total_fee, pay_type, sellerId, target_type, out_trade_no, status)
		}
	}

	//判断是否已做操作

	//判断该笔订单是否在商户网站中已经做过处理
	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
	//如果有做过处理，不执行商户的业务程序

	//注意：
	//1、开通了普通即时到账，买家付款成功后。
	//该种交易状态只在两种情况下出现
	//2、开通了高级即时到账，从该笔交易成功时间算起，过了签约时的可退款时限（如：三个月以内可退款、一年以内可退款等）后。

	if trade_status == "TRADE_FINISHED" {

	}
}
