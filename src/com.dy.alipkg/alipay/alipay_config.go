package alipay

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

type AlipayConfig struct {
	//↓↓↓↓↓↓↓↓↓↓请在这里配置您的基本信息↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
	//合作身份者id，以2088开头的16位纯数字
	Partner string
	//安全检验码，以数字和字母组成的32位字符
	Key string
	//商户的私钥（后缀是.pen）文件相对路径
	//如果签名方式设置为“0001”时，请设置该参数
	Private_key_path []byte
	//支付宝公钥（后缀是.pen）文件相对路径
	//如果签名方式设置为“0001”时，请设置该参数
	Ali_public_key_path []byte
	//↑↑↑↑↑↑↑↑↑↑请在这里配置您的基本信息↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
	//签名方式 不需修改
	Sign_type string
	//字符编码格式 目前支持 gbk 或 utf-8
	Input_charset string
	//ca证书路径地址，用于curl中ssl校验
	//请保证cacert.pem文件在当前文件夹目录中
	Cacert string
	//访问模式,根据自己的服务器是否支持ssl访问，若支持请选择https；若不支持请选择http
	Transport string
	//商户私钥
	Private_key *rsa.PrivateKey
	//支付宝公钥 （验签）
	Public_key *rsa.PublicKey
	//服务类型（移动支付）
	Service string
	//卖家支付宝账号
	Seller_id string
	//异步通知URL
	Notify_url string
	//
	Return_url string
	//支付类型  1
	Payment_type string
	//显示订单消息页面
	Show_order_url string
	//操作中断返回地址
	Wap_merchant_url string
	//页面跳转同步通知页面路径
	Wap_callback_url string
	//Wap_Service
	Wap_Service string
}

var AWebConfig = &AlipayConfig{
	Partner:   "2088501949844011",
	Key:       "viz4safb1zazb5bqeraujlg79agfcj20",
	Sign_type: "MD5",
	Private_key_path: []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDAymEQFouLpQ7a8dhynKVLDD9T0yxjY9LeTQk0Y+97rf9sjk22
OdeNvVpTnEMd1GHjQL2GU/YFQOMG2gIbkLMUashQGkNGGUIQa/owF4Us3vIhaENw
gEWg3ybpuYwO2QVNXD3CpWwTuQ/KXllKq/n5M1vZgB7vzW5SW2Ll7tcwdQIDAQAB
AoGAA0Gn4BCxbLtuA86//DefJyRe0XiosEkYX77nX+YsPmeS/+9rap+Rf4hqRS1H
iPJ0cNVAHku+xRnye8Qk8vg/vhoE40HZqIp8cUAb2x/+V5hPYX8VT8Joby8y37kk
JwaSR4BbNq5r2h7OTmLAIyLBOCWfPIUxm092NS9N+Hh9gCECQQDvg3b0u3299yEo
tx9yz1rJo5DPQw/j/9lHI4/dlwmU7qgoitDD0CoHm57MDoP0rfodmqRQuu39WLvQ
RSAVsKjJAkEAzg+YdKOskpEQ5gqVpnPG1QX+zJ2q6IrJwoz3Vz7ApZuZuTMNQP6I
pRk+DwyYQjMLJZ9pWXMnQqrUZsXymvMMTQJBAN0Ed3UGvtXZgqPLuB198s03Pp3D
yJogsiPPUSauBY2FXW9sCZFoOiCjHjo2+tWC5dcDqLVGie7LPyTFZQIe7uECQQCh
9341CEXoOqBWwTEZ0d0SqhC6Up/5AsZN5rGv4QJZ6bFt8ePUF33ej7XI7HRtbGOx
odiT18CfVVl4otTr2GjtAkBlAFytEjIbtdpYDE1F7ANFoVFGy6jeCCs7re1D2BTp
pIGJlxkkxC/Loi2ZnztwbY/DbOCQpN9dvSEiwzsOQyWL
-----END RSA PRIVATE KEY-----
	`),
	Ali_public_key_path: []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnxj/9qwVfgoUh/y2W89L6BkRA
FljhNhgPdyPuBV64bfQNN1PjbCzkIM6qRdKBoLPXmKKMiFYnkd6rAoprih3/PrQE
B/VsW8OoM8fxn67UDYuyBTqA23MML9q1+ilIZwBC2AQ2UBVOrFXfFl75p6/B5Ksi
NG9zpgmLCUYuLkxpLQIDAQAB
-----END PUBLIC KEY-----
`),
	Input_charset: "utf-8",
	//Input_charset:  "GBK",
	Cacert:         "Cacert",
	Transport:      "http",
	Service:        "create_direct_pay_by_user", //即时到帐
	Seller_id:      "itdayang@gmail.com",
	Notify_url:     "/api/pub/alipayWeb/notify",
	Return_url:     "/api/pub/alipayWeb/return",
	Payment_type:   "1", //即时到帐
	Show_order_url: "/paymentStatus.html",
}

var AMobileConfig = &AlipayConfig{
	Partner:   "2088501949844011",
	Key:       "viz4safb1zazb5bqeraujlg79agfcj02",
	Sign_type: "RSA",
	Private_key_path: []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDAymEQFouLpQ7a8dhynKVLDD9T0yxjY9LeTQk0Y+97rf9sjk22
OdeNvVpTnEMd1GHjQL2GU/YFQOMG2gIbkLMUashQGkNGGUIQa/owF4Us3vIhaENw
gEWg3ybpuYwO2QVNXD3CpWwTuQ/KXllKq/n5M1vZgB7vzW5SW2Ll7tcwdQIDAQAB
AoGAA0Gn4BCxbLtuA86//DefJyRe0XiosEkYX77nX+YsPmeS/+9rap+Rf4hqRS1H
iPJ0cNVAHku+xRnye8Qk8vg/vhoE40HZqIp8cUAb2x/+V5hPYX8VT8Joby8y37kk
JwaSR4BbNq5r2h7OTmLAIyLBOCWfPIUxm092NS9N+Hh9gCECQQDvg3b0u3299yEo
tx9yz1rJo5DPQw/j/9lHI4/dlwmU7qgoitDD0CoHm57MDoP0rfodmqRQuu39WLvQ
RSAVsKjJAkEAzg+YdKOskpEQ5gqVpnPG1QX+zJ2q6IrJwoz3Vz7ApZuZuTMNQP6I
pRk+DwyYQjMLJZ9pWXMnQqrUZsXymvMMTQJBAN0Ed3UGvtXZgqPLuB198s03Pp3D
yJogsiPPUSauBY2FXW9sCZFoOiCjHjo2+tWC5dcDqLVGie7LPyTFZQIe7uECQQCh
9341CEXoOqBWwTEZ0d0SqhC6Up/5AsZN5rGv4QJZ6bFt8ePUF33ej7XI7HRtbGOx
odiT18CfVVl4otTr2GjtAkBlAFytEjIbtdpYDE1F7ANFoVFGy6jeCCs7re1D2BTp
pIGJlxkkxC/Loi2ZnztwbY/DbOCQpN9dvSEiwzsOQyWL
-----END RSA PRIVATE KEY-----
`),
	Ali_public_key_path: []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnxj/9qwVfgoUh/y2W89L6BkRA
FljhNhgPdyPuBV64bfQNN1PjbCzkIM6qRdKBoLPXmKKMiFYnkd6rAoprih3/PrQE
B/VsW8OoM8fxn67UDYuyBTqA23MML9q1+ilIZwBC2AQ2UBVOrFXfFl75p6/B5Ksi
NG9zpgmLCUYuLkxpLQIDAQAB
-----END PUBLIC KEY-----
`),
	Input_charset: "UTF-8",
	Cacert:        "Cacert",
	Transport:     "http",
	Service:       "mobile.securitypay.pay",
	Seller_id:     "itdayang@gmail.com",
	Notify_url:    "/alipayMobile/notify",
	Payment_type:  "1",
}

var AWapConfig = &AlipayConfig{
	Partner:   "2088501949844011",
	Key:       "viz4safb1zazb5bqeraujlg79agfcj20",
	Sign_type: "MD5",
	Private_key_path: []byte(`
-----BEGIN RSA PRIVATE KEY-----
xxxxxxxxxxx
-----END RSA PRIVATE KEY-----
`),
	Ali_public_key_path: []byte(`
-----BEGIN PUBLIC KEY-----
xxxxxxxxxxx
-----END PUBLIC KEY-----
`),
	Input_charset: "utf-8",
	//	Cacert:        "Cacert",
	Transport:   "http",
	Service:     "alipay.wap.auth.authAndExecute",
	Wap_Service: "alipay.wap.trade.create.direct",
	Seller_id:   "itdayang@gmail.com",
	Notify_url:  "/api/pub/alipayWap/notify",
	//	Payment_type:  "1",
	Wap_merchant_url: "/api/pub/alipayWap/merchant",
	//页面跳转同步通知页面路径
	Wap_callback_url: "/api/pub/alipayWap/callback",
	Show_order_url:   "/paymentStatus.html",
}

func Init(alipayConfig *AlipayConfig) error {

	log.Println("init rsakeys begin")

	block, _ := pem.Decode(alipayConfig.Private_key_path)
	if block == nil {
		log.Println("rsaSign private_key error")
		return fmt.Errorf("rsaSign pem.Decode error")
	}
	var err error
	//把RSA私钥转换成PKCS8格式
	alipayConfig.Private_key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("rsaSign ParsePKIXPublicKey error : %v\n", err)
		return err
	}

	block, _ = pem.Decode(alipayConfig.Ali_public_key_path)
	if block == nil {
		log.Println("public key error")
		return err
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println("rsaSign ParsePKIXPublicKey error : %v\n", err)
		return err
	}

	alipayConfig.Public_key = pubInterface.(*rsa.PublicKey)
	log.Println("init rsakeys success ")
	return err
}
