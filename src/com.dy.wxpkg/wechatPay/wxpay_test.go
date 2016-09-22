package wechatPay

import (
	"fmt"
	"testing"
)

/*
 <xml>
 <return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<appid><![CDATA[wx270243445e3233ee]]></appid>
<mch_id><![CDATA[1288157701]]></mch_id>
<device_info><![CDATA[013467007045764]]></device_info>
<nonce_str><![CDATA[rsfx6xfxxTDHqX0b]]></nonce_str>
<sign><![CDATA[3174FECC76CE18F06EAD3BAD97DCDB43]]></sign>
<result_code><![CDATA[SUCCESS]]></result_code>
<prepay_id><![CDATA[wx201601081103033c9db561090996187444]]></prepay_id>
<trade_type><![CDATA[APP]]></trade_type>
</xml>

*/
func TestUnifiedorders(t *testing.T) {
	// https://pay.weixin.qq.com/wiki/tools/signverify/
	// 验证地址
	t.Log("TestNewUnifiedorder")
	u := NewUnifiedorder(GWxPayConfig)
	// 随机字符串
	u.Nonce_str = md5String(NewOrderNo())
	// 商品描述
	u.Body = "testWXPay"
	// 商户订单号
	u.Out_trade_no = NewOrderNo()
	// 总金额
	u.Total_fee = "1"
	// 终端IP
	u.Spbill_create_ip = "127.0.0.1"
	// 通知地址 Config 统一配置  接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	// u.Notify_url = "xxxx://xxx.xxx.xxx.xxx"
	// 交易类型  JSAPI--公众号支付、NATIVE--原生扫码支付、APP--app支付  MICROPAY--刷卡支付
	//u.Trade_type = "NATIVE"
	u.Trade_type = "APP"
	// 商品ID
	u.Product_id = "123456"
	//====== 选填
	// 设备号
	u.Device_info = "013467007045764"
	// 商品详情
	u.Detail = "xxoo xxoo"
	// 商品标记
	//u.Goods_tag = "超🐂B"

	uresp, err := u.TakeOrder(GWxPayConfig)
	if err != nil {
		t.Fatalf("TakeOrder fail %v", err)
	}
	_, err = u.GenPayReq(GWxPayConfig, uresp.Prepay_id, u.Out_trade_no)
	if err != nil {
		t.Fatalf("GenPayReq fail %v", err)
	}

}

// func TestUnifiedorderResponse(t *testing.T) {
// 	t.Log("TestUnifiedorderResponse")

// 	tstring := []byte(`<xml>
//    <return_code><![CDATA[SUCCESS]]></return_code>
//    <return_msg><![CDATA[OK]]></return_msg>
//    <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
//    <mch_id><![CDATA[10000100]]></mch_id>
//    <nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str>
//    <sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign>
//    <result_code><![CDATA[SUCCESS]]></result_code>
//    <prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id>
//    <trade_type><![CDATA[JSAPI]]></trade_type>
// </xml>`)

// 	_, err := ParseUResponse(tstring)
// 	if err != nil {
// 		t.Fatalf("ParseUnifiedorderResponse err %v", err)
// 	}

// }

// <xml><appid><![CDATA[wx270243445e3233ee]]></appid>
// <bank_type><![CDATA[CFT]]></bank_type>
// <cash_fee><![CDATA[1]]></cash_fee>
// <fee_type><![CDATA[CNY]]></fee_type>
// <is_subscribe><![CDATA[N]]></is_subscribe>
// <mch_id><![CDATA[1288157701]]></mch_id>
// <nonce_str><![CDATA[ff4ef04a8c68e42d9b4fc052b24bc7b7]]></nonce_str>
// <openid><![CDATA[ojs26tyfGSxu3IdT2RgzVGs-HgVE]]></openid>
// <out_trade_no><![CDATA[2016021718351688391]]></out_trade_no>
// <result_code><![CDATA[SUCCESS]]></result_code>
// <return_code><![CDATA[SUCCESS]]></return_code>
// <sign><![CDATA[C2ADC80F45AC22188CE7118E1C064159]]></sign>
// <time_end><![CDATA[20160217183528]]></time_end>
// <total_fee>1</total_fee>
// <trade_type><![CDATA[APP]]></trade_type>
// <transaction_id><![CDATA[1004430666201602173341556530]]></transaction_id>
// </xml>

func TestNewNaviteNotify(t *testing.T) {
	// 	rtn := `<xml>
	//   <appid><![CDATA[wx2421b1c4370ec43b]]></appid>
	//   <attach><![CDATA[支付测试]]></attach>
	//   <bank_type><![CDATA[CFT]]></bank_type>
	//   <fee_type><![CDATA[CNY]]></fee_type>
	//   <is_subscribe><![CDATA[Y]]></is_subscribe>
	//   <mch_id><![CDATA[10000100]]></mch_id>
	//   <nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
	//   <openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
	//   <out_trade_no><![CDATA[1409811653]]></out_trade_no>
	//   <result_code><![CDATA[SUCCESS]]></result_code>
	//   <return_code><![CDATA[SUCCESS]]></return_code>
	//   <sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
	//   <sub_mch_id><![CDATA[10000100]]></sub_mch_id>
	//   <time_end><![CDATA[20140903131540]]></time_end>
	//   <total_fee>1</total_fee>
	//   <trade_type><![CDATA[JSAPI]]></trade_type>
	//   <transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
	// </xml>`

	rtn := `<xml><appid><![CDATA[wx270243445e3233ee]]></appid>
	<bank_type><![CDATA[CFT]]></bank_type>
	<cash_fee><![CDATA[1]]></cash_fee>
	<fee_type><![CDATA[CNY]]></fee_type>
	<is_subscribe><![CDATA[N]]></is_subscribe>
	<mch_id><![CDATA[1288157701]]></mch_id>
	<nonce_str><![CDATA[ff4ef04a8c68e42d9b4fc052b24bc7b7]]></nonce_str>
	<openid><![CDATA[ojs26tyfGSxu3IdT2RgzVGs-HgVE]]></openid>
	<out_trade_no><![CDATA[2016021718351688391]]></out_trade_no>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<sign><![CDATA[C2ADC80F45AC22188CE7118E1C064159]]></sign>
	<time_end><![CDATA[20160217183528]]></time_end>
	<total_fee>1</total_fee>
	<trade_type><![CDATA[APP]]></trade_type>
	<transaction_id><![CDATA[1004430666201602173341556530]]></transaction_id>
	</xml>`

	u, sbool, err := NewNaviteNotify([]byte(rtn), GWxPayConfig)
	if err != nil {
		t.Error("解析微信请求错误: ", err.Error())
	}
	// if !sbool {
	// 	t.Error("验证微信签名错误 : ", err)
	// }
	if sbool {
		t.Error("验证微信签名错误")
	}
	fmt.Println("out_trade_no is ", u.Out_trade_no)

}
