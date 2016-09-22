package demo

import (
	"fmt"
	//alipay "github.com/ljy2010a/go_alipay"
	"com.dy.aliPkgT/alipay"
	//"github.com/Centny/gwf/routing"
	//"github.com/Centny/gwf/util"
	"log"
	"net/http"
)

/**
请求支付
*/
func AlipayWebRequest(w http.ResponseWriter, r *http.Request) {

	alipayR := &alipay.AlipayWebRequest{
		OutTradeNo: NewOrderNo(), // 订单号
		Subject:    `迟到扣200`,     // 商品名称
		TotalFee:   0.01,         // 价格
		Body:       "id=39940",
		ShowUrl:    "http://www.alipay.com/",
	}
	//alipay.AWebConfig.Return_url = fmt.Sprintf("%s%s", alipay.AWebConfig.Return_url, "id=39940")
	// 输出的是 html 页面，会自动跳转到支付界面
	err := alipay.AlipayWebRequestForm(alipay.AWebConfig, alipayR, w)
	if err != nil {
		return
	}
	return
}

/*
Remote call
*/
/*
func AlipayRemoteRequestForm(alipayConfig *AlipayConfig, r *AlipayWebRequest) (htm string, err error) {
	p := Kvpairs{
		Kvpair{`total_fee`, fmt.Sprintf("%.2f", r.TotalFee)},
		Kvpair{`subject`, r.Subject},
		Kvpair{`body`, r.Body},
		Kvpair{`show_url`, r.ShowUrl},
		Kvpair{`out_trade_no`, r.OutTradeNo},
		Kvpair{`service`, alipayConfig.Service},
		Kvpair{`partner`, alipayConfig.Partner},
		Kvpair{`payment_type`, alipayConfig.Payment_type},
		Kvpair{`notify_url`, alipayConfig.Notify_url},
		Kvpair{`return_url`, alipayConfig.Return_url},
		Kvpair{`seller_email`, alipayConfig.Seller_id},
		Kvpair{`_input_charset`, alipayConfig.Input_charset},
		// Kvpair{`anti_phishing_key`,}
		// Kvpair{`exter-invoke_ip`,}
		//	Kvpair{`price`, "0.01"},
	}

	paraFilter(&p)
	argSort(&p)
	sign := md5Sign(createLinkStringNoUrl(&p), alipayConfig.Key)

	p = append(p, Kvpair{`sign`, sign})
	p = append(p, Kvpair{`sign_type`, `MD5`})
	//p = append(p, Kvpair{`sign_type`, `GBK`})

	// htm = fmt.Sprintf("%s", `<html><head>
	// <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	// </head><body>`)
	// htm += fmt.Sprintf(`<form name='alipaysubmit' action='%s_input_charset=%s' method='post'> `, alipayGatewayNew, alipayConfig.Input_charset)
	// for _, kv := range p {
	// 	htm += fmt.Sprintf(`<input type='hidden' name='%s' value='%s' />`, kv.K, kv.V)
	// }
	// //
	// htm += fmt.Sprintf(`</form>`)
	// htm += fmt.Sprintf(`<script>document.forms['alipaysubmit'].submit();</script>`)
	// htm += fmt.Sprintf(`</body></html>`)

	return htm, nil
}
*/
/*
callback
*/
func Callback() (bl bool, err error) {
	log.Println("callback")
	var strurl string = `http://rcp.dev.jxzy.com/usr/purchase-course?id=39940`
	// _sql := `select aval from ods_order_env where akey in (select p_from from ods_order_item where ono=?)`
	// if err := common.DbConn().QueryRow(_sql, ono).Scan(&strurl); err != nil {
	// 	log.E("query aval err in ods_order_env")
	// 	return false, err
	// }
	// fmt.Printf("strurl:%s\n", strurl)
	// if strurl != "" {
	_, err = http.Get(strurl)
	if err != nil {
		log.Println("callback:postfrom err")
		return false, err
	}
	// }
	//如果没返回，就不管
	return true, nil
}

//支付宝异步通知处理
func AlipayWebNotify(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AlipayWebNotify Begin")

	callbackMsg := "fail"
	defer func() {
		log.Println("AlipayWebNotify Notify End")
		log.Println("callbackMsg to alipay in notify_url:", callbackMsg)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, callbackMsg)
	}()

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.PostForm = nil
	r.ParseForm()

	log.Println("==========================================================")
	log.Println("AlipayWebNotify Request :%v", r)
	log.Println("==========================================================")

	if err := alipay.VerifyWebNotify(r, alipay.AWebConfig); err != nil {
		//验证失败
		log.Println("verify notify fail:%s", err)
		//w.Write([]byte(callbackMsg))
		return
	}

	trade_status := r.FormValue("trade_status")
	out_trade_no := r.FormValue("out_trade_no")
	buyer_email := r.FormValue("buyer_email")
	subject := r.FormValue("subject")

	log.Println("trade_status is : %v ", trade_status)
	log.Println("out_trade_no is : %v ", out_trade_no)
	log.Println("buyer_email is : %v ", buyer_email)
	log.Println("subject is : %v ", subject)

	var total_fee float64
	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	//判断该笔订单是否在商户网站中已经做过处理
	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
	//如果有做过处理，不执行商户的业务程序

	//注意：
	//该种交易状态只在一种情况下出现——开通了高级即时到账，买家付款成功后。

	if trade_status == "TRADE_SUCCESS" {

		log.Println("success")
		Callback()

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
	//	echo "success";		//请不要修改或删除
	callbackMsg = "success"
	//w.Write([]byte(callbackMsg))
	return
}

/*
test
*/
func AlipayWebTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Test...")
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	strUrl := fmt.Sprintf("location.href='%s'", "http://www.baidu.com")
	callbackMsg := "<html>"
	callbackMsg += "<head>"
	callbackMsg += "<script type='text/javascript'>"
	callbackMsg += strUrl
	//callbackMsg += "location.href='http://www.baidu.com'"
	//callbackMsg += "javascript:history.back(-1);"
	//callbackMsg += "location.href=document.referrer;"
	//callbackMsg += "javascript:history.go(-1);"
	callbackMsg += "</script>"
	callbackMsg += "</head>"
	callbackMsg += "<body>"

	callbackMsg += "</body>"
	callbackMsg += "</html>"
	w.Write([]byte(callbackMsg))
	//fmt.Fprint(w, callbackMsg)
}

//支付宝 同步通知处理
func AlipayWebReturn(w http.ResponseWriter, r *http.Request) {
	log.Println("AlipayWebReturn Begin")

	//var callbackMsg = "fail"
	var callbackMsg string
	defer func() {
		log.Println("AlipayWebReturn End")
		log.Println("callbackMsg to alipay in return_url: ", callbackMsg)
		//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		//fmt.Fprint(w, callbackMsg)
	}()

	//	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	r.PostForm = nil
	r.ParseForm()

	log.Println("==========================================================")
	log.Println("AlipayWebReturn Request :%v", r)
	log.Println("==========================================================")

	if err := alipay.VerifyWebNotify(r, alipay.AWebConfig); err != nil {
		//验证失败
		log.Println("verify notify fail")
		//callbackMsg = "verify notify fail"
		return
	}

	trade_status := r.FormValue("trade_status")
	// out_trade_no := r.FormValue("out_trade_no")
	// buyer_email := r.FormValue("buyer_email")
	// subject := r.FormValue("subject")
	// log.Println("buyer_email is : %v ", buyer_email)
	// log.Println("subject is : %v ", subject)
	// log.Println("trade_status is : %v ", trade_status)
	// log.Println("out_trade_no is : %v ", out_trade_no)

	// var total_fee float64
	// fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	if trade_status == "TRADE_SUCCESS" {

		//todo : deal the order
		log.Println("TRADE_SUCCESS in return_url")
	}

	if trade_status == "TRADE_FINISHED" {

	}

	//callbackMsg = ""
	callbackMsg = "<html>"
	callbackMsg += "<head>"
	callbackMsg += "<script type='text/javascript'>"
	//callbackMsg += "location.href='http://wwww.baidu.com'"
	//callbackMsg += "javascript:history.go(-5);"
	callbackMsg += "location.href='http://rcp.dev.jxzy.com/courseDetail.html?id=39940'"
	callbackMsg += "</script>"
	callbackMsg += "</head>"
	callbackMsg += "<body>"

	callbackMsg += "</body>"
	callbackMsg += "</html>"
	w.Write([]byte(callbackMsg))
	return
}
