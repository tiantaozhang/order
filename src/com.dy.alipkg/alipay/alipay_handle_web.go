package alipay

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	fmt.Printf("VerifyWebReturn sign is : %v , sign_type is %v\n", sign, sign_type)

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
	fmt.Printf("r in WebNotify:%s\n", r)
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

	log.Println("VerifyWebNotify prestr is : ", prestr)
	fmt.Printf("VerifyWebNotify sign is : %s, sign_type is %v\n", sign, sign_type)

	strSign := md5Sign(prestr, alipayConfig.Key)
	fmt.Printf("MD5:%s\n", strSign)

	switch sign_type {
	case "MD5":
		if strSign != sign {
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
	log.Println("VerifyWebNotify responseTxt is: %v", responseTxt)
	if err != nil {
		fmt.Println("responseTxt err:", err.Error())
		return err
	}

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.Println("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.Println("VerifyWebNotify responseTxt verify success ")
	return nil
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
