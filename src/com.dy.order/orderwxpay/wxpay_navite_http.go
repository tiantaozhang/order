package orderwxpay

import (
	"com.dy.order/common"
	order "com.dy.order/orderModel"
	"com.dy.order/orderQr"
	"com.dy.wxpkg/wechatPay"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"io/ioutil"
	//	"net/http"
	"strconv"
)

// func HttpRespE(w http.ResponseWriter, msg string, code int) {
// 	http.Error(w, msg, code)
// }

// func HttpRespJE(w http.ResponseWriter, msg string) {
// 	msgerr := map[string]interface{}{}
// 	msgerr["code"] = "1"
// 	msgerr["data"] = msg
// 	msgbyte, _ := json.Marshal(msgerr)
// 	w.Write([]byte(msgbyte))
// }

// func HttpResp(w http.ResponseWriter, msg string) {
// 	w.Write([]byte(msg))
// }

// 发送请求
// func WxPayNavite(w http.ResponseWriter, r *http.Request) {
// 	log.D("WxPayNavite Begin")

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
// 	u.Trade_type = "NATIVE"
// 	//u.Trade_type = "APP"
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
// 		HttpRespE(w, err.Error(), 500)
// 		return
// 	}

// 	HttpResp(w, uresp.Code_url)

// 	return
// }

/*
remote call
参数: WXRemoteReqStruct json
返回：string
*/
func WxNativeRemoteCall(req string) (string, error) {

	var err error

	log.D("begin WxNativeRemoteCall remote request")
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
	log.D("integral:%f,%f", js.TotalFee, float64(js.Integral)/100.0)
	_integral := float64(js.Integral) / 100.0
	_needpay := js.TotalFee - _integral
	s := fmt.Sprintf("%.2f", _needpay)
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", errors.New("TotalFee err")
	}

	u := wechatPay.NewUnifiedorder(wechatPay.GWxNativePayConfig)
	// 随机字符串
	u.Nonce_str = wechatPay.Md5String(wechatPay.NewOrderNo())
	// 商品描述
	u.Body = js.Subject
	// 商户订单号
	u.Out_trade_no = wechatPay.NewOrderNo()
	// // 总金额
	u.Total_fee = fmt.Sprintf("%d", int(i*100))
	// // 终端IP
	u.Spbill_create_ip = "127.0.0.1"
	// // 交易类型  JSAPI--公众号支付、NATIVE--原生扫码支付、APP--app支付  MICROPAY--刷卡支付
	u.Trade_type = "NATIVE"
	// // 商品ID
	//u.Product_id = "123456"
	// //detail
	u.Detail = js.Body
	//u.Device_info = "123456"
	// //js.Ono = u.Out_trade_no
	fmt.Printf("Total_fee:%s\n", u.Total_fee)

	//sync uap
	// if err := order.SynUser(js.Seller, js.Buyer); err != nil {
	// 	return "", err
	// }
	//预先检测
	if js.TotalFee < 0 || (js.TotalFee == 0 && js.Integral == 0) {
		return "", errors.New("invalid total_fee")
	}

	//全积分暂不支持
	uresp, err := u.TakeOrder(wechatPay.GWxNativePayConfig)
	if err != nil {
		log.E("TakeOrder err:", err.Error())
		return "", err
	}

	if _, err := DealWXOrder(js, u.Out_trade_no, uresp.Prepay_id); err != nil {
		log.E("DealWXOrder err:", err.Error())
		return "", err
	}

	//qr gen
	data, err := orderQr.GenQr(uresp.Code_url, u.Out_trade_no)
	if err != nil {
		log.E("genQr err:%s", err.Error())
		return "", err
	}

	return data, nil

}

//接收通知
// func WxPayNaviteNotify(w http.ResponseWriter, r *http.Request) {
// 	log.D("WxPayNaviteNotify Begin")

// 	n := &wechatPay.NotyfyCallback{}
// 	n.Return_code = "FAIL"
// 	defer func() {
// 		log.D("WxPayNaviteNotify End")
// 		log.D("NotyfyCallback to WX : %v", n.ToXML())
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
// 	return
// }

func WxPayWebNotify(hs *routing.HTTPSession) routing.HResult {
	log.D("WxPayWebNotify Begin")

	n := &wechatPay.NotyfyCallback{}
	n.Return_code = "FAIL"
	defer func() {
		log.D("WxPayAppNotify End")
		log.D("NotyfyCallback to WX : ", n.ToXML())
		hs.W.Header().Set("Content-Type", "application/xml ")
		fmt.Fprint(hs.W, n.ToXML())
	}()

	bodyByte, _ := ioutil.ReadAll(hs.R.Body)
	u, sbool, err := wechatPay.NewNaviteNotify(bodyByte, wechatPay.GWxNativePayConfig)
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
