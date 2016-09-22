package orderalipay

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	//"net/http"
	"com.dy.order/common"
	"com.dy.order/conf"
	"com.dy.order/orderModel"
	"com.dy.order/orderwxpay"
	"com.dy.order/testconf"
	"com.dy.wxpkg/wechatPay"
	"database/sql"
	"encoding/json"
	"github.com/Centny/DEM"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	uap_cf "org.cny.uap/conf"
	"org.cny.uap/uap"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

// var strono string
var gono string = NewOrderNo()
var gonoPaid string = NewOrderNo()

var t_conn *sql.DB
var t_db *sql.DB

func insertData() {
	time.Now()
	sql_ := `INSERT INTO ods_order ( ono, buyer, seller, total_price, type, status, time, return_url, expand, wno, add1, add2)VALUES ( ?, 267250, 438982,0.01, 'N', 'NOT_PAY', ?, 'http://rcp.dev.jxzy.com/questionPoolDetailNew.html?id=40021&eid=54891', 'id=40021&token=b4561e0e8c3e185e9ef858cc54cad5f1-9933ec14-6590-4548-9613-0cda727bfbe4', NULL, NULL, NULL)`
	_, err := t_conn.Exec(sql_, gono, time.Now())
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	sql_ = `INSERT INTO ods_order ( ono, buyer, seller, total_price, type, status, time, return_url, expand, wno, add1, add2)VALUES ( ?, 267250, 438982,0.01, 'N', 'PAID', ?, 'http://rcp.dev.jxzy.com/questionPoolDetailNew.html?id=40021&eid=54891', 'id=40021&token=b4561e0e8c3e185e9ef858cc54cad5f1-9933ec14-6590-4548-9613-0cda727bfbe4', NULL, NULL, NULL)`
	_, err = t_conn.Exec(sql_, gonoPaid, time.Now())
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	//fmt.Println("sql_ ods_order_item")
	sql_ = `INSERT INTO ods_order_item ( ono, p_name, p_id, p_type, p_img, p_from,status, time, add1, add2)VALUES ( ?, 'testSyn', 267250, 10, NULL, 'TEST', 'N',?,NULL, NULL)`
	_, err = t_conn.Exec(sql_, gono, time.Now())
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	sql_ = `INSERT INTO uap_attr( a_key,oid,owner,ns,integral,add1, add2)VALUES ('USER_INTEGRAL',267250,'USER',DEFAULT,1000,NULL, NULL)`
	_, err = t_conn.Exec(sql_)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	sql_ = `insert into ods_record(name,type,money,uid,pay_type,target_id,ono,status,add1,add2) values('寻龙诀','INCOME',0.01,267250,'ALIPAY',267250,?,'NOT_PAY',NULL,NULL)`
	_, err = t_conn.Exec(sql_, gono)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	_, err = t_conn.Exec(`insert into ods_order_env(akey,aval,type,status) values('TEST','/usr/purchase-course?','PAID_CB','N')`)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	// t_conn.Exec("INSERT INTO `ods_order_item` (`tid`, `ono`, `oid`, `p_name`, `p_id`, `p_type`, `p_img`, `p_count`, `p_from`, `notified`, `price`, `type`, `status`, `time`, `add1`, `add2`)VALUES(11415, '2016012618114979877', 0, '第二个付费课程哈哈哈', 40025, '10', 'http://u.dev.jxzy.com/fLRlKA==', 1, 'RCP', 0, 0.01, 'N', 'N', '2016-01-26 18:11:49', NULL, NULL)")
	// t_conn.Exec("INSERT INTO `ods_order` (`tid`, `ono`, `buyer`, `seller`, `total_price`, `type`, `status`, `time`, `return_url`, `expand`, `wno`, `add1`, `add2`)VALUES(9516, '2016012710275443433', 453177, 267250, 0.01, 'N', 'PAID', '2016-01-27 10:43:45', 'http://rcp.dev.jxzy.com/pay-result.html?id=40025&eid=0', 'id=40025&token=8dad86e559c60901330b126861a447ae-3a2aa1fe-778d-4477-8da7-d969b73de7e7', NULL, NULL, NULL)")
}

func deleteData() {
	sql_ := `delete from ods_order`
	_, err := t_conn.Exec(sql_)
	if err != nil {
		fmt.Println(err.Error())
	}
	t_conn.Exec("delete from ods_order_item")
	t_conn.Exec("delete from ods_record")
	t_conn.Exec("delete from ods_order_env where akey='TEST'")
	t_conn.Exec("delete from uap_attr where oid=267250")

}

//var db *sql.DB
func InitDB() error {
	// err := common.Init("mysql", "cny:123@tcp(192.168.2.57:3306)/ucs?charset=utf8&loc=Local")
	// if err != nil {
	// 	fmt.Println("init db err:", err)
	// 	return err
	// }
	// err = common.Init("mysql", "cny:123@tcp(192.168.2.57:3306)/orderv2?charset=utf8&loc=Local")

	// if err != nil {
	// 	fmt.Println("init db err:", err)
	// }

	cfg := "../../../conf/order.properties"
	err := conf.Cfg.InitWithFilePath(cfg)
	uap_cf.Cfg = conf.Cfg
	//conf.Cfg.Print()
	if err != nil {
		fmt.Println(cfg)
		//panic(err)
		return err
	}
	DEM.ShowLog(true)
	DEM.G_Dn = "mysql"
	g_dsn := conf.ORDER_DB_CONN()
	fmt.Println("DEM.G_Dsn:", g_dsn)
	//orderv2 ---> orderv2_test
	s := strings.Split(g_dsn, "?")
	if len(s) < 1 {
		panic("数据库地址格式有误")
	}
	DEM.G_Dsn = s[0] + "_test?" + s[1]
	fmt.Println("DEM.G_Dsn:", DEM.G_Dsn)
	fmt.Println(common.Init("DEM", DEM.G_Dsn))
	t_conn = DEM.LAST
	t_db = common.DbConn()

	uap.InitDb(common.DbConn)
	insertData()

	return nil
}

func AssemblePara() AlipayRemoteReqStruct {

	m := AlipayRemoteReqStruct{
		Ono:        "",
		Buyer:      267250,
		Seller:     438982,
		Subject:    "寻龙诀",
		TotalFee:   0.01,
		Body:       "迟到扣200",
		Type:       "N",
		Status:     "NOT_PAY",
		Return_url: "http://rcp.dev.jxzy.com/courseDetail.html?id=40040",
		Expand:     "id=40040&token=4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35",
		Integral:   0,
	}
	//orderI := []map[string]interface{}{}
	for i := 1; i < 2; i++ {
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

	return m
}
func AssemErrPara() AlipayRemoteReqStruct {
	em := AlipayRemoteReqStruct{
		Ono:        "",
		Buyer:      267250,
		Seller:     438982,
		Subject:    "寻龙诀",
		TotalFee:   0.01,
		Body:       "迟到扣200",
		Type:       "N",
		Status:     "NOT_PAY",
		Return_url: "http://rcp.dev.jxzy.com/courseDetail.html?id=40040",
		Expand:     "id=40040&token=4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35",
		Integral:   0,
	}
	return em
}

func TestWebPay(t *testing.T) {
	//m := map[string]interface{}{"ono": "", "buyer": 54321, "seller": 44444, "subject": "天地一号", "totalFee": 0.01, "body": "迟到扣200", "type": "N", "status": "NOT_PAY"}
	//orderI := []Item{}
	//	jsstr := `{"ono": "", "buyer": 54321, "seller": 44444, "subject": "天地一号", "totalFee": 0.01, "body": "迟到扣200", "type": "N", "status": "NOT_PAY","item":[{"ono":"","oid":i,"p_name":str,"p_id":i,"p_type":"","p_count":1,"p_from":"RCP","notify":0,"price":0.01,"type":"N","status":"N",}],"refund":[""]}`
	if err := InitDB(); err != nil {
		t.Error(err.Error())
	}

	InitAlipayConfig()
	m := AssemblePara()

	bytes, _ := json.Marshal(m)
	//fmt.Printf("json:m,%s\n", bytes)

	convey.Convey("json err", t, func() {
		_, err := AlipayRemoteRequest(string(""))
		convey.So(err, convey.ShouldNotBeNil)
		_, err = GetRsaSignJson(string(""))
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("return url nil", t, func() {
		m.Return_url = ""
		bytes1, _ := json.Marshal(m)
		_, err := AlipayRemoteRequest(string(bytes1))
		convey.So(err, convey.ShouldNotBeNil)
		_, err = GetRsaSignJson(string(bytes1))
		convey.So(err, convey.ShouldNotBeNil)
	})
	m.Return_url = "http://rcp.dev.jxzy.com/courseDetail.html?id=40040"
	convey.Convey("dem sync uap err", t, func() {
		DEM.Evb.SetErrs(DEM.OPEN_ERR | DEM.BEGIN_ERR)
		strhtm, err := AlipayRemoteRequest(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(strhtm, convey.ShouldBeBlank)
		strhtm, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(strhtm, convey.ShouldBeBlank)
		DEM.Evb.ClsQErr()
	})
	DEM.Evb.SetErrs(0)
	// convey.Convey("dem checkExist err", t, func() {
	// 	DEM.Evb.SetErrs(0)
	// 	DEM.Evb.AddQErr3(".*select.*count.*")
	// 	_, err := AlipayRemoteRequest(string(bytes))
	// 	convey.So(err, convey.ShouldNotBeNil)

	// 	DEM.Evb.ClsQErr()
	// })

	convey.Convey("RemoteRequestTypeN", t, func() {

		strhtm, err := AlipayRemoteRequest(string(bytes))
		// strono = GetCurrentOno()

		if err != nil {
			fmt.Println(err)
		}
		convey.So(err, convey.ShouldBeNil)
		convey.So(strhtm, convey.ShouldNotBeBlank)
		//fmt.Printf("htm:%s\n", strhtm)

		strhtm, err = GetRsaSignJson(string(bytes))
		if err != nil {
			fmt.Println(err)
		}
		convey.So(err, convey.ShouldBeNil)
		convey.So(strhtm, convey.ShouldNotBeBlank)
	})
	convey.Convey("RemoteRequestTotalFee--err", t, func() {
		//m["totalFee"] = -0.01
		m.TotalFee = -0.01
		bytes, _ = json.Marshal(m)
		strhtm, err := AlipayRemoteRequest(string(bytes))

		convey.So(err, convey.ShouldNotBeNil)
		convey.So(strhtm, convey.ShouldBeBlank)

		strhtm, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(strhtm, convey.ShouldBeBlank)
	})

	convey.Convey("RemoteRequestFormat--err", t, func() {
		//m["totalFee"] = "-0.01"
		m.TotalFee = 0
		bytes, _ = json.Marshal(m)
		strhtm, err := AlipayRemoteRequest(string(bytes))

		convey.So(err, convey.ShouldNotBeNil)
		convey.So(strhtm, convey.ShouldBeBlank)

		strhtm, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(strhtm, convey.ShouldBeBlank)

	})

	convey.Convey("RemoteRequest Integral and alipay", t, func() {
		//m["totalFee"] = "-0.01"
		m.TotalFee = 0.02
		m.Integral = 1
		bytes, _ = json.Marshal(m)
		strhtm, err := AlipayRemoteRequest(string(bytes))

		convey.So(err, convey.ShouldBeNil)
		convey.So(strhtm, convey.ShouldNotBeBlank)

		strhtm, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldBeNil)
		convey.So(strhtm, convey.ShouldNotBeBlank)
	})

	convey.Convey("RemoteRequest Integral total", t, func() {
		//m["totalFee"] = "-0.01"
		m.TotalFee = 0.01
		m.Integral = 1
		bytes, _ = json.Marshal(m)
		strhtm, err := AlipayRemoteRequest(string(bytes))

		convey.So(err, convey.ShouldBeNil)
		convey.So(strhtm, convey.ShouldNotBeBlank)

		strhtm, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldBeNil)
		convey.So(strhtm, convey.ShouldNotBeBlank)
	})

	convey.Convey("RemoteRequest alipay only", t, func() {
		//m["totalFee"] = "-0.01"
		m.TotalFee = 0.01
		m.Integral = 0
		bytes, _ = json.Marshal(m)
		strhtm, err := AlipayRemoteRequest(string(bytes))

		convey.So(err, convey.ShouldBeNil)
		convey.So(strhtm, convey.ShouldNotBeBlank)

		strhtm, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldBeNil)
		convey.So(strhtm, convey.ShouldNotBeBlank)
	})
	convey.Convey("orderitem err", t, func() {
		em := AssemErrPara()
		by, _ := json.Marshal(em)
		_, err := AlipayRemoteRequest(string(by))
		convey.So(err, convey.ShouldNotBeNil)
		_, err = GetRsaSignJson(string(by))
		convey.So(err, convey.ShouldNotBeNil)

	})
	// convey.Convey("dem sync uap err", t, func() {
	// 	DEM.Evb.SetErrs(DEM.TX_COMMIT_ERR)
	// 	_, err := AlipayRemoteRequest(string(bytes))
	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	_, err = GetRsaSignJson(string(bytes))
	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	DEM.Evb.ClsQErr()
	// })
	// DEM.Evb.SetErrs(0)
	convey.Convey("RemoteRequest Integral and alipay--err", t, func() {
		m.TotalFee = 0.02
		m.Integral = 3
		bytes, _ = json.Marshal(m)
		_, err := AlipayRemoteRequest(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)

		_, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
	})

	convey.Convey("RemoteRequest InsertOrderItem --err", t, func() {
		m.TotalFee = 0.01
		m.Integral = 0
		bytes, _ = json.Marshal(m)

		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert into ods_order_item.*")
		_, err := AlipayRemoteRequest(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		_, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("RemoteRequest InsertOrderItem --err", t, func() {

		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert into ods_record.*")
		_, err := AlipayRemoteRequest(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		_, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("RemoteRequest InsertOrderItem --err", t, func() {
		m.TotalFee = 0.02
		m.Integral = 1
		bytes, _ = json.Marshal(m)
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert into ods_record.*")
		_, err := AlipayRemoteRequest(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		_, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("RemoteRequest InsertOrderItem --err", t, func() {
		m.TotalFee = 0.01
		m.Integral = 0
		bytes, _ = json.Marshal(m)
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert into ods_order\\(ono,buyer,seller,total_price,type,status,return_url,expand\\) value.*")
		_, err := AlipayRemoteRequest(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		_, err = GetRsaSignJson(string(bytes))
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})

	// convey.Convey("RemoteRequest seller not in --err", t, func() {
	// 	//m["totalFee"] = "-0.01"
	// 	m.Return_url = "http://rcp.dev.jxzy.com/courseDetail.html?id=40040"
	// 	m.TotalFee = 0.01
	// 	m.Integral = 0
	// 	m.Buyer = 200353
	// 	m.Seller = 1298712343
	// 	bytes, _ = json.Marshal(m)
	// 	strhtm, err := AlipayRemoteRequest(string(bytes))

	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)

	// 	strhtm, err = GetRsaSignJson(string(bytes))
	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)
	// })
	// convey.Convey("RemoteRequest buyer not in --err", t, func() {
	// 	//m["totalFee"] = "-0.01"
	// 	m.Return_url = "http://rcp.dev.jxzy.com/courseDetail.html?id=40040"
	// 	m.TotalFee = 0.01
	// 	m.Integral = 0
	// 	m.Buyer = 0
	// 	m.Seller = 1298712343
	// 	bytes, _ = json.Marshal(m)
	// 	strhtm, err := AlipayRemoteRequest(string(bytes))

	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)

	// 	strhtm, err = GetRsaSignJson(string(bytes))
	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)
	// })
	// convey.Convey("RemoteRequest no integral but comsume in --err", t, func() {
	// 	//m["totalFee"] = "-0.01"
	// 	m.Seller = 438982
	// 	m.Buyer = 200353
	// 	m.Integral = 123
	// 	bytes, _ = json.Marshal(m)
	// 	strhtm, err := AlipayRemoteRequest(string(bytes))

	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)

	// 	strhtm, err = GetRsaSignJson(string(bytes))
	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)
	// })
	// convey.Convey("RemoteRequest integral not enought --err", t, func() {
	// 	//m["totalFee"] = "-0.01"
	// 	m.Seller = 438982
	// 	m.Buyer = 267250
	// 	m.Integral = 1233
	// 	bytes, _ = json.Marshal(m)
	// 	strhtm, err := AlipayRemoteRequest(string(bytes))

	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)

	// 	strhtm, err = GetRsaSignJson(string(bytes))
	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)
	// })
	// convey.Convey("RemoteRequest orderitem --err", t, func() {
	// 	//m["totalFee"] = "-0.01"
	// 	m.Seller = 438982
	// 	m.Buyer = 267250
	// 	m.Integral = 0
	// 	item := []orderModel.Item{}
	// 	m.OrderItem = item
	// 	bytes, _ = json.Marshal(m)
	// 	strhtm, err := AlipayRemoteRequest(string(bytes))

	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)

	// 	strhtm, err = GetRsaSignJson(string(bytes))
	// 	convey.So(err, convey.ShouldNotBeNil)
	// 	convey.So(strhtm, convey.ShouldBeBlank)
	// })

}

func TestConfirmWXPay(t *testing.T) {
	convey.Convey("confirmPay err", t, func() {
		u := &wechatPay.Unifiedorder{}
		_, err := confirmWXPay(way[2], u, gono)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("confirmPay err", t, func() {
		u := &wechatPay.Unifiedorder{}
		_, err := confirmWXPay(way[3], u, gono)
		convey.So(err, convey.ShouldNotBeNil)
	})
}

func TestConfirmOrderAliPay(t *testing.T) {

	if err := InitDB(); err != nil {
		fmt.Println("InitDB err:", err)
		t.Error(err.Error())
	}
	InitAlipayConfig()
	orderwxpay.InitWxConfig()
	var ono string
	var rjs struct {
		Code string `json:"code"`
		Data Msg    `json:"data"`
	}
	for j := 0; j < 4; j++ {
		m := AlipayRemoteReqStruct{
			Ono:        "",
			Buyer:      267250,
			Seller:     438982,
			Subject:    "testcorfirm",
			TotalFee:   0.01,
			Body:       "迟到扣200",
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

		bytes, _ := json.Marshal(m)
		if j < 3 {
			// if 0 == j {
			if _, err := AlipayRemoteRequest(string(bytes)); err != nil {

				t.Error("request err:", err.Error())
			}
			ono = GetCurrentOno()
			// }
		} else if j < 5 {
			fmt.Println("============test wx re-pay==============")
			if js, err := orderwxpay.WxMoblieRemoteCall(string(bytes)); err != nil {

				t.Error("request err:", err.Error())
			} else {
				if err := json.Unmarshal([]byte(js), &rjs); err != nil {
					t.Error("Unmarshal err")
				}
				ono = rjs.Data.Ono
			}
		}

		fmt.Println("ono:", ono)
		var str string
		var err error

		if j == 0 {
			if str, err = ConfirmOrderPay(ono, way[j], ""); err != nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
		} else if 1 == j {
			//ALIPAY err
			DEM.Evb.SetErrs(0)
			DEM.Evb.AddQErr3(".*select total_price from ods_order where ono.*")
			if str, err = ConfirmOrderPay(ono, way[j], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			DEM.Evb.ClsQErr()

			DEM.Evb.SetErrs(0)
			DEM.Evb.AddQErr3(".*select name from ods_record where ono.*")
			if str, err = ConfirmOrderPay(ono, way[j], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			DEM.Evb.ClsQErr()

			if str, err = ConfirmOrderPay(ono, way[j], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err != nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			if str, err = ConfirmOrderPay("", way[j], ""); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			if str, err = ConfirmOrderPay("9999999", way[j], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			if str, err = ConfirmOrderPay(gonoPaid, way[j], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
		} else if 2 == j {

			//Check is paid
			DEM.Evb.SetErrs(0)
			DEM.Evb.AddQErr3(".*select status from.*")
			if str, err = ConfirmOrderPay(ono, way[3], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			DEM.Evb.ClsQErr()

			DEM.Evb.SetErrs(0)
			DEM.Evb.AddQErr3(".*select count.*")
			if str, err = ConfirmOrderPay(ono, way[3], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			DEM.Evb.ClsQErr()

			DEM.Evb.SetErrs(0)
			DEM.Evb.AddQErr3(".*select expand from ods_order where ono.*")
			if str, err = ConfirmOrderPay(ono, way[3], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			DEM.Evb.ClsQErr()

			DEM.Evb.SetErrs(0)
			DEM.Evb.AddQErr3(".*update ods_order set expand.*")
			if str, err = ConfirmOrderPay(ono, way[3], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			DEM.Evb.ClsQErr()

			if str, err = ConfirmOrderPay(ono, way[3], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err != nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			if str, err = ConfirmOrderPay(ono, "heheda", "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
		}
		if 3 == j {
			//WX err
			DEM.Evb.SetErrs(0)
			DEM.Evb.AddQErr3(".*select total_price from ods_order where ono.*")
			if str, err = ConfirmOrderPay(ono, way[2], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			DEM.Evb.ClsQErr()

			DEM.Evb.SetErrs(0)
			DEM.Evb.AddQErr3(".*select name from ods_record where ono.*")
			if str, err = ConfirmOrderPay(ono, way[2], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err == nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
			DEM.Evb.ClsQErr()

			if str, err = ConfirmOrderPay(ono, way[2], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err != nil {
				fmt.Println("corfirm pay err: ", err)
				t.Error("corfirm pay err: ", err.Error())
			}
		}

		fmt.Println(str)

	}

	//deleteDataOfPerf()
}

type Msg struct {
	Appid     string `json:"appid"`
	Noncestr  string `json:"noncestr"`
	Package   string `json:package`
	Partnerid string `json:partnerid`
	Prepayid  string `json:prepayid`
	Sign      string `json:sign`
	Timestamp string `json:timestamp`
	Ono       string `json:ono`
}

func TestAlipayWebReturn(t *testing.T) {
	//var v url.Values = make(map[string][]string)
	ts := httptest.NewServer(http.HandlerFunc(AlipayWebReturn))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+`?buyer_email=327468120%40qq.com&buyer_id=2088502994384781&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%252Fvwbh3InUFZ5MZmpINkGNvruvoiQtqmvUFHqbJsaJgn3gM5vzcZMCtfR8&notify_time=2016-01-26+19%3A24%3A54&notify_type=trade_status_sync&out_trade_no=2016012619235921896&payment_type=1&seller_email=itdayang%40gmail.com&seller_id=2088501949844011&subject=%E6%B5%8B%E8%AF%95123&total_fee=0.01&trade_no=2016012621001004780097729115&trade_status=TRADE_SUCCESS&sign=0d0349e4f7beba835996789cd3fded8d&sign_type=MD5`, strings.NewReader(""))
	fmt.Println("http.NewRequest:", req)
	s, err := http.DefaultClient.Do(req)
	fmt.Println("s:", s)
	if err != nil {
		t.Error("http.DefaultClient.Do:", err.Error())
	}
	got, _ := ioutil.ReadAll(s.Body)
	fmt.Printf("io got:%s\n", got)
	//在此向支付宝发校验一定出错
	if strings.Contains(string(got), "html") == true {
		t.Error("AlipayWebReturn error")
	}
	// req, _ := http.NewRequest("POST", UnifiedorderApi, strings.NewReader(reqParam))
	// req.Header.Set("Content-Type", "application/xml ")
	//resp, err := http.DefaultClient.Do(req)

	//?buyer_email=327468120%40qq.com&buyer_id=2088502994384781&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%252Fvwbh3InUFZ5MZmpINkGNvruvoiQtqmvUFHqbJsaJgn3gM5vzcZMCtfR8&notify_time=2016-01-26+19%3A24%3A54&notify_type=trade_status_sync&out_trade_no=2016012619235921896&payment_type=1&seller_email=itdayang%40gmail.com&seller_id=2088501949844011&subject=%E6%B5%8B%E8%AF%95123&total_fee=0.01&trade_no=2016012621001004780097729115&trade_status=TRADE_SUCCESS&sign=0d0349e4f7beba835996789cd3fded8d&sign_type=MD5
}

//982c22047683947542fc8cdb5153f2dm0o   noticeID
//RqPnCoPT3K9%2Fvwbh3InUFZ5KraeT0DKYTlycJYRIIIQ9046FX2Et9bohMniCY3mN%2FEug  returnID
//web mobile
func TestAlipayWebNotify(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(AlipayWebNotify))
	defer ts.Close()
	var v url.Values = make(map[string][]string)
	v.Add("buyer_email", `327468120@qq.com`)
	v.Add("buyer_id", `2088502994384781`)
	v.Add("exterface", `create_direct_pay_by_user`)
	v.Add("is_success", `T`)
	v.Add("notify_id", `RqPnCoPT3K9%2Fvwbh3InUFZ5MZmpINkGNvruvoiQtqmvUFHqbJsaJgn3gM5vzcZMCtfR8`)
	v.Add("notify_time", `2016-01-26 19:24:54`)
	v.Add("notify_type", `trade_status_sync`)
	v.Add("out_trade_no", `2016012619235921896`)
	v.Add("payment_type", `1`)
	v.Add("seller_email", `itdayang@gmail.com`)
	v.Add("seller_id", `2088501949844011`)
	v.Add("subject", `测试123`)
	v.Add("total_fee", `0.01`)
	v.Add("trade_no", `2016012621001004780097729115`)
	v.Add("trade_status", `TRADE_SUCCESS`)
	v.Add("sign", `0d0349e4f7beba835996789cd3fded8d`)
	v.Add("sign_type", `MD5`)
	req, _ := http.NewRequest("POST", ts.URL, strings.NewReader(v.Encode()))
	s, _ := http.DefaultClient.Do(req)
	got, _ := ioutil.ReadAll(s.Body)
	fmt.Printf("io got:%s\n", got)
	if string(got) != "fail" {
		t.Error("notify getResponse error")
	}

}

// _input_charset="UTF-8"
// &it_b_pay="30m"
// &notify_url="http://order.dev.jxzy.com/alipay-mobile-notify"
// &out_trade_no="2016011118105694060"
// &partner="2088501949844011"
// &payment_type="1"
// &seller_id="itdayang@gmail.com"
// &service="mobile.securitypay.pay"
// &subject="学好语文"
// &total_fee="0.01"
// &sign="XVJS8z8IiJG2ZKvfF%2BVmhcScPxtVtw1H1QBgkmzhL7TH101Fu7%2F6Iebadj2VJuIg81iCfK7C0HyiAPKG8w9xuRo7v8Ac9Nd4q79PjABxBQTWoIu7zchIgTjpUQtYf1mq3aeiMxpn06nauqUZ5OIRbOTq%2BWzeK8aDp0f6WV5k9j4%3D"
// &sign_type="RSA"

// mysign: cmWaW8BvWgdvfNKStRNdU3EJOZTaZjZILpKkc3cleTbwE/4l6Hz9OSijQeINfbRpY+Q9bJHNDcxA7o2bfc8VnFLvzZCMKjWaVReVDWDlBOW1xzIz94D/SDJM6LkSW+C3+cSqkl8r6ycTuM8FWhG+63i1NFUJIpw1aA+MO3JpYMA=
// _input_charset="UTF-8"
// &it_b_pay="30m"
// &notify_url="http://order.dev.jxzy.com/alipay-mobile-notify"
// &out_trade_no="2016012717202477158"
// &partner="2088501949844011"
// &payment_type="1"
// &seller_id="itdayang@gmail.com"
// &service="mobile.securitypay.pay"
// &subject="订单测试001"
// &total_fee="0.01"
// &sign="cmWaW8BvWgdvfNKStRNdU3EJOZTaZjZILpKkc3cleTbwE%2F4l6Hz9OSijQeINfbRpY%2BQ9bJHNDcxA7o2bfc8VnFLvzZCMKjWaVReVDWDlBOW1xzIz94D%2FSDJM6LkSW%2BC3%2BcSqkl8r6ycTuM8FWhG%2B63i1NFUJIpw1aA%2BMO3JpYMA%3D"
// &sign_type="RSA"

func TestAlipayMobileNotify(t *testing.T) {
	InitAlipayConfig()
	tsm := httptest.NewServer(http.HandlerFunc(AlipayMobileNotify))
	defer tsm.Close()
	var v url.Values = make(map[string][]string)
	v.Add("_input_charset", `UTF-8`)
	v.Add("it_b_pay", `30m`)
	v.Add("notify_url", `http://order.dev.jxzy.com/alipay-mobile-notify`)
	v.Add("out_trade_no", `2016012717202477158`)
	v.Add("partner", `2088501949844011`)
	v.Add("payment_type", `1`)
	v.Add("seller_id", `itdayang@gmail.com`)
	v.Add("service", `mobile.securitypay.pay`)
	v.Add("subject", `订单测试001`)
	v.Add("total_fee", `0.01`)
	// sgn, _ := url.QueryUnescape(`XVJS8z8IiJG2ZKvfF%2BVmhcScPxtVtw1H1QBgkmzhL7TH101Fu7%2F6Iebadj2VJuIg81iCfK7C0HyiAPKG8w9xuRo7v8Ac9Nd4q79PjABxBQTWoIu7zchIgTjpUQtYf1mq3aeiMxpn06nauqUZ5OIRbOTq%2BWzeK8aDp0f6WV5k9j4%3D`)
	// fmt.Printf("==sgn:%s\n", sgn)
	v.Add("sign", `cmWaW8BvWgdvfNKStRNdU3EJOZTaZjZILpKkc3cleTbwE/4l6Hz9OSijQeINfbRpY+Q9bJHNDcxA7o2bfc8VnFLvzZCMKjWaVReVDWDlBOW1xzIz94D/SDJM6LkSW+C3+cSqkl8r6ycTuM8FWhG+63i1NFUJIpw1aA+MO3JpYMA=`)
	v.Add("sign_type", `RSA`)

	reqm, _ := http.NewRequest("POST", tsm.URL, strings.NewReader(v.Encode()))
	s, _ := http.DefaultClient.Do(reqm)
	got, _ := ioutil.ReadAll(s.Body)
	fmt.Printf("io got:%s\n", got)
	if string(got) != "fail" {
		t.Error("notify getResponse error")
	}

}

func TestGetRsaSign(t *testing.T) {
	InitAlipayConfig()
	var v url.Values = make(map[string][]string)
	convey.Convey("TestGetRsaSign tradeNo nil ", t, func() {
		hs, g_rs := testconf.HsBuilder("GET", "http://test.com", v, 100, "")
		GetRsaSign(hs)
		var parse = make(map[string]interface{})
		json.Unmarshal(g_rs.Bytes(), &parse)
		convey.So(parse["code"], convey.ShouldEqual, 1)
	})
	convey.Convey("TestGetRsaSign totalFee err ", t, func() {
		v.Add("tradeNo", "123")
		v.Add("totalFee", "a123")
		hs, g_rs := testconf.HsBuilder("GET", "http://test.com", v, 100, "")
		GetRsaSign(hs)
		var parse = make(map[string]interface{})
		json.Unmarshal(g_rs.Bytes(), &parse)
		convey.So(parse["code"], convey.ShouldEqual, 1)
	})
	convey.Convey("TestGetRsaSign totalFee err ", t, func() {
		v.Set("totalFee", "-0.01")
		hs, g_rs := testconf.HsBuilder("GET", "http://test.com", v, 100, "")
		GetRsaSign(hs)
		var parse = make(map[string]interface{})
		json.Unmarshal(g_rs.Bytes(), &parse)
		convey.So(parse["code"], convey.ShouldEqual, 1)
	})
	convey.Convey("TestGetRsaSign totalFee N ", t, func() {
		v.Set("totalFee", "0.01")
		hs, g_rs := testconf.HsBuilder("GET", "http://test.com", v, 100, "")
		GetRsaSign(hs)
		var parse = make(map[string]interface{})
		json.Unmarshal(g_rs.Bytes(), &parse)
		convey.So(parse["code"], convey.ShouldEqual, 0)
	})

}

func TestTestAlipay(t *testing.T) {
	convey.Convey("TestTestAlipay  N ", t, func() {
		var v url.Values = make(map[string][]string)
		hs, g_rs := testconf.HsBuilder("GET", "http://test.com", v, 0, "")
		TestAlipay(hs)
		var parse = make(map[string]interface{})
		json.Unmarshal(g_rs.Bytes(), &parse)
		//fmt.Printf("parse:%s\n", string(g_rs.Bytes()))
		convey.So(string(g_rs.Bytes()), convey.ShouldNotBeBlank)
	})
}

func TestMobilePayTest(t *testing.T) {
	InitAlipayConfig()
	convey.Convey("test MobilePayTest", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(MobilePayTest))
		defer ts.Close()
		req, _ := http.NewRequest("GET", ts.URL, strings.NewReader(""))
		s, _ := http.DefaultClient.Do(req)
		got, _ := ioutil.ReadAll(s.Body)
		convey.So(string(got), convey.ShouldNotBeBlank)
	})
}

func CleanData(ono string) {
	db := common.DbConn()
	_sql := `delete from ods_order where ono=?`
	_, err := db.Exec(_sql, ono)
	if err != nil {
		fmt.Println("delete from ods_order err,ono=", ono)
		return
	}
	_sql = `delete from ods_order_item where ono=?`
	_, err = db.Exec(_sql, ono)
	if err != nil {
		fmt.Println("delete from ods_order_item err,ono=", ono)
		return
	}
	_sql = `delete from ods_record where ono=?`
	_, err = db.Exec(_sql, ono)
	if err != nil {
		fmt.Println("delete from ods_record err,ono=", ono)
		return
	}

}

func deleteDataOfPerf() {
	db := common.DbConn()
	_sql := `delete from ods_order where ono in (select ono from ods_order_item where p_from='TEST')`
	_, err := db.Exec(_sql)
	if err != nil {
		fmt.Println("delete from ods_order err in deleteData")
		return
	}
	_sql = `delete from ods_record where ono in (select ono from ods_order_item where p_from='TEST')`
	_, err = db.Exec(_sql)
	if err != nil {
		fmt.Println("delete from ods_record err in deleteData")
		return
	}
	_sql = `delete from ods_order_item where p_from='TEST'`
	_, err = db.Exec(_sql)
	if err != nil {
		fmt.Println("delete from ods_order_item err in deleteData")
		return
	}
}

func testAlipay(i int64, db *sql.DB, t *testing.T) {
	m := AlipayRemoteReqStruct{
		Ono:        "",
		Buyer:      267250,
		Seller:     438982,
		Subject:    "天地一号",
		TotalFee:   0.01,
		Body:       "迟到扣200",
		Type:       "N",
		Status:     "NOT_PAY",
		Return_url: "http://rcp.dev.jxzy.com/courseDetail.html?id=40040",
		Expand:     "id=40040&token=4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35",
	}
	orderi := orderModel.Item{
		Ono:      "",
		Oid:      int64(100000) + i,
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

	if _, err := AlipayRemoteRequest(string(bytes)); err != nil {
		t.Error(err.Error())
		return
	}
	var bl bool
	var err error
	var ono string
	_sql := `select ono from ods_order_item where oid=?`
	err = db.QueryRow(_sql, int64(100000)+i).Scan(&ono)
	//err = db.QueryRow(_sql, int64(10000)+i).Scan(&ono)

	if err != nil {
		t.Error(err.Error())
		return
	}
	if bl, err = orderModel.AlipayPaySuccess("天地一号", "N", 0.01, "ALIPAY", "267250", "USER", ono, "PAID"); bl != true {
		orderModel.IfAlipaySuccessFail("天地一号", "N", 0.01, "ALIPAY", "267250", "USER", ono, "N")
	}
	if err != nil {
		t.Error(err.Error())
		return
	}

}

func TestPerformance(t *testing.T) {
	// if err := InitDB(); err != nil {
	// 	fmt.Println("InitDB err:", err)
	// 	t.Error(err.Error())
	// }
	InitAlipayConfig()
	runtime.GOMAXPROCS(runtime.NumCPU())
	tlog := "pay_test.log"
	os.Remove(tlog)
	tc := 1
	db := common.DbConn()
	used, _ := DoPerf(tc, tlog, func(i int64) {
		testAlipay(i, db, t)
	})
	per := int64(0)
	if 0 != tc {
		per = used / int64(tc)
	}

	fmt.Printf(`
--------------------------------------------
--------------------------------------------
-->used:%vms,count:%d,per:%vms
--------------------------------------------
--------------------------------------------
		`, used, tc, per)
	if per > 3000 {
		//deleteDataOfPerf()
		t.Error("time out")
		return
	}
	//deleteDataOfPerf()
	deleteData()
}
