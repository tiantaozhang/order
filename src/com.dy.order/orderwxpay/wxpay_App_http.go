package orderwxpay

import (
	//"encoding/json"
	"com.dy.order/common"
	order "com.dy.order/orderModel"
	"com.dy.wxpkg/wechatPay"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"io/ioutil"
	//"net/http"
	"strconv"
	"time"
)

// 发送请求
/*
此函数只做测试用,勿使用也勿删
*/
// func WxPayApp(w http.ResponseWriter, r *http.Request) {
// 	log.D("WxPayApp Begin")

// 	u := wechatPay.NewUnifiedorder(wechatPay.GWxPayConfig)
// 	// 随机字符串
// 	u.Nonce_str = wechatPay.Md5String(wechatPay.NewOrderNo())
// 	// 商品描述
// 	u.Body = "testWXPay"
// 	// 商户订单号
// 	u.Out_trade_no = wechatPay.NewOrderNo()
// 	// 总金额
// 	u.Total_fee = "1"
// 	// 终端IP
// 	u.Spbill_create_ip = "127.0.0.1"
// 	// 通知地址 Config 统一配置  接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
// 	// u.Notify_url = "xxxx://xxx.xxx.xxx.xxx"
// 	// 交易类型  JSAPI--公众号支付、NATIVE--原生扫码支付、APP--app支付  MICROPAY--刷卡支付
// 	u.Trade_type = "APP"
// 	// 商品ID
// 	u.Product_id = "123456"
// 	//====== 选填
// 	// 设备号
// 	u.Device_info = "123456"
// 	// 商品详情
// 	u.Detail = "xxoo xxoo"
// 	// 附加数据
// 	//u.Attach = "你大爷我"
// 	// 货币类型
// 	//u.Fee_type = "CNY"
// 	// 交易起始时间  订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。
// 	// u.Time_start = "xxxxxxxxx"
// 	// 交易结束时间   订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。注意：最短失效时间间隔必须大于5分钟
// 	// u.Time_expire = "xxxxxxxxx"
// 	// 商品标记
// 	//u.Goods_tag = "超🐂B"
// 	// 指定支付方式
// 	//u.Limit_pay = "no_credit"
// 	// 用户标识
// 	// u.Openid = "没标识"

// 	uresp, err := u.TakeOrder(wechatPay.GWxPayConfig)
// 	if err != nil {
// 		//HttpRespE(w, err.Error(), 500)
// 		HttpRespJE(w, err.Error())
// 		return
// 	}
// 	//write database

// 	//app
// 	jp, err := u.GenPayReq(wechatPay.GWxPayConfig, uresp.Prepay_id, u.Out_trade_no)
// 	if err != nil {
// 		//HttpRespE(w, err.Error(), 500)
// 		HttpRespJE(w, err.Error())
// 		return
// 	}
// 	HttpResp(w, string(jp))

// 	return
// }

//WxMoblieRemoteCall
//weixin pay
//
//
//@url,remote request
// WXRemoteReqStruct json
//@arg,接口参数的详细描述
//  req    R    WXRemoteReqStruct json

/*
   {//样例数据
       	m := WXRemoteReqStruct{
			Ono:        "",
			Buyer:      267250,
			Seller:     438982,
			Subject:    "testWXPay",
			TotalFee:   fee,
			Body:       "test",
			Type:       "N",
			Status:     "NOT_PAY",
			Return_url: "http://rcp.dev.jxzy.com/courseDetail.html?id=40040",
			Expand:     "id=40040&token=4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35",
		}
		for i := 1; i < 3; i++ {
			stri := fmt.Sprintf("%d", i)
			str := "物品" + stri
			orderi := orderModel.Item{
				Ono:      "",
				Oid:      int64(i),
				P_name:   str,
				P_id:     int64(i),
				P_type:   "",
				P_img:    `http://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E7%99%BE%E5%BA%A6%E5%9B%BE%E7%89%87&pn=2&spn=0&di=171315887930&pi=&rn=1&tn=baiduimagedetail&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=1879444470%2C3904781009&os=340336596%2C2044119696&simid=4219135247%2C874483244&adpicid=0&ln=30&fr=ala&sme=&cg=&bdtype=0&objurl=http%3A%2F%2Fd.hiphotos.baidu.com%2Fzhidao%2Fpic%2Fitem%2F6d81800a19d8bc3e4a4c8226838ba61ea9d34592.jpg&fromurl=ippr_z2C%24qAzdH3FAzdH3Fzit1w5_z%26e3Bkwt17_z%26e3Bv54AzdH3Fq7jfpt5gAzdH3Fc0lclac9l_z%26e3Bip4s&gsm=0`,
				P_count:  1,
				P_from:   "TEST",
				Notified: 0,
				Price:    0.01,
				Type:     "N",
				Status:   "N",
			}
			m.OrderItem = append(m.OrderItem, orderi)
		}
   }
*/
//@ret,接口返回数据描述
//  str  S   string
//  err  O   error

//@tag,微信支付远程调用
//@author,zhang,2016-02-18

func WxMoblieRemoteCall(req string) (string, error) {

	var err error

	log.D("begin WxMoblieRemoteCall remote request", "heheda")
	fmt.Printf("req:%s\n", []byte(req))
	var js WXRemoteReqStruct
	err = json.Unmarshal([]byte(req), &js)
	if err != nil {
		log.E("format err")
		return "", err
	}
	fmt.Printf("js:%s\n", js)

	if err := order.CheckParas(order.CommonRemoteReqStruct(js)); err != nil {
		return "", err
	}

	//detect integral
	if err = order.DetectIntegral(common.DbConn(), js.Integral, js.Buyer, js.TotalFee); err != nil {
		return "", err
	}
	//conv totalfee
	log.D("integral:%f,%f", js.TotalFee, float64(js.Integral)/100.0)
	_integral := float64(js.Integral) / 100.0
	_needpay := js.TotalFee - _integral
	s := fmt.Sprintf("%.2f", _needpay)
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", errors.New("TotalFee err")
	}

	u := wechatPay.NewUnifiedorder(wechatPay.GWxPayConfig)
	// 随机字符串
	u.Nonce_str = wechatPay.Md5String(wechatPay.NewOrderNo())
	// 商品描述
	u.Body = js.Subject
	// 商户订单号
	u.Out_trade_no = wechatPay.NewOrderNo()
	// 总金额
	u.Total_fee = fmt.Sprintf("%d", int(i*100))
	// 终端IP
	u.Spbill_create_ip = "127.0.0.1"
	// 交易类型  JSAPI--公众号支付、NATIVE--原生扫码支付、APP--app支付  MICROPAY--刷卡支付
	u.Trade_type = "APP"
	//detail
	u.Detail = js.Body
	//js.Ono = u.Out_trade_no
	fmt.Printf("Total_fee:%s\n", u.Total_fee)

	//sync uap
	// if err := order.SynUser(js.Seller, js.Buyer); err != nil {
	// 	return "", err
	// }
	if js.TotalFee < 0 || (js.TotalFee == 0 && js.Integral == 0) {
		return "", errors.New("invalid total_fee")
	}

	//全积分暂不支持
	uresp, err := u.TakeOrder(wechatPay.GWxPayConfig)
	if err != nil {
		log.E("TakeOrder err:", err.Error())
		return "", err
	}

	if _, err := DealWXOrder(js, u.Out_trade_no, uresp.Prepay_id); err != nil {
		log.E("DealWXOrder err:", err.Error())
		return "", err
	}
	//app
	returnJs, err := u.GenPayReq(wechatPay.GWxPayConfig, uresp.Prepay_id, u.Out_trade_no)
	if err != nil {
		log.E("GenPayReq err:", err.Error())
		return "", err
	}
	return string(returnJs), nil

}

//接收通知
// func WxPayAppNotify(w http.ResponseWriter, r *http.Request) {
// 	log.D("WxPayAppNotify Begin")

// 	n := &wechatPay.NotyfyCallback{}
// 	n.Return_code = "FAIL"
// 	defer func() {
// 		log.D("WxPayAppNotify End")
// 		log.D("NotyfyCallback to WX : ", n.ToXML())
// 		w.Header().Set("Content-Type", "application/xml ")
// 		fmt.Fprint(w, n.ToXML())
// 	}()

// 	bodyByte, _ := ioutil.ReadAll(r.Body)
// 	u, sbool, err := wechatPay.NewNaviteNotify(bodyByte, wechatPay.GWxPayConfig)
// 	if err != nil {
// 		log.D("解析微信请求错误 : ", err)
// 		return
// 	}
// 	if !sbool {
// 		log.D("验证微信签名错误 : ", err)
// 		return
// 	}

// 	log.D("在这处理订单")
// 	log.D("Out_trade_no : ", u.Out_trade_no)

// 	n.Return_code = "SUCCESS"
// 	n.Return_msg = "OK"
// 	return
// }

//接收微信回调通知
func WxPayMoblieNotify(hs *routing.HTTPSession) routing.HResult {
	log.D("WxPayAppNotify Begin")

	n := &wechatPay.NotyfyCallback{}
	n.Return_code = "FAIL"
	defer func() {
		log.D("WxPayAppNotify End")
		log.D("NotyfyCallback to WX : ", n.ToXML())
		hs.W.Header().Set("Content-Type", "application/xml ")
		fmt.Fprint(hs.W, n.ToXML())
	}()

	bodyByte, _ := ioutil.ReadAll(hs.R.Body)

	u, sbool, err := wechatPay.NewNaviteNotify(bodyByte, wechatPay.GWxPayConfig)
	if err != nil {
		log.E("解析微信请求错误 : ", err)
		return routing.HRES_RETURN
	}
	if !sbool {
		log.E("验证微信签名错误 : ", err)
		return routing.HRES_RETURN
	}

	status := "PAID"
	log.D("trade_status is : %v ", u.Result_code)
	log.D("total_fee is %v ", u.Total_fee)
	log.D("Out_trade_no %v: ", u.Out_trade_no)
	log.D("Transaction_id:", u.Transaction_id)

	log.D("MOBILEPAY TRADE_SUCCESS,处理订单中...")
	if bl, _ := WXpayPaySuccess(u.Out_trade_no, status, u.Transaction_id); bl != true {
		IfWXPaySuccessFail(u.Out_trade_no, status, u.Transaction_id)
	}

	n.Return_code = "SUCCESS"
	n.Return_msg = "OK"
	return routing.HRES_RETURN

}

/*
func: paid failed
*/
func IfWXPaySuccessFail(ono string, status string, transaction_id string) {
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
					log.D("the %d times ccb wx fail", i)
					i++
				} else {
					log.D("kill timer wx fail")
					timer.Stop()
					break breakf
				}
				if bl, _ := WXpayPaySuccess(ono, status, transaction_id); bl {
					log.D("success ccb wx fail")
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
func WXpayPaySuccess(ono string, status string, transaction_id string) (bool, error) {
	callback := false
	defer func() {
		if callback == true {
			if bl, _ := order.Callback(ono); bl != true {
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
								log.D("kill timer wx")
								break breakf
							} else {
								log.D("the %d times ccb wx", i)
								i++
							}
							if bl, _ := order.Callback(ono); bl {
								log.D("success ccb wx")
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
		log.E("uid or target_id not exist")
		tx.Rollback()
		return false, errors.New("uid or target not exist")
	}

	//integral 积分暂不支持
	//_n_integral := ^js.Integral + 1
	// _sql = `select r.money,o.buyer from ods_record r join ods_order o on r.uid=o.buyer and r.ono=o.ono  where o.ono=? and r.pay_type='大洋币'`
	// err = tx.QueryRow(_sql, ono).Scan(&imoney, &buyer)
	// if err != nil {
	// 	log.E("Query ods_record money error %v", err.Error())
	// 	tx.Rollback()
	// 	return false, errors.New("Query ods_record money error")
	// }
	// _n_integral := ^int64(imoney) + 1
	// if bl, err := order.UpdateIntegral(tx, _n_integral, buyer); bl != true {
	// 	log.E("UpdateIntegral: %v", err.Error())
	// 	return false, err
	// }
	//buyer-->uid  record
	Type := "PAID"
	sts := "PAID"
	if bl, err := order.UpdateRecord(tx, Type, sts, uid, ono); bl != true {
		return false, err
	}
	//seller-->uid  record
	Type = "INCOME"
	if bl, err := order.UpdateRecord(tx, Type, sts, target_id, ono); bl != true {
		return false, err
	}

	_sql = `update ods_order set status=? where ono =?`
	_, err = tx.Exec(_sql, status, ono)
	if err != nil {
		log.E("Add ods_record error %v", err.Error())
		tx.Rollback()
		return false, err
	}

	_sql = `update ods_order set wno=? where ono =?`
	_, err = tx.Exec(_sql, transaction_id, ono)
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
