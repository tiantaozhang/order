package orderalipay

import (
	"com.dy.alipkg/alipay"
	"com.dy.order/common"
	"com.dy.wxpkg/wechatPay"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	//"log"
	//"com.dy.tool/dbMgr"
	"database/sql"
	"net/http"
	//"net/url"
	// "com.dy.order/conf"
	"com.dy.order/orderModel"
	// "io/ioutil"
	//"org.cny.uap/sync"
	"com.dy.order/orderQr"
	"strconv"
	"strings"
	// "time"
)

// func NowTime() string {
// 	return fmt.Sprintf("%s", time.Now().Format("20060102150405"))
// }

/**
请求支付
*/
var c_ono string

//A:ali N:native M:mobile W:WX
var way = []string{"AN", "AM", "WM", "WN"}

// 远程调用不用该接口
// func AlipayWebRequest(w http.ResponseWriter, r *http.Request) {

// 	alipayR := &alipay.AlipayWebRequest{
// 		OutTradeNo: NewOrderNo(), // 订单号
// 		Subject:    `迟到扣200`,     // 商品名称
// 		TotalFee:   0.01,         // 价格
// 		Body:       "test web Chinese",
// 	}

// 	// 输出的是 html 页面，会自动跳转到支付界面
// 	err := alipay.AlipayWebRequestForm(alipay.AWebConfig, alipayR, w)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

/*
remote call
参数:AlipayRemoteReqStruct json
返回：string
*/
func AlipayRemoteRequest(req string) (htm string, err error) {
	// js := AlipayRemoteReqStruct{}
	log.D("begin remote request")
	fmt.Printf("req:%s\n", []byte(req))
	var js AlipayRemoteReqStruct
	err = json.Unmarshal([]byte(req), &js)
	if err != nil {
		log.E("format err")
		return "", err
	}
	//log.D("json format right")
	if err := orderModel.CheckParas(orderModel.CommonRemoteReqStruct(js)); err != nil {
		return "", err
	}

	fmt.Println(js)
	ono := NewOrderNo()
	SetCurrentOne(ono)

	//sync uap
	if err := orderModel.SynUser(js.Seller, js.Buyer); err != nil {
		return "", err
	}

	db := common.DbConn()

	isOnoExist, err := CheckIsOnoExist(db, js.Ono, "ods_order")
	//ono is nil 为正常情况
	if err != nil {
		if err.Error() == "ono is nil" {

		} else {
			log.E("%s", err.Error())
			return "", err
		}
	}
	//检测积分
	if err = orderModel.DetectIntegral(db, js.Integral, js.Buyer, js.TotalFee); err != nil {
		return "", err
	}

	//integral 假定1=＝1分
	_, _needpay, _payway := orderModel.PayTypeAndNeedPay(js.Integral, js.TotalFee)

	log.D("integral:%f,%f", js.TotalFee, float64(js.Integral)/100.0)

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("start a transaction err,%s\n", err.Error())
		return "", errors.New("start a transaction err")
	}

	//OrderType := js.Type
	//if OrderType == "REFUND" {
	// ----------暂时用不到,以后有refund 请勿删除-----------
	// log.D("Type Refund")
	// if js.Ono == "" {
	// 	log.E("Ono is nil")
	// 	tx.Rollback()
	// 	return "", errors.New("Ono is nil")
	// }
	// if bl, err := CheckIsExist(db, js.Ono, 0, "ods_order"); bl != true {
	// 	log.E("Ono is not exist in db")
	// 	tx.Rollback()
	// 	return "", errors.New("ono is not exist:" + err.Error())
	// }
	// for i := 0; i < len(js.OrderItem); i++ {
	// 	if bl, _ := CheckIsExist(db, js.Ono, js.OrderItem[i].P_id, "ods_order_item"); bl != true {
	// 		log.E("ono or oid is not exist in order_item")
	// 		tx.Rollback()
	// 		return "", errors.New("ono or oid is not exist in order_item")
	// 	}
	// }
	// //更改order_item
	// for i := 0; i < len(js.OrderItem); i++ {
	// 	_sql := `update ods_order_item set type ='REFUNG' where ono=? and p_id=?`
	// 	_, err = tx.Exec(_sql, js.OrderItem[i].Ono, js.OrderItem[i].P_id)
	// 	if err != nil {
	// 		log.E("exec order_item refund err")
	// 		tx.Rollback()
	// 		return "", err
	// 	}
	// }

	// //insert refund
	// for i := 0; i < len(js.OrderRefund); i++ {
	// 	_sql := `insert into ods_order_refund(ono,item,content,imgs,status) value(?,?,?,?,?)`
	// 	_, err = tx.Exec(_sql, js.OrderRefund[i].Ono, js.OrderRefund[i].Item, js.OrderRefund[i].Imgs, js.OrderRefund[i].Status)
	// 	if err != nil {
	// 		log.E("Add OrderItem error %v", err.Error())
	// 		tx.Rollback()
	// 		return "", err
	// 	}
	// }
	// //record
	// //uid-->buyer  target_id-->seller
	// if bl, err := InsertRecord(tx, js.Subject, "REFUND", js.TotalFee, js.Buyer, "ALIPAY", js.Seller, "USER", ono, "NOT_PAY"); bl != true {
	// 	return "", err
	// }
	// //uid-->seller  target_id-->buyer
	// if bl, err := InsertRecord(tx, js.Subject, "REFUND", js.TotalFee, js.Seller, "ALIPAY", js.Buyer, "USER", ono, "NOT_PAY"); bl != true {
	// 	return "", err
	// }
	//} else {
	log.D("Type N")
	if js.TotalFee < 0 || (js.TotalFee == 0 && js.Integral == 0) {
		log.E("消费金额或积分有误:%v", js.TotalFee)
		tx.Rollback()
		return "", errors.New("TotalFee error")
	}

	//同一订单
	if isOnoExist == false {
		//insert item
		log.D("length orderItem:%d", len(js.OrderItem))
		log.D("%v", js.OrderItem)
		if len(js.OrderItem) < 1 {
			log.E("OrderItem is nil")
			tx.Rollback()
			return "", errors.New("OrderItem is nil")
		}
		//detect paid_cb
		if err := CheckPaidcb(js.OrderItem[0].P_from); err != nil {
			log.E("p_from err:%s", js.OrderItem[0].P_from)
			tx.Rollback()
			return "", errors.New("p_from err，系统暂不支持")
		}

		if err := orderModel.InsertOrderItem(tx, orderModel.CommonRemoteReqStruct(js), ono); err != nil {
			log.E("InsertOrderItem err:%v", err.Error())
			return "", err
		}

		//_payway 0:正常支付 1:有积分有💰 2:全积分
		if 2 == _payway || 1 == _payway {
			// //uid-->buyer  target_id-->seller
			// if bl, err := InsertRecord(tx, js.Subject, "INCOME", float64(js.Integral), js.Buyer, "大洋币", js.Seller, "USER", ono, "NOT_PAY"); bl != true {
			// 	log.E("insertRecord: %v", err.Error())
			// 	return "", err
			// }
			// //uid-->seller  target_id-->buyer
			// if bl, err := InsertRecord(tx, js.Subject, "PAY", float64(js.Integral), js.Seller, "大洋币", js.Buyer, "USER", ono, "NOT_PAY"); bl != true {
			// 	log.E("insertRecord: %v", err.Error())
			// 	return "", err
			// }
			// //deduct
			// _n_integral := ^js.Integral + 1
			// if bl, err := UpdateIntegral(tx, _n_integral, js.Buyer); bl != true {
			// 	log.E("UpdateIntegral: %v", err.Error())
			// 	return "", err
			// }
			if err := orderModel.InsertWithIntegral(tx, orderModel.CommonRemoteReqStruct(js), ono); err != nil {
				log.E("InsertWithIntegral: %v", err.Error())
				return "", err
			}

		}
		if 0 == _payway || 1 == _payway {
			// //record
			if err := orderModel.InsertWithRMB(tx, orderModel.CommonRemoteReqStruct(js), ono, _needpay, "ALIPAYW"); err != nil {
				log.E("InsertWithRMB: %v", err.Error())
				return "", err
			}
		}
	}

	//}

	//insert order
	//	status := "NOT_PAY"
	if isOnoExist == false {
		// _sql := `insert into ods_order(ono,buyer,seller,total_price,type,status,return_url,expand) value(?,?,?,?,?,?,?,?)`
		// _, err = tx.Exec(_sql, ono, js.Buyer, js.Seller, js.TotalFee, js.Type, js.Status, js.Return_url, js.Expand)
		// if err != nil {
		// 	log.E("Add Order error %v", err.Error())
		// 	tx.Rollback()
		// 	return "", err
		// }
		if err := orderModel.InsertOdsOrder(tx, orderModel.CommonRemoteReqStruct(js), ono); err != nil {
			log.E("InsertOdsOrder: %v", err.Error())
			return "", err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.E("AddOrder commit error %v", err.Error())
		tx.Rollback()
		return "", err
	}

	if _payway == 2 {
		//total integral, update data

		return "integral total", nil
	}

	alipayR := &alipay.AlipayWebRequest{
		OutTradeNo: ono,         // 订单号
		Subject:    js.Subject,  // 商品名称
		TotalFee:   js.TotalFee, // 价格
		Body:       js.Body,
	}
	// 输出的是 html 页面，会自动跳转到支付界面
	htm, err = alipay.AlipayRemoteRequestForm(alipay.AWebConfig, alipayR)
	if err != nil {
		return "", err
	}
	return htm, nil
}

//confirm pay
func ConfirmOrderPay(ono, payType, token string) (string, error) {
	log.D("=============begin corfirmPay===============")
	if ono == "" || payType == "" {
		log.E("ono or payType is nil")
		return "", errors.New("ono or payType is nil")
	}
	isIn := false //in way

	for _, v := range way {
		if payType == v {
			isIn = true
			log.D("payType match")
			break
		}
	}
	if isIn != true {
		log.E(payType + " is not support")
		return "", errors.New(payType + " is not support")
	}
	db := common.DbConn()
	//以后会增加状态，暂时用不到
	isOnoExist, err := CheckIsOnoExist(db, ono, "ods_order")
	if isOnoExist != true || err != nil {
		log.E("ono is not exist")
		return "", errors.New("ono is not exist")
	}
	//
	// if isOnoExist != true {
	// 	log.E("ono is not exist")
	// 	return "", errors.New("ono is not exist")
	// }
	status, err := CheckIsPay(db, ono, "ods_order")
	if err != nil {
		log.E(err.Error())
		return "", err
	}
	if status == "PAID" || status == "paid" {
		log.D("%v has paid", ono)
		return "", errors.New(ono + " has paid")
	}
	if token != "" {

		oldExpand, err := getExpandByOno(db, ono)

		if err != nil {
			log.E("%s", err.Error())
			return "", err
		}
		s_expand := strings.Split(oldExpand, "&")
		for i := 0; i < len(s_expand); i++ {
			ss := strings.Split(s_expand[i], "=")
			if "token" == ss[0] || "TOKEN" == ss[0] {
				s_expand[i] = ss[0] + "=" + token
				break
			}
		}
		expand := strings.Join(s_expand, "&")
		log.D("expand:%v", expand)
		if bl, err := updateExpand(db, ono, expand); bl != true {
			log.E("updateExpand err: ", err.Error())
			return "", err
		}
	}

	switch payType {
	case way[0]:
		fallthrough
	case way[1]:
		r, err := getDataAli(db, ono)
		if err != nil {
			return "", err
		}
		return confirmAliPay(payType, r)
	case way[2]:
		u, err := getDataWX(db, ono, payType)
		if err != nil {
			return "", err
		}
		return confirmWXPay(payType, u, ono)
	case way[3]:
		//WN
		u, err := getDataWX(db, ono, payType)
		if err != nil {
			return "", err
		}
		return confirmWXPay(payType, u, ono)
	default:
		return "", errors.New(payType + " is not support")
	}

}

func confirmAliPay(payType string, r *alipay.AlipayWebRequest) (string, error) {
	if payType == way[0] {
		htm, err := alipay.AlipayRemoteRequestForm(alipay.AWebConfig, r)
		if err != nil {
			log.E(err.Error())
			return "", err
		}
		return htm, nil
	} else if payType == way[1] {
		orderinfo := alipay.AlipayMobileRsaSign(alipay.AlipayMobileRequest(*r), alipay.AMobileConfig)

		fmt.Println("orderinfo", orderinfo)

		orderInfoJson := map[string]interface{}{}
		orderInfoJson["code"] = "0"
		orderInfoJson["data"] = orderinfo
		returnJs, err := json.Marshal(orderInfoJson)
		if err != nil {
			log.E("json: ", err.Error())
			return "", err
		}
		return string(returnJs), nil
	}
	return "", errors.New("payType is not marry")
}

func confirmWXPay(payType string, u *wechatPay.Unifiedorder, ono string) (string, error) {

	// uresp, err := u.TakeOrder(wechatPay.GWxPayConfig)
	// if err != nil {
	// 	log.E("TakeOrder err:", err.Error())
	// 	return "", err
	// }

	if payType == way[2] {
		//app
		uresp, err := u.TakeOrder(wechatPay.GWxPayConfig)
		if err != nil {
			log.E("TakeOrder err:", err.Error())
			return "", err
		}
		returnJs, err := u.GenPayReq(wechatPay.GWxPayConfig, uresp.Prepay_id, ono)
		if err != nil {
			log.E("GenPayReq err:", err.Error())
			return "", err
		}
		return string(returnJs), nil
	} else {
		//native
		uresp, err := u.TakeOrder(wechatPay.GWxNativePayConfig)
		if err != nil {
			log.E("TakeOrder err:", err.Error())
			return "", err
		}
		data, err := orderQr.GenQr(uresp.Code_url, ono)
		if err != nil {
			log.E("genQr err:%s", err.Error())
			return "", err
		}
		return data, nil
	}

}

func getDataAli(db *sql.DB, ono string) (*alipay.AlipayWebRequest, error) {

	var totalFee float64
	var subject string
	_sql := `select total_price from ods_order where ono=?`
	if err := db.QueryRow(_sql, ono).Scan(&totalFee); err != nil {
		log.E("query total_price err in ods_order:%v", err.Error())
		return &alipay.AlipayWebRequest{}, err
	}
	_sql = `select name from ods_record where ono=?`
	if err := db.QueryRow(_sql, ono).Scan(&subject); err != nil {
		log.E("query name err in ods_record:%v", err.Error())
		return &alipay.AlipayWebRequest{}, err
	}

	r := &alipay.AlipayWebRequest{
		OutTradeNo: ono,      // 订单号
		Subject:    subject,  // 商品名称
		TotalFee:   totalFee, // 价格
	}
	return r, nil
}

func getDataWX(db *sql.DB, ono string, payType string) (*wechatPay.Unifiedorder, error) {

	var totalFee float64
	var subject string
	_sql := `select total_price from ods_order where ono=?`
	if err := db.QueryRow(_sql, ono).Scan(&totalFee); err != nil {
		log.E("query total_price err in ods_order:%v", err.Error())
		return &wechatPay.Unifiedorder{}, err
	}
	_sql = `select name from ods_record where ono=?`
	if err := db.QueryRow(_sql, ono).Scan(&subject); err != nil {
		log.E("query name err in ods_record:%v", err.Error())
		return &wechatPay.Unifiedorder{}, err
	}
	s := fmt.Sprintf("%.2f", totalFee)
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return &wechatPay.Unifiedorder{}, errors.New("TotalFee err")
	}

	//
	if payType == way[2] {
		u := wechatPay.NewUnifiedorder(wechatPay.GWxPayConfig)
		u.Nonce_str = wechatPay.Md5String(wechatPay.NewOrderNo())
		u.Body = subject
		u.Out_trade_no = ono
		//fmt.Println("TotalFee:", totalFee*100)
		u.Total_fee = fmt.Sprintf("%d", int(i*100))
		u.Spbill_create_ip = "127.0.0.1"
		// 交易类型  JSAPI--公众号支付、NATIVE--原生扫码支付、APP--app支付  MICROPAY--刷卡支付
		u.Trade_type = "APP"
		// uresp, err := u.TakeOrder(wechatPay.GWxPayConfig)
		// if err != nil {
		// 	log.E("TakeOrder err:", err.Error())
		// 	return "", err
		// }
		//fmt.Println("u.TotalFee:", u.Total_fee)

		return u, nil
	} else {
		u := wechatPay.NewUnifiedorder(wechatPay.GWxNativePayConfig)
		// 随机字符串
		u.Nonce_str = wechatPay.Md5String(wechatPay.NewOrderNo())
		// 商品描述
		u.Body = subject
		// 商户订单号
		u.Out_trade_no = ono
		// // 总金额
		u.Total_fee = fmt.Sprintf("%d", int(i*100))
		// // 终端IP
		u.Spbill_create_ip = "127.0.0.1"
		// // 交易类型  JSAPI--公众号支付、NATIVE--原生扫码支付、APP--app支付  MICROPAY--刷卡支付
		u.Trade_type = "NATIVE"

		return u, nil
	}

}

//支付宝异步通知处理
func AlipayWebNotify(w http.ResponseWriter, r *http.Request) {
	log.D("AlipayWebNotify Begin")

	var callbackMsg = "fail"
	defer func() {
		log.D("AlipayWebNotify Notify End")
		log.D("callbackMsg to alipay notifyW: %v", callbackMsg)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, callbackMsg)
	}()

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.PostForm = nil
	r.ParseForm()

	log.D("==========================================================")
	log.D("AlipayWebNotify Request :%v", r)
	log.D("==========================================================")

	if err := alipay.VerifyWebNotify(r, alipay.AWebConfig); err != nil {
		//验证失败
		log.D("verify notify fail")
		return
	}
	orderModel.DealAliNotify(r, `ALIPAY`, "web")

	callbackMsg = "success"
	return
}

//支付宝 同步通知处理
func AlipayWebReturn(w http.ResponseWriter, r *http.Request) {
	log.D("AlipayWebReturn Begin")

	var callbackMsg = "验证失败，请联系客服"
	defer func() {
		log.D("AlipayWebReturn End")
		log.D("callbackMsg to alipay return: %v", callbackMsg)
		fmt.Fprint(w, callbackMsg)
	}()
	//	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	r.PostForm = nil
	r.ParseForm()

	log.D("==========================================================")
	log.D("AlipayWebReturn Request :%v", r)
	log.D("==========================================================")

	if err := alipay.VerifyWebNotify(r, alipay.AWebConfig); err != nil {
		//验证失败
		log.D("verify webreturn fail")
		callbackMsg = "同步跳转检验失败，请联系客服"
		return
	}
	callbackMsg = orderModel.DealAliReturn(r)

	return
}

// func GetIntegral(db *sql.DB, uid int64) (int64, error) {
// 	var _score int64 = 0
// 	_sql := `select integral from uap_attr where oid=?`
// 	if err := db.QueryRow(_sql, uid).Scan(&_score); err != nil {
// 		log.E("getIntegral :%v", err.Error())
// 		return _score, err
// 	}
// 	return _score, nil

// }

func CheckIsOnoExist(db *sql.DB, ono string, schema string) (bl bool, err error) {
	IdCount := int64(0)
	if ono != "" {
		//checkIdSql := "select count(*) from " + schema + " where ono=?"
		checkIdSql := fmt.Sprintf("%s%s%s%s%s", "select count(*) from ", schema, " where ono='", ono, "'")
		if err = db.QueryRow(checkIdSql).Scan(&IdCount); err != nil {
			log.E("%s", err.Error())
			return false, err
		}
	} else {
		return false, errors.New("ono is nil")
	}

	if IdCount == 0 {
		log.E("idcount:%d", IdCount)
		return false, nil
	}
	return true, nil
}

//check is paid
func CheckIsPay(db *sql.DB, ono string, schema string) (status string, err error) {

	if bl, _ := CheckIsOnoExist(db, ono, schema); bl != true {
		return "", errors.New("ono not exist")
	}
	_sql := fmt.Sprintf("%s%s%s%s%s", "select status from ", schema, " where ono='", ono, "'")
	err = db.QueryRow(_sql).Scan(&status)
	if err != nil {
		return "", err
	}
	return
}

// func IfAlipaySuccessFail(name string, Type string, money float64, pay_type string, targetid string, target_type string, ono string, status string) {
// 	//	conn = dbMgr.DbConn()
// 	log.D("数据库出错，后续处理中")
// 	go func() {
// 		timer := time.NewTicker(10 * time.Second)
// 		i := 0
// 		for {
// 			select {
// 			case <-timer.C:
// 				if i < 5 {
// 					i++
// 				} else {
// 					return
// 				}
// 				if bl, _ := AlipayPaySuccess(name, Type, money, pay_type, targetid, target_type, ono, status); bl {
// 					timer.Stop()
// 					return
// 				}
// 			}
// 		}
// 	}()
// }

/*
callback
*/
// func Callback(ono string) (bl bool, err error) {
// 	log.D("callback")
// 	var strurl1 string
// 	var strurl2 string

// 	strurl := fmt.Sprintf("%s%s", `http://`, conf.Rcp_host())
// 	db := common.DbConn()
// 	_sql := `select aval from ods_order_env where akey in (select p_from from ods_order_item where ono=?)`
// 	if err := db.QueryRow(_sql, ono).Scan(&strurl1); err != nil {
// 		log.E("query aval err in ods_order_env")
// 		return false, err
// 	}
// 	_sql = `select expand from ods_order where ono=?`
// 	if err := db.QueryRow(_sql, ono).Scan(&strurl2); err != nil {
// 		log.E("query aval err in ods_order")
// 		return false, err
// 	}
// 	strurl = strurl + strurl1 + strurl2
// 	//Rcp_host

// 	fmt.Printf("strurl:%s\n", strurl)
// 	if strurl != "" {
// 		res, err := http.Get(strurl)
// 		if err != nil {
// 			log.E("callback:get err:%s", err.Error())
// 			return false, err
// 		}
// 		defer res.Body.Close()
// 		got, err := ioutil.ReadAll(res.Body)
// 		var cb CBStruct
// 		if err := json.Unmarshal(got, &cb); err != nil {
// 			log.E("json err:%s", err.Error())
// 			return false, nil
// 		}
// 		if cb.Code == int64(0) {
// 			return true, nil
// 		}
// 		return false, nil
// 	}
// 	//如果没返回，就不管
// 	return true, nil
// 	//id=40001&token=5e6248a918eb211ab85381c6499adeb8-db481955-9910-4db3-aa6c-b401f3831743
// }

/*
func:after pay success,but not refund
*/
// func AlipayPaySuccess(name string, Type string, money float64, pay_type string, targetid string, target_type string, ono string, status string) (bool, error) {
// 	callback := false
// 	defer func() {
// 		if callback == true {
// 			if bl, _ := Callback(ono); bl != true {
// 				log.D("callback=false")
// 				go func(ono string) {
// 					log.D("回调出错，后续处理中")
// 					timer := time.NewTimer(5 * time.Second)
// 					i := 0
// 					for {
// 						select {
// 						case <-timer.C:
// 							if i >= 5 {
// 								timer.Stop()
// 								return
// 							} else {
// 								i++
// 							}
// 							if bl, _ := Callback(ono); bl {
// 								timer.Stop()
// 								return
// 							}
// 						}
// 					}
// 				}(ono)
// 			} else {
// 				log.D("callback=true")
// 			}
// 		}
// 	}()

// 	db := common.DbConn()
// 	tx, _ := db.Begin()
// 	var uid int64
// 	var target_id int64
// 	// var imoney float64
// 	// var buyer int64
// 	_sql := `select buyer,seller from  ods_order o join ods_record r  where r.ono=o.ono and r.ono =? order by r.tid asc`
// 	err := tx.QueryRow(_sql, ono).Scan(&uid, &target_id)
// 	//err1 := tx.QueryRow(_sql, ono).Scan(&target_id)
// 	if err != nil /*|| err1 != nil*/ {
// 		log.E("Query ods_record uid ,target_id error %v", err.Error())
// 		tx.Rollback()
// 		return false, errors.New("Query record uid,target_id error")
// 	}
// 	if uid == 0 || target_id == 0 {
// 		log.E("uid or target_id not exist")
// 		tx.Rollback()
// 		return false, errors.New("uid or target not exist")
// 	}
// 	//integral

// 	//buyer-->uid  record
// 	Type = "PAID"
// 	sts := "PAID"
// 	if bl, err := UpdateRecord(tx, Type, sts, uid, ono); bl != true {
// 		return false, err
// 	}
// 	//seller-->uid  record
// 	Type = "INCOME"
// 	if bl, err := UpdateRecord(tx, Type, sts, target_id, ono); bl != true {
// 		return false, err
// 	}

// 	_sql = `update ods_order set status=? where ono =?`
// 	_, err = tx.Exec(_sql, status, ono)
// 	if err != nil {
// 		log.E("Add ods_record error %v", err.Error())
// 		tx.Rollback()
// 		return false, err
// 	}
// 	_sql = `update ods_order set wno=NULL where ono =?`
// 	_, err = tx.Exec(_sql, ono)
// 	if err != nil {
// 		log.E("Add ods_record error %v", err.Error())
// 		tx.Rollback()
// 		return false, err
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		log.E("AlipayPaySuccess commit error %v", err.Error())
// 		tx.Rollback()
// 		return false, err
// 	}
// 	callback = true
// 	return true, nil

// }

// func InsertRecord(tx *sql.Tx, name string, Type string, money float64, uid int64, pay_type string, target_id int64, target_type string, ono string, status string) (bool, error) {
// 	_sql := `insert into ods_record(name,type,money,uid,pay_type,target_id,target_type,ono,status) value(?,?,?,?,?,?,?,?,?)`
// 	//buyer
// 	_, err := tx.Exec(_sql, name, Type, money, uid, pay_type, target_id, target_type, ono, status)
// 	if err != nil {
// 		log.E("Add ods_record buyer error %v", err.Error())
// 		tx.Rollback()
// 		return false, err
// 	}
// 	//seller
// 	return true, nil
// }
// func UpdateRecord(tx *sql.Tx, Type string, status string, uid int64, ono string) (bool, error) {

// 	_sql := `update ods_record set type=?,status=? where ono=? and uid=?`
// 	_, err := tx.Exec(_sql, Type, status, ono, uid)
// 	if err != nil {
// 		log.E("update ods_record buyer or seller error %v", err.Error())
// 		tx.Rollback()
// 		return false, err
// 	}
// 	return true, nil
// }

/*
_integral 负数减积分
*/
// func UpdateIntegral(tx *sql.Tx, _integral int64, tid int64) (bool, error) {

// 	_score, err := GetIntegral(common.DbConn(), tid)
// 	if err != nil {
// 		tx.Rollback()
// 		log.E("GetIntegral err: %v", err.Error())
// 		return false, err
// 	}
// 	_updateIntegral := _score + _integral
// 	_sql := `update uap_attr set integral=? where oid=?`
// 	_, err = tx.Exec(_sql, _updateIntegral, tid)
// 	if err != nil {
// 		log.E("update uap_attr integral %v", err.Error())
// 		tx.Rollback()
// 		return false, err
// 	}
// 	return true, nil
// }

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

func updateExpand(db *sql.DB, ono string, expand string) (bl bool, err error) {
	bl = false
	_sql := `update ods_order set expand=? where ono=?`
	_, err = db.Exec(_sql, expand, ono)
	if err != nil {
		log.E(err.Error())
		return
	}
	bl = true
	return
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
	return nil
}

func getExpandByOno(db *sql.DB, ono string) (expand string, err error) {

	_sql := `select expand from ods_order where ono=?`
	if err = db.QueryRow(_sql, ono).Scan(&expand); err != nil {
		log.E("expand err %s", err.Error())
		return "", err
	}
	return

}

func SetCurrentOne(ono string) {
	c_ono = ono
}
func GetCurrentOno() string {
	return c_ono
}

func TestAlipay(hs *routing.HTTPSession) routing.HResult {

	args, _ := alipay.GenTestData()

	args_ := ""
	for k, v := range args {
		args_ += k + "=" + v + "&"
	}
	fmt.Println("https://www.alipay.com/cooperate/gateway.do?" + args_)
	res, _ := util.HPost("https://www.alipay.com/cooperate/gateway.do", args)

	hs.W.Write([]byte(res))

	return routing.HRES_RETURN
	//	return common.MsgRes(hs, parse)
}
