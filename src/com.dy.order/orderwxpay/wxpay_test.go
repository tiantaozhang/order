package orderwxpay

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	//"net/http"
	//"net/http/httptest"
	"com.dy.order/common"
	"com.dy.order/conf"
	//"database/sql"
	"encoding/json"
	//"github.com/Centny/gwf/tutil"
	"com.dy.order/orderModel"
	//"com.dy.order/testconf"
	"bytes"
	"github.com/Centny/gwf/routing/httptest"
	//"net/url"
	"database/sql"
	"encoding/xml"
	"github.com/Centny/DEM"
	"github.com/Centny/gwf/tutil"
	uap_cf "org.cny.uap/conf"
	"org.cny.uap/uap"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

var gono string = orderModel.NewOrderNo()
var gonoPaid string = orderModel.NewOrderNo()
var gerrno string = orderModel.NewOrderNo()
var transaction_id string = "1004430666201602173341556530"
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
	sql_ = `INSERT INTO ods_order ( ono, buyer, seller, total_price, type, status, time, return_url, expand, wno, add1, add2)VALUES ( ?, 0, 0,0.01, 'N', 'NOT_PAY', ?, 'http://rcp.dev.jxzy.com/questionPoolDetailNew.html?id=40021&eid=54891', 'id=40021&token=b4561e0e8c3e185e9ef858cc54cad5f1-9933ec14-6590-4548-9613-0cda727bfbe4', NULL, NULL, NULL)`
	_, err = t_conn.Exec(sql_, gerrno, time.Now())
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
	sql_ = `insert into ods_record(name,type,money,uid,pay_type,target_id,ono,status,add1,add2) values('寻龙诀','INCOME',0.01,267250,'ALIPAY',267250,?,'NOT_PAY',NULL,NULL)`
	_, err = t_conn.Exec(sql_, gerrno)
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

func TestWXPay(t *testing.T) {
	if err := InitDB(); err != nil {
		fmt.Println("InitDB err:", err)
		t.Error(err.Error())
	}
	InitWxConfig()
	var fee float64 = 0.01
	//onoArray := []string{}
	for j := 0; j < 4; j++ {
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
		var rjs struct {
			Code string `json:"code"`
			Data Msg    `json:"data"`
		}
		if j == 0 {
			m.Seller = 0
			bytes, _ := json.Marshal(m)
			if _, err := WxMoblieRemoteCall(string(bytes)); err == nil {

				t.Error("request err:", err.Error())
			}
			if _, err := WxNativeRemoteCall(string(bytes)); err == nil {

				t.Error("request err:", err.Error())
			}
			//marshal err

			if _, err := WxMoblieRemoteCall(string("")); err == nil {

				t.Error("request err:", err.Error())
			}
			if _, err := WxNativeRemoteCall(string("")); err == nil {

				t.Error("request err:", err.Error())
			}
			//fee
			m.Seller = 438982
			m.TotalFee = -1
			bytes, _ = json.Marshal(m)
			if _, err := WxMoblieRemoteCall(string(bytes)); err == nil {

				t.Error("request err:", err.Error())
			}
			if _, err := WxNativeRemoteCall(string(bytes)); err == nil {

				t.Error("request err:", err.Error())
			}

		} else if j == 1 {
			m.Integral = 100000
			bytes, _ := json.Marshal(m)
			if _, err := WxMoblieRemoteCall(string(bytes)); err == nil {

				t.Error("request err:", err.Error())
			}
			if _, err := WxNativeRemoteCall(string(bytes)); err == nil {

				t.Error("request err:", err.Error())
			}
		} else if j == 2 {
			m.Ono = ""
			bytes, _ := json.Marshal(m)
			if _, err := WxMoblieRemoteCall(string(bytes)); err != nil {
				t.Error("request err:%s", err.Error())
			}
			if _, err := WxNativeRemoteCall(string(bytes)); err != nil {

				t.Error("request err:", err.Error())
			}
		} else {
			bytes, _ := json.Marshal(m)
			if js, err := WxMoblieRemoteCall(string(bytes)); err != nil {

				t.Error("request err:", err.Error())
			} else {
				if err := json.Unmarshal([]byte(js), &rjs); err != nil {
					t.Error("Unmarshal err")
				}
			}
			if _, err := WxNativeRemoteCall(string(bytes)); err != nil {

				t.Error("request err:", err.Error())
			}
			fmt.Println("rjs:", rjs)
		}

		// var str string
		// var err error

		// if j == 0 {
		// 	if str, err = ConfirmOrderPay(ono, way[j], ""); err != nil {
		// 		fmt.Println("corfirm pay err: ", err)
		// 		t.Error("corfirm pay err: ", err.Error())
		// 	}
		// } else if 1 == j {
		// 	if str, err = ConfirmOrderPay(ono, way[j], "4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"); err != nil {
		// 		fmt.Println("corfirm pay err: ", err)
		// 		t.Error("corfirm pay err: ", err.Error())
		// 	}
		// }

		// fmt.Println(str)

	}

	//deleteDataOfPerf()

}

func TestWxPayMoblieNotify(t *testing.T) {
	// if err := InitDB(); err != nil {
	// 	fmt.Println("InitDB err:", err)
	// 	t.Error(err.Error())
	// }
	// InitWxConfig()

	//var v url.Values = make(map[string][]string)
	convey.Convey("TestWxpay success ", t, func() {

		tx := `<xml><appid><![CDATA[wx270243445e3233ee]]></appid>
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
		ts := httptest.NewMuxServer()
		ts.Mux.HFunc(".*", WxPayMoblieNotify)
		ts.Mux.ShowLog = true
		omg, err := ts.PostN("", "application/xml", bytes.NewBuffer([]byte(tx)))
		if err != nil {
			fmt.Printf("%s/n", err.Error())
			//panic(err)
		}
		fmt.Println("%v", omg)
		var xmlO struct {
			Suc string `xml:"return_code"`
			ReM string `xml:"return_msg"`
		}
		if err := xml.Unmarshal([]byte(omg), &xmlO); err != nil {
			t.Error("%s", err.Error())
		}
		convey.So(xmlO.Suc, convey.ShouldEqual, "SUCCESS")
		//	convey.So(omg["code"], convey.ShouldEqual, "0")
	})

	convey.Convey("TestWxpay failed ", t, func() {

		tx := `<xml><appid><![CDATA[wx270243445e3233ee]]></appid>
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
	<total_fee>1.11</total_fee>
	<trade_type><![CDATA[APP]]></trade_type>
	<transaction_id><![CDATA[1004430666201602173341556530]]></transaction_id>
	</xml>`
		ts := httptest.NewMuxServer()
		ts.Mux.HFunc(".*", WxPayMoblieNotify)
		ts.Mux.ShowLog = true
		omg, err := ts.PostN("", "application/xml", bytes.NewBuffer([]byte(tx)))
		if err != nil {
			fmt.Printf("%s/n", err.Error())
			//panic(err)
		}
		fmt.Println("%v", omg)
		var xmlO struct {
			Suc string `xml:"return_code"`
			ReM string `xml:"return_msg"`
		}
		if err := xml.Unmarshal([]byte(omg), &xmlO); err != nil {
			t.Error("%s", err.Error())
		}
		convey.So(xmlO.Suc, convey.ShouldEqual, "FAIL")
		//	convey.So(omg["code"], convey.ShouldEqual, "0")
	})
	convey.Convey("TestNativeWxpay success ", t, func() {

		tx := `<xml><appid><![CDATA[wx270243445e3233ee]]></appid>
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
		ts := httptest.NewMuxServer()
		ts.Mux.HFunc(".*", WxPayWebNotify)
		ts.Mux.ShowLog = true
		omg, err := ts.PostN("", "application/xml", bytes.NewBuffer([]byte(tx)))
		if err != nil {
			fmt.Printf("%s/n", err.Error())
			//panic(err)
		}
		fmt.Println("%v", omg)
		var xmlO struct {
			Suc string `xml:"return_code"`
			ReM string `xml:"return_msg"`
		}
		if err := xml.Unmarshal([]byte(omg), &xmlO); err != nil {
			t.Error("%s", err.Error())
		}
		convey.So(xmlO.Suc, convey.ShouldEqual, "SUCCESS")
		//	convey.So(omg["code"], convey.ShouldEqual, "0")
	})
	convey.Convey("TestWxpay failed ", t, func() {

		tx := `<xml><appid><![CDATA[wx270243445e3233ee]]></appid>
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
	<total_fee>1.11</total_fee>
	<trade_type><![CDATA[APP]]></trade_type>
	<transaction_id><![CDATA[1004430666201602173341556530]]></transaction_id>
	</xml>`
		ts := httptest.NewMuxServer()
		ts.Mux.HFunc(".*", WxPayWebNotify)
		ts.Mux.ShowLog = true
		omg, err := ts.PostN("", "application/xml", bytes.NewBuffer([]byte(tx)))
		if err != nil {
			fmt.Printf("%s/n", err.Error())
			//panic(err)
		}
		fmt.Println("%v", omg)
		var xmlO struct {
			Suc string `xml:"return_code"`
			ReM string `xml:"return_msg"`
		}
		if err := xml.Unmarshal([]byte(omg), &xmlO); err != nil {
			t.Error("%s", err.Error())
		}
		convey.So(xmlO.Suc, convey.ShouldEqual, "FAIL")
		//	convey.So(omg["code"], convey.ShouldEqual, "0")
	})
}

func TestAfterSuccess(t *testing.T) {
	// if err := InitDB(); err != nil {
	// 	fmt.Println("InitDB err:", err)
	// 	t.Error(err.Error())
	// }
	// InitWxConfig()

	convey.Convey("...", t, func() {
		var bl bool
		var err error
		IfWXPaySuccessFail(gono, "PAID", transaction_id)
		if bl, err = WXpayPaySuccess(gono, "PAID", transaction_id); bl != true {
			fmt.Println("err:", err)
			IfWXPaySuccessFail(gono, "PAID", transaction_id)
		}
		//test paid fail
		//	CleanData(strono)
		convey.So(err, convey.ShouldBeNil)
		convey.So(bl, convey.ShouldBeTrue)
	})
	convey.Convey("...", t, func() {
		var bl bool
		var err error
		IfWXPaySuccessFail(gerrno, "PAID", transaction_id)
		if bl, err = WXpayPaySuccess(gerrno, "PAID", transaction_id); bl != true {
			fmt.Println("err:", err)
		}
		//test paid fail
		//	CleanData(strono)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("update err", t, func() {

		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*update ods_record set.*")
		_, err := WXpayPaySuccess(gono, "PAID", transaction_id)
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("update err", t, func() {

		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*update ods_order set status.*")
		_, err := WXpayPaySuccess(gono, "PAID", transaction_id)
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("update err", t, func() {

		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*update ods_order set wno.*")
		bl, err := WXpayPaySuccess(gono, "PAID", transaction_id)
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(bl, convey.ShouldBeFalse)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("commit err", t, func() {

		DEM.Evb.SetErrs(DEM.TX_COMMIT_ERR)
		//	DEM.Evb.AddQErr3(".* .*")
		bl, err := WXpayPaySuccess(gono, "PAID", transaction_id)
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(bl, convey.ShouldBeFalse)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("...", t, func() {

		var bl bool
		var err error
		if bl, err = WXpayPaySuccess(gono, "PAID", transaction_id); bl != true {
			fmt.Println("err:", err)
			IfWXPaySuccessFail(gono, "PAID", transaction_id)
		}
		IfWXPaySuccessFail(gono, "PAID", transaction_id)
		//test paid fail
		convey.So(err, convey.ShouldNotBeNil)
	})

	timer := time.NewTicker(5 * time.Second)
	select {
	case <-timer.C:
		timer.Stop()
	}

}

func TestIfWXPaySuccessFail(t *testing.T) {
	IfWXPaySuccessFail(gono, "PAID", transaction_id)
	timer := time.NewTicker(5 * time.Second)
	select {
	case <-timer.C:
		timer.Stop()
	}
}

func AssemblePara() WXRemoteReqStruct {
	m := WXRemoteReqStruct{
		Ono:        "",
		Buyer:      267250,
		Seller:     438982,
		Subject:    "testWXPay",
		TotalFee:   float64(0.01),
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
	return m

}

func TestDealWXOrder(t *testing.T) {
	// if err := InitDB(); err != nil {
	// 	fmt.Println("InitDB err:", err)
	// 	t.Error(err.Error())
	// }
	var str string
	var err error
	convey.Convey("DealWXOrder  N", t, func() {
		DEM.Evb.SetErrs(0)
		m := AssemblePara()
		str, err = DealWXOrder(m, gono, "")
		str, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("DealWXOrder  insert into ods_order err", t, func() {
		m := AssemblePara()

		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert into ods_order\\(ono,buyer,seller,total_price,type,status,return_url,expand,wno\\).*")
		_, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("DealWXOrder begin err", t, func() {
		m := AssemblePara()

		DEM.Evb.SetErrs(DEM.OPEN_ERR | DEM.BEGIN_ERR)
		// DEM.Evb.AddQErr3(".*insert into ods_order_item.*")
		_, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("DealWXOrder TotalFee err", t, func() {
		DEM.Evb.SetErrs(0)
		m := AssemblePara()
		m.TotalFee = 0
		m.Integral = 0
		_, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("DealWXOrder orderItem err", t, func() {
		m := WXRemoteReqStruct{
			Ono:        "",
			Buyer:      267250,
			Seller:     438982,
			Subject:    "testWXPay",
			TotalFee:   float64(0.01),
			Body:       "test",
			Type:       "N",
			Status:     "NOT_PAY",
			Return_url: "http://rcp.dev.jxzy.com/courseDetail.html?id=40040",
			Expand:     "id=40040&token=4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35",
		}
		_, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("DealWXOrder insert ods_order_item err", t, func() {
		m := AssemblePara()

		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert into ods_order_item.*")
		_, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("DealWXOrder  InsertWithIntegral err", t, func() {
		m := AssemblePara()

		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert into ods_record.*")
		_, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("DealWXOrder payway=1 InsertWithIntegral err", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert into ods_record.*")
		m := AssemblePara()
		m.TotalFee = 0.02
		m.Integral = 1
		_, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
	})

	convey.Convey("DealWXOrder  commit err", t, func() {
		m := AssemblePara()

		DEM.Evb.SetErrs(DEM.TX_COMMIT_ERR)
		_, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("DealWXOrder  sync err", t, func() {
		m := AssemblePara()
		m.Buyer = 12332112321
		str, err = DealWXOrder(m, gono, "123")
		convey.So(err, convey.ShouldNotBeNil)
	})
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

func InitDB() error {
	// err := common.Init("mysql", "cny:123@tcp(192.168.2.57:3306)/ucs?charset=utf8&loc=Local")
	// if err != nil {
	// 	fmt.Println("init db err:", err)
	// 	return err
	// }
	// err = common.Init("mysql", "cny:123@tcp(192.168.2.57:3306)/orderv2?charset=utf8&loc=Local")
	// //var err error
	// // db, err = sql.Open("mysql", mydbname)("root:123456@tcp(127.0.0.1:3306)/mydb?charset=utf8")

	// if err != nil {
	// 	fmt.Println("init db err:", err)
	// }
	// cfg := "../../../conf/order.properties"
	// //cfg := "/Users/xxx/code/go/src/order/conf/order.properties"
	// err = conf.Cfg.InitWithFilePath(cfg)
	// uap_cf.Cfg = conf.Cfg
	// //conf.Cfg.Print()
	// if err != nil {
	// 	fmt.Println(cfg)
	// 	//panic(err)
	// 	return err
	// }
	// uap.InitDb(common.DbConn)
	// return nil

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

func testN(str string) {
	tx := `<xml><appid><![CDATA[wx270243445e3233ee]]></appid>
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
	ts := httptest.NewMuxServer()
	ts.Mux.HFunc(".*", WxPayMoblieNotify)
	ts.Mux.ShowLog = true
	omg, err := ts.PostN("", "application/xml", bytes.NewBuffer([]byte(tx)))
	if err != nil {
		fmt.Printf("%s/n", err.Error())
		//panic(err)
	}
	fmt.Println("%v", omg)
	var xmlO struct {
		Suc string `xml:"return_code"`
		ReM string `xml:"return_msg"`
	}
	if err := xml.Unmarshal([]byte(omg), &xmlO); err != nil {
		fmt.Printf("%s", err.Error())
	}

}

func testWXpay(i int, db *sql.DB, t *testing.T) {
	m := AssemblePara()
	jm, _ := json.Marshal(m)
	strW, err := WxMoblieRemoteCall(string(jm))
	if err != nil {
		fmt.Println("WxMoblieRemoteCall err", err)
		return
	}
	//
	testN(strW)
}

func TestPerformance(t *testing.T) {
	// if err := InitDB(); err != nil {
	// 	fmt.Println("InitDB err:", err)
	// 	t.Error(err.Error())
	// }
	// InitWxConfig()
	runtime.GOMAXPROCS(runtime.NumCPU())
	tlog := "pay_test.log"
	os.Remove(tlog)
	tc := 10
	db := common.DbConn()
	used, _ := tutil.DoPerf(tc, tlog, func(i int) {
		testWXpay(i, db, t)
	})
	per := used / int64(tc)

	fmt.Printf(`
--------------------------------------------
--------------------------------------------
-->used:%vms,count:%d,per:%vms
--------------------------------------------
--------------------------------------------
		`, used, tc, per)
	if per > 30000 {
		//deleteDataOfPerf()
		deleteData()
		t.Error("time out")
		return
	}
	//deleteDataOfPerf()
	deleteData()
}
