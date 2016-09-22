package orderalipay

import (
	"encoding/json"
	"fmt"
	//"github.com/ljy2010a/go_alipay"
	"com.dy.alipkg/alipay"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	//"log"
	"math/rand"
	"net/http"
	"strconv"
	//"strings"
	"com.dy.order/common"
	"com.dy.order/orderModel"
	"errors"
	"time"
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

//POST
//支付宝回调处理
func AlipayMobileNotify(w http.ResponseWriter, r *http.Request) {
	log.D("-----AlipayMobileNotify Notify Begin------")

	var callbackMsg = "fail"
	defer func() {
		log.D("alipay Notify End")
		log.D("callbackMsg to alipayM : %v", callbackMsg)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, callbackMsg)
	}()

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.PostForm = nil
	r.ParseForm()

	fmt.Println("==========================================================")
	fmt.Println("Request :%v", r)
	fmt.Println("==========================================================")

	if err := alipay.VerifyMobileNotify(r, alipay.AMobileConfig); err != nil {
		//验证失败
		log.D("verify notify fail")
		return
	}
	orderModel.DealAliNotify(r, `ALIPAY`, "mobile")

	callbackMsg = "success"
	return
}

//获取支付宝签名,整个字符串&
func GetRsaSign(hs *routing.HTTPSession) routing.HResult {

	var (
		orderNos string
		//	orderNo  []string
		err error
	)

	err = hs.ValidCheckVal(`
		tradeNo,R|S,L:0;
		`, &orderNos)
	if err != nil {
		return hs.MsgResE(1, err.Error())
	}

	fmt.Println("orderNos = ", orderNos)

	//tradeNo := hs.R.FormValue("tradeNo")
	sumstr := hs.R.FormValue("totalFee")
	sum, err := strconv.ParseFloat(sumstr, 64)
	if err != nil {
		return hs.MsgResE(1, err.Error())
	}
	if sum <= 0 {
		return hs.MsgResE(1, fmt.Sprintf("%s", "totalFee < 0"))
	}

	amr := alipay.AlipayMobileRequest{}
	amr.OutTradeNo = NewOrderNo()
	amr.Subject = "掌上学园商品"
	amr.Body = "掌上学园商品"
	amr.TotalFee = sum

	fmt.Println(amr)
	//fmt.Printf("alipay.AMobileConfig:%s\n", alipay.AMobileConfig)
	orderinfo := alipay.AlipayMobileRsaSign(amr, alipay.AMobileConfig)

	fmt.Println("orderinfo", orderinfo)
	//fmt.Println(hs.MsgRes(orderinfo))
	return hs.MsgRes(orderinfo)
	//return routing.HRES_RETURN
}

//获取支付宝签名,整个字符串&
func GetRsaSignJson(req string) (string, error) {

	var err error

	log.D("begin GetRsaSignJson remote request")
	fmt.Printf("req:%s\n", []byte(req))
	var js AlipayRemoteReqStruct
	err = json.Unmarshal([]byte(req), &js)
	if err != nil {
		log.E("format err")
		return "", err
	}
	fmt.Printf("js:%s\n", js)

	if err := orderModel.CheckParas(orderModel.CommonRemoteReqStruct(js)); err != nil {
		return "", err
	}

	//detect integral
	if err = orderModel.DetectIntegral(common.DbConn(), js.Integral, js.Buyer, js.TotalFee); err != nil {
		return "", err
	}

	ono := NewOrderNo()
	//deal req data
	if _, err := DealReqData(&js, ono); err != nil {
		log.E("DealReqData err:", err)
		return "", err
	}
	//全积分暂不支持
	//
	orderinfo := ""
	if js.TotalFee > 0 {
		amr := alipay.AlipayMobileRequest{}
		amr.OutTradeNo = ono
		amr.Subject = js.Subject // 商品名称
		amr.Body = js.Body
		amr.TotalFee = js.TotalFee // 价格

		orderinfo = alipay.AlipayMobileRsaSign(amr, alipay.AMobileConfig)

		fmt.Println("orderinfo:", orderinfo)
	} else {
		orderinfo = "success"
	}
	orderInfoJson := map[string]interface{}{}
	orderInfoJson["code"] = "0"
	orderInfoJson["data"] = orderinfo
	returnJs, err := json.Marshal(orderInfoJson)
	if err != nil {
		log.E("json:", err.Error())
		return "", err
	}
	return string(returnJs), nil
}

func DealReqData(js *AlipayRemoteReqStruct, ono string) (string, error) {

	//sync uap
	if err := orderModel.SynUser(js.Seller, js.Buyer); err != nil {
		return "", err
	}

	db := common.DbConn()
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("start a transaction err,%s\n", err.Error())
		return "", err
	}
	//isOnoExist, isExistErr := CheckIsExist(db, js.Ono, 0, "ods_order")
	isOnoExist, _ := orderModel.CheckIsExist(db, js.Ono, 0, "ods_order")
	//integral 假定1=＝1分

	_, _needpay, _payway := orderModel.PayTypeAndNeedPay(js.Integral, js.TotalFee)

	//======以后会有refund 请勿删除=====
	//OrderType := js.Type
	// if OrderType == "REFUND" {
	// 	log.D("Type Refund")
	// 	if js.Ono == "" {
	// 		log.E("Ono is nil")
	// 		tx.Rollback()
	// 		return "", errors.New("Ono is nil")
	// 	}
	// 	if isOnoExist != true {
	// 		log.E("Ono is not exist in db")
	// 		tx.Rollback()
	// 		return "", errors.New("ono is not exist:" + isExistErr.Error())
	// 	}
	// 	for i := 0; i < len(js.OrderItem); i++ {
	// 		if bl, _ := CheckIsExist(db, js.Ono, js.OrderItem[i].P_id, "ods_order_item"); bl != true {
	// 			log.E("ono or oid is not exist in order_item")
	// 			tx.Rollback()
	// 			return "", errors.New("ono or oid is not exist in order_item")
	// 		}
	// 	}
	// 	//更改order_item
	// 	for i := 0; i < len(js.OrderItem); i++ {
	// 		_sql := `update ods_order_item set type ='REFUNG' where ono=? and p_id=?`
	// 		_, err = tx.Exec(_sql, js.OrderItem[i].Ono, js.OrderItem[i].P_id)
	// 		if err != nil {
	// 			log.E("exec order_item refund err")
	// 			tx.Rollback()
	// 			return "", err
	// 		}
	// 	}

	// 	//insert refund
	// 	for i := 0; i < len(js.OrderRefund); i++ {
	// 		_sql := `insert into ods_order_refund(ono,item,content,imgs,status) value(?,?,?,?,?)`
	// 		_, err = tx.Exec(_sql, js.OrderRefund[i].Ono, js.OrderRefund[i].Item, js.OrderRefund[i].Imgs, js.OrderRefund[i].Status)
	// 		if err != nil {
	// 			log.E("Add OrderItem error %v", err.Error())
	// 			tx.Rollback()
	// 			return "", err
	// 		}
	// 	}
	// 	//record
	// 	//uid-->buyer  target_id-->seller
	// 	if bl, err := InsertRecord(tx, js.Subject, "REFUND", js.TotalFee, js.Buyer, "ALIPAY", js.Seller, "USER", ono, "NOT_PAY"); bl != true {
	// 		return "", err
	// 	}
	// 	//uid-->seller  target_id-->buyer
	// 	if bl, err := InsertRecord(tx, js.Subject, "REFUND", js.TotalFee, js.Seller, "ALIPAY", js.Buyer, "USER", ono, "NOT_PAY"); bl != true {
	// 		return "", err
	// 	}

	// } else {
	log.D("Type N")
	if js.TotalFee < 0 || (js.TotalFee == 0 && js.Integral == 0) {
		log.E("消费金额<=0:%v", js.TotalFee)
		return "", errors.New("TotalFee error")
	}
	//同一订单
	if !isOnoExist {
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
			log.E("p_from err")
			tx.Rollback()
			return "", errors.New("p_from err")
		}

		if err := orderModel.InsertOrderItem(tx, orderModel.CommonRemoteReqStruct(*js), ono); err != nil {
			log.E("InsertOrderItem err:%v", err.Error())
			return "", err
		}
		//_payway 0:正常支付 1:有积分有💰 2:全积分
		if 2 == _payway || 1 == _payway {
			// //uid-->buyer  target_id-->seller
			if err := orderModel.InsertWithIntegral(tx, orderModel.CommonRemoteReqStruct(*js), ono); err != nil {
				log.E("InsertWithIntegral: %v", err.Error())
				return "", err
			}
		}
		if 0 == _payway || 1 == _payway {
			//record
			//uid-->buyer  target_id-->seller
			// if bl, err := InsertRecord(tx, js.Subject, "INCOME", _needpay, js.Buyer, "ALIPAY", js.Seller, "USER", ono, "NOT_PAY"); bl != true {
			// 	return "", err
			// }
			// //uid-->seller  target_id-->buyer
			// if bl, err := InsertRecord(tx, js.Subject, "PAY", _needpay, js.Seller, "ALIPAY", js.Buyer, "USER", ono, "NOT_PAY"); bl != true {
			// 	return "", err
			// }
			if err := orderModel.InsertWithRMB(tx, orderModel.CommonRemoteReqStruct(*js), ono, _needpay, "ALIPAYM"); err != nil {
				log.E("InsertWithRMB: %v", err.Error())
				return "", err
			}
		}
	}
	// }
	//insert order
	//	status := "NOT_PAY"
	if !isOnoExist {
		// _sql := `insert into ods_order(ono,buyer,seller,total_price,type,status,return_url,expand) value(?,?,?,?,?,?,?,?)`
		// _, err = tx.Exec(_sql, ono, js.Buyer, js.Seller, js.TotalFee, js.Type, js.Status, js.Return_url, js.Expand)
		// if err != nil {
		// 	log.E("Add Order error %v", err.Error())
		// 	tx.Rollback()
		// 	return "", err
		// }
		if err := orderModel.InsertOdsOrder(tx, orderModel.CommonRemoteReqStruct(*js), ono); err != nil {
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
	return "", nil
}

func MobilePayTest(w http.ResponseWriter, r *http.Request) {
	m := AlipayRemoteReqStruct{
		Ono:        "",
		Buyer:      267250,
		Seller:     438982,
		Subject:    "功夫熊猫3",
		TotalFee:   0.01,
		Body:       "指纹打卡",
		Type:       "N",
		Status:     "NOT_PAY",
		Return_url: "http://rcp.dev.jxzy.com/questionPoolDetailNew.html?id=40781&eid=55651",
		Expand:     "id=40781&token=2b8294d19c4a0c9729b685abbd1e3341-9767721f-75bc-47f2-b46a-f1377bae47f4",
	}
	orderi := orderModel.Item{
		Ono:      "",
		Oid:      int64(10000),
		P_name:   "balabalabala",
		P_id:     int64(11111),
		P_type:   "",
		P_count:  1,
		P_from:   "TEST",
		Notified: 0,
		Price:    0.01,
		Type:     "N",
		Status:   "N",
	}
	m.OrderItem = append(m.OrderItem, orderi)
	bytes, _ := json.Marshal(m)
	var str string
	var err error
	if str, err = GetRsaSignJson(string(bytes)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(w, str)

}

/*
_input_charset="UTF-8"
&body="掌上学园商品"
&notify_url="http%3A%2F%2Fmall.kuxiao.cn%2FalipayMobile%2Fnotify"
&out_trade_no="2015122215125373432"
&partner="2088501949844011"
&payment_type="1"
&seller_id="itdayang@gmail.com"
&service="mobile.securitypay.pay"
&subject="掌上学园商品"
&total_fee="0.01"
&sign="EDH9zspz8IyKAoJPKoHePpTzooYp5IXEFAmBIrnW3D50ESj3gsXoPZ%2FxWN4Qmcwav0m69F7whckY23hiWoLUirC%2BEBM9JHfYCqUUsQvDd6GNcsul9y31cBB9bAbs%2FJkBqJWSdG%2FEEBfwHBgflvNZKStBoqwqE6CDWDojyi8RLhg%3D"
&sign_type="RSA"",
*/
