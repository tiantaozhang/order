package orderalipay

import (
	"com.dy.alipkg/alipay"
	"com.dy.order/conf"
	"fmt"
)

var private_key = `
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`

var public_key = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`

var partner = "xxxxxxxxx"

var key = "xxxxxxxxx"

var seller = "xxxxxxxxx@gmail.com"

func InitAlipayConfig() {

	host := conf.Order_host()
	//host := "14.23.162.170:12000"
	//host := "192.168.10.118:8888"

	// alipay.AWebConfig = &alipay.AlipayConfig{
	// 	Partner: partner,
	// 	Key:     key,
	// 	//	Sign_type:           "MD5",
	// 	Sign_type: "RSA",
	// 	//	Private_key_path:    []byte(private_key),
	// 	//	Ali_public_key_path: []byte(public_key),
	// 	//Input_charset:       "utf-8",
	// 	Input_charset:  "GBK",
	// 	Cacert:         "Cacert",
	// 	Transport:      "http",
	// 	Service:        "create_direct_pay_by_user",
	// 	Seller_id:      seller,
	// 	Payment_type:   "1",
	// 	Show_order_url: "/paymentStatus.html",
	// }

	// alipay.AMobileConfig = &alipay.AlipayConfig{
	// 	Partner:             partner,
	// 	Key:                 key,
	// 	Sign_type:           "RSA",
	// 	Private_key_path:    []byte(private_key),
	// 	Ali_public_key_path: []byte(public_key),
	// 	Input_charset:       "UTF-8",
	// 	Cacert:              "Cacert",
	// 	Transport:           "http",
	// 	Service:             "mobile.securitypay.pay",
	// 	Seller_id:           seller,
	// 	Payment_type:        "1",
	// }

	// alipay.AWapConfig = &alipay.AlipayConfig{
	// 	Partner:             partner,
	// 	Key:                 key,
	// 	Sign_type:           "MD5",
	// 	Private_key_path:    []byte(private_key),
	// 	Ali_public_key_path: []byte(public_key),
	// 	Input_charset:       "utf-8",
	// 	Transport:           "http",
	// 	Service:             "alipay.wap.auth.authAndExecute",
	// 	Wap_Service:         "alipay.wap.trade.create.direct",
	// 	Seller_id:           seller,
	// 	Show_order_url:      "/paymentStatus.html",
	// }

	alipay.AMobileConfig.Notify_url = fmt.Sprintf("http://%v/alipay-mobile-notify", host)
	alipay.AMobileConfig.Return_url = fmt.Sprintf("http://%v/alipay-mobile-return", host)
	//alipay.AMobileConfig.Return_url = fmt.Sprintf("http://%v/alipay-mobile-return", "192.168.2.174")
	err := alipay.Init(alipay.AMobileConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// alipay.AWapConfig.Notify_url = fmt.Sprintf("http://%v/alipay-wap-notify", host)
	// alipay.AWapConfig.Wap_merchant_url = fmt.Sprintf("http://%v/merchant", host)
	// alipay.AWapConfig.Wap_callback_url = fmt.Sprintf("http://%v/alipay-wap-callback", host)
	// alipay.AWapConfig.Show_order_url = fmt.Sprintf("http://%v/orderDetail", host)

	// alipay.Init(alipay.AWapConfig)

	alipay.AWebConfig.Notify_url = fmt.Sprintf("http://%v/alipay-web-notify", host)
	alipay.AWebConfig.Return_url = fmt.Sprintf("http://%v/alipay-web-return", host)
	alipay.AWebConfig.Show_order_url = fmt.Sprintf("http://%v/orderDetail", host)

	alipay.Init(alipay.AWebConfig)

}
