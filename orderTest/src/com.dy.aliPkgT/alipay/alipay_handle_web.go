package alipay

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

/*
Remote call
*/
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

	htm = fmt.Sprintf("%s", `<html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head><body>`)
	htm += fmt.Sprintf(`<form name='alipaysubmit' action='%s_input_charset=%s' method='post'> `, alipayGatewayNew, alipayConfig.Input_charset)
	for _, kv := range p {
		htm += fmt.Sprintf(`<input type='hidden' name='%s' value='%s' />`, kv.K, kv.V)
	}
	//
	htm += fmt.Sprintf(`</form>`)
	htm += fmt.Sprintf(`<script>document.forms['alipaysubmit'].submit();</script>`)
	htm += fmt.Sprintf(`</body></html>`)
	return htm, nil
}

func AlipayWebRequestForm(alipayConfig *AlipayConfig, r *AlipayWebRequest, w io.Writer) error {
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

	fmt.Fprintln(w, `<html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head><body>`)
	fmt.Fprintf(w, `<form name='alipaysubmit' action='%s_input_charset=%s' method='post'> `, alipayGatewayNew, alipayConfig.Input_charset)
	for _, kv := range p {
		fmt.Fprintf(w, `<input type='hidden' name='%s' value='%s' />`, kv.K, kv.V)
	}
	//
	fmt.Fprintf(w, `</form>`)
	fmt.Fprintln(w, `<script>document.forms['alipaysubmit'].submit();</script>`)
	fmt.Fprintln(w, `</body></html>`)
	return nil
}

// func simulateAlipay(ty []string) error {
// 	var v url.Values = make(map[string][]string)
// 	url_ := fmt.Sprintf("http://192.168.2.174:8888/AliPayRequest")
// 	v.Add("_input_charset", ty[0])
// 	v.Add("notify_url", ty[1])
// 	v.Add("out_trade_no", ty[2])
// 	v.Add("partner", ty[3])
// 	v.Add("payment_type", ty[4])
// 	v.Add("return_url", ty[5])
// 	v.Add("seller_email", ty[6])
// 	v.Add("service", ty[7])
// 	v.Add("subject", ty[8])
// 	v.Add("total_fee", ty[9])
// 	v.Add("sign", ty[10])
// 	v.Add("sign_type", ty[11])
// 	req, _ := http.NewRequest("POST", url_, strings.NewReader(v.Encode()))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
// 	if resp, err := http.DefaultClient.Do(req); err != nil {
// 		log.E("err:%v", err)
// 		return err
// 	} else {
// 		b, _ := ioutil.ReadAll(resp.Body)
// 		parse := make(map[string]interface{})
// 		json.Unmarshal(b, &parse)
// 		log.D("SendNotify resp:%v,resp.StatusCode:%d", string(b), resp.StatusCode)
// 		if fmt.Sprintf("%v", parse["code"]) != "0" {
// 			log.E("error :%v", parse["msg"])
// 			return fmt.Errorf("%v", parse["msg"])
// 		} else {
// 			return nil
// 		}
// 	}
// }

/**
 * 验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWebReturn(r *http.Request, alipayConfig *AlipayConfig) error {
	log.Println("VerifyWebReturn begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	for k := range r.Form {
		v := r.PostForm.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		}
		*p = append(*p, Kvpair{k, v})
	}
	//除去待签名参数数组中的空值和签名参数
	paraFilter(p)

	//对待签名参数数组排序
	argSort(p)

	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(p)

	log.Println("VerifyWebReturn prestr is : %v ", prestr)
	log.Println("VerifyWebReturn sign is : %v  , sign_type is %v", sign, sign_type)

	switch sign_type {
	case "MD5":
		if md5Sign(prestr, alipayConfig.Key) != sign {
			return fmt.Errorf("sign invalid")
		}
		break
	default:
		return fmt.Errorf("no right sign_type")
	}

	log.Println("VerifyWebReturn success")

	notify_id := r.FormValue("notify_id")
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）(1分钟认证)
	responseTxt, err := getResponse(notify_id, alipayConfig)
	if err != nil {
		return err
	}
	log.Println("VerifyWebReturn responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.Println("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.Println("VerifyWebReturn responseTxt verify success ")
	return nil

}

/**
 * 针对notify_url验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWebNotify(r *http.Request, alipayConfig *AlipayConfig) error {
	log.Println("VerifyWebNotify begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	for k := range r.Form {
		v := r.Form.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		}
		*p = append(*p, Kvpair{k, v})
	}
	//除去待签名参数数组中的空值和签名参数
	paraFilter(p)

	//对待签名参数数组排序
	argSort(p)

	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(p)

	log.Println("VerifyWebNotify prestr is : %v ", prestr)
	log.Println("VerifyWebNotify sign is : %v  , sign_type is %v", sign, sign_type)

	switch sign_type {
	case "MD5":
		if md5Sign(prestr, alipayConfig.Key) != sign {
			return fmt.Errorf("sign invalid")
		}
		break
	default:
		return fmt.Errorf("no right sign_type")
	}

	log.Println("VerifyWebNotify success")

	notify_id := r.FormValue("notify_id")
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）(1分钟认证)
	responseTxt, err := getResponse(notify_id, alipayConfig)
	if err != nil {
		return err
	}
	log.Println("VerifyWebNotify responseTxt is: ", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.Println("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.Println("VerifyWebNotify responseTxt verify success ")
	return nil
}

/*
http://192.168.2.174:8888
/alipay-web-return
?body=test+web+Chinese&
buyer_email=327468120%40qq.com&
buyer_id=2088502994384781&
exterface=create_direct_pay_by_user&
is_success=T&
notify_id=RqPnCoPT3K9%252Fvwbh3InVZ3r3Q8XImQft4yHdn52q1w4C8YioyZQ8YGGBWRFLjAeogS0l&
notify_time=2015-12-04+12%3A16%3A11&notify_type=trade_status_sync&
out_trade_no=2015120412150153694&payment_type=1&
seller_email=itdayang%40gmail.com&
seller_id=2088501949844011&
subject=迟到扣200&
total_fee=0.01&
trade_no=2015120421001004780215834730&
trade_status=TRADE_SUCCESS&
sign=c0a5dbab3a9e707367fc8aad962a5288&
sign_type=MD5
*/
func GenAndNotify(r *http.Request, alipayConfig *AlipayConfig) (string, url.Values, error) {

	p := &Kvpairs{}
	nv := url.Values{}
	sign := ""
	sign_type := "MD5"
	for k := range r.Form {
		v := r.PostForm.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		case "return_url":
			continue
		case "notify_url":
			continue
		}
		*p = append(*p, Kvpair{k, v})
		nv.Set(k, v)
	}
	*p = append(*p, Kvpair{`exterface`, "create_direct_pay_by_user"})
	*p = append(*p, Kvpair{`is_success`, "T"})
	// *p = append(*p, Kvpair{`notify_id`, "RqPnCoPT3K9%252Fvwbh3InVZ3r3Q8XImQft4yHdn52q1w4C8YioyZQ8YGGBWRFLjAeogS0l"})
	// *p = append(*p, Kvpair{`notify_time`, "2015-12-14+12%3A16%3A11"})
	*p = append(*p, Kvpair{`notify_id`, "RqPnCoPT3K9252Fvwbh"})
	*p = append(*p, Kvpair{`notify_time`, "2015-12-14-121611"})
	*p = append(*p, Kvpair{`trade_status`, "TRADE_SUCCESS"})
	nv.Set(`exterface`, "create_direct_pay_by_user")
	nv.Set(`is_success`, "T")
	nv.Set(`notify_id`, "RqPnCoPT3K9252Fvwbh")
	nv.Set(`notify_time`, "2015-12-14-121611")
	nv.Set(`trade_status`, "TRADE_SUCCESS")

	paraFilter(p)
	argSort(p)
	sign = md5Sign(createLinkStringNoUrl(p), alipayConfig.Key)

	fmt.Printf("p:%s\n", p)
	fmt.Printf("nv:%s\n", nv)
	*p = append(*p, Kvpair{`sign`, sign})
	*p = append(*p, Kvpair{`sign_type`, sign_type})
	nv.Set(`sign_type`, "MD5")
	nv.Set(`sign`, sign)

	str := createLinkStringNoUrl(p)
	strUrl := fmt.Sprintf("%s%s", alipayConfig.Return_url, str)
	fmt.Printf("strUrl:%s\n", strUrl)
	//同步通知
	resp, err := http.Get(strUrl)
	if err != nil {
		fmt.Printf("return_url get error:%s\n", err.Error())
		return "", nv, err
	}
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return "", nv, err
	}
	fmt.Printf("got:%s\n", string(got))
	return string(got), nv, nil

}

func GenTestData() (map[string]string, error) {

	p := &Kvpairs{}
	nv := map[string]string{}

	sign := ""
	sign_type := "MD5"

	*p = append(*p, Kvpair{`total_fee`, "0.01"})
	*p = append(*p, Kvpair{`_input_charset`, "utf-8"})
	*p = append(*p, Kvpair{`_output_charset`, "utf-8"})
	*p = append(*p, Kvpair{`notify_url`, "http://order.dev.jxzy.com"})
	*p = append(*p, Kvpair{`out_trade_no`, "2015121518304888932"})
	*p = append(*p, Kvpair{`partner`, "2088501949844011"})
	*p = append(*p, Kvpair{`payment_type`, "1"})
	*p = append(*p, Kvpair{`return_url`, "http://order.dev.jxzy.com/alipay-web-return"})
	*p = append(*p, Kvpair{`seller_email`, "itdayang@gmail.com"})
	*p = append(*p, Kvpair{`service`, "create_direct_pay_by_user"})
	*p = append(*p, Kvpair{`subject`, "测试课程"})

	paraFilter(p)
	argSort(p)
	sign = md5Sign(createLinkStringNoUrl(p), "viz4safb1zazb5bqeraujlg79agfcj02")

	fmt.Printf("p:%s\n", p)
	fmt.Printf("nv:%s\n", nv)
	*p = append(*p, Kvpair{`sign`, sign})
	*p = append(*p, Kvpair{`sign_type`, sign_type})

	nv[`total_fee`] = "0.01"
	nv[`_input_charset`] = "utf-8"
	nv[`_output_charset`] = "utf-8"
	nv[`notify_url`] = "http://order.dev.jxzy.com"
	nv[`out_trade_no`] = "2015121518304888932"
	nv[`partner`] = "2088501949844011"
	nv[`payment_type`] = "1"
	nv[`return_url`] = "http://order.dev.jxzy.com/alipay-web-return"
	nv[`seller_email`] = "itdayang@gmail.com"
	nv[`service`] = "create_direct_pay_by_user"
	nv[`subject`] = "测试课程"
	nv[`sign_type`] = "MD5"
	nv[`sign`] = sign

	return nv, nil

}
