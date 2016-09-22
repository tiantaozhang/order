package orderModel

import (
	"com.dy.order/common"
	"com.dy.order/conf"
	"database/sql"
	"fmt"
	"github.com/Centny/DEM"
	//"github.com/Centny/gwf/tutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartystreets/goconvey/convey"
	uap_cf "org.cny.uap/conf"
	"org.cny.uap/uap"
	//"runtime"
	"strings"
	"testing"
	"time"
)

var t_conn *sql.DB
var t_db *sql.DB

var ono string = NewOrderNo()

func initDB() {

	//cfg := "/Users/xxx/code/go/src/order/conf/order.properties"
	cfg := "../../../conf/order.properties"
	err := conf.Cfg.InitWithFilePath(cfg)
	uap_cf.Cfg = conf.Cfg
	//conf.Cfg.Print()
	if err != nil {
		fmt.Println(cfg)
		panic(err)
		return
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

	// t_open := DEM.OpenDem()
	// fmt.Println("t_open:", t_open)
	// cny:123@tcp(192.168.2.57:3306)/orderv2_test?charset=utf8&loc=Local
	uap.InitDb(common.DbConn)
	return
}

func insertData() {
	time.Now()
	sql_ := `INSERT INTO ods_order ( ono, buyer, seller, total_price, type, status, time, return_url, expand, wno, add1, add2)VALUES ( ?, 452638, 267250, 0.01, 'N', 'NOT_PAY', ?, 'http://rcp.dev.jxzy.com/questionPoolDetailNew.html?id=40021&eid=54891', 'id=40021&token=b4561e0e8c3e185e9ef858cc54cad5f1-9933ec14-6590-4548-9613-0cda727bfbe4', NULL, NULL, NULL)`
	_, err := t_conn.Exec(sql_, ono, time.Now())
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	//fmt.Println("sql_ ods_order_item")
	sql_ = `INSERT INTO ods_order_item ( ono, p_name, p_id, p_type, p_img, p_from,status, time, add1, add2)VALUES ( ?, 'testSyn', 41890, 10, NULL, 'TEST', 'N',?,NULL, NULL)`
	_, err = t_conn.Exec(sql_, ono, time.Now())
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	sql_ = `INSERT INTO uap_attr( a_key,oid,owner,ns,integral,add1, add2)VALUES ('USER_INTEGRAL',41890,'USER',DEFAULT,1000,NULL, NULL)`
	_, err = t_conn.Exec(sql_)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	sql_ = `insert into ods_record(name,type,money,uid,pay_type,target_id,ono,status,add1,add2) values('寻龙诀','INCOME',0.01,41890,'ALIPAY',267250,?,'NOT_PAY',NULL,NULL)`
	_, err = t_conn.Exec(sql_, ono)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	_, err = t_conn.Exec(`insert into ods_order_env(akey,aval,type,status) values('TEST','127.0.0.1','PAID_CB','N')`)
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
	t_conn.Exec("delete from uap_attr where oid=41890")

}

func TestCheckIsExistAndInsertDB(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	initDB()
	fmt.Println("ono:", ono)
	insertData()
	convey.Convey("CheckIsExist", t, func() {
		isOnoExist, _ := CheckIsExist(t_db, ono, 0, "ods_order")
		convey.So(isOnoExist, convey.ShouldBeTrue)
	})
	convey.Convey("checkExist err", t, func() {
		isOnoExist, _ := CheckIsExist(t_db, ono, 41890, "ods_order_item")
		convey.So(isOnoExist, convey.ShouldBeTrue)
	})
	convey.Convey("checkExist err", t, func() {
		isOnoExist, _ := CheckIsExist(t_db, "", 0, "ods_order")
		convey.So(isOnoExist, convey.ShouldBeFalse)
	})
	convey.Convey("checkExist err1", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*select count.*")
		isOnoExist, _ := CheckIsExist(DEM.OpenDem(), ono, 0, "ods_order")
		convey.So(isOnoExist, convey.ShouldBeFalse)
		DEM.Evb.ClsQErr()
	})

}

func TestIntegral(t *testing.T) {
	convey.Convey("GetIntegral N", t, func() {
		igral, err := GetIntegral(t_conn, 41890)
		convey.So(err, convey.ShouldBeNil)
		convey.So(igral, convey.ShouldEqual, 1000)
	})
	convey.Convey("DetectIntegral N", t, func() {
		err := DetectIntegral(t_conn, 90, 41890, 90)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("Integral not enough", t, func() {
		err := DetectIntegral(t_conn, 10000, 41890, 90)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("DetectIntegral buyer err", t, func() {
		err := DetectIntegral(t_conn, 10000, 0, 90)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("DetectIntegral totalFee err", t, func() {
		err := DetectIntegral(t_conn, 100, 41890, 0.9)
		convey.So(err, convey.ShouldNotBeNil)
	})
}

func TestPay(t *testing.T) {
	convey.Convey("PayTypeAndNeedPay total igral", t, func() {
		_, _needpay, _payway := PayTypeAndNeedPay(100, 1.0)
		convey.So(_needpay, convey.ShouldEqual, 0)
		convey.So(_payway, convey.ShouldEqual, 2)
	})
	convey.Convey("PayTypeAndNeedPay total igral", t, func() {
		_, _, _payway := PayTypeAndNeedPay(100, 1.1)
		convey.So(_payway, convey.ShouldEqual, 1)
	})
	convey.Convey("PayTypeAndNeedPay total igral", t, func() {
		_, _needpay, _payway := PayTypeAndNeedPay(0, 1.1)
		convey.So(_needpay, convey.ShouldEqual, float64(1.1))
		convey.So(_payway, convey.ShouldEqual, 0)
	})
}

func AssemblePara() CommonRemoteReqStruct {
	m := CommonRemoteReqStruct{
		Ono:        "",
		Buyer:      41890,
		Seller:     438982,
		Subject:    "寻龙诀",
		TotalFee:   0.11,
		Body:       "迟到扣200",
		Type:       "N",
		Status:     "NOT_PAY",
		Return_url: "http://rcp.dev.jxzy.com/courseDetail.html?id=40040",
		Expand:     "id=40040&token=4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35",
		Integral:   10,
	}
	//orderI := []map[string]interface{}{}
	for i := 1; i < 2; i++ {
		stri := fmt.Sprintf("%d", i)
		str := "物品" + stri
		orderi := Item{
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

func TestTx(t *testing.T) {
	// initDB()
	// fmt.Println("ono:", ono)
	// insertData()

	m := AssemblePara()
	tx, err := t_db.Begin()
	if err != nil {
		fmt.Printf("start a transaction err,%s\n", err.Error())
	}
	convey.Convey("insertOdsOrder N", t, func() {
		err := InsertOdsOrder(tx, m, ono)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("insertOdsOrder err", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert.*ods_order.*")
		err := InsertOdsOrder(tx, m, ono)
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})

	convey.Convey("InsertOrderItem N", t, func() {
		tx, err = t_db.Begin()
		if err != nil {
			fmt.Printf("start a transaction err,%s\n", err.Error())
		}
		err := InsertOrderItem(tx, m, ono)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("InsertOrderItem err", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert.*ods_order_item.*")
		err := InsertOrderItem(tx, m, ono)
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})

	convey.Convey("InsertWithIntegral N", t, func() {
		tx, err = t_db.Begin()
		if err != nil {
			fmt.Printf("start a transaction err,%s\n", err.Error())
		}
		err := InsertWithIntegral(tx, m, ono)
		convey.So(err, convey.ShouldBeNil)
		igral := 0
		if err := tx.QueryRow(`select integral from uap_attr where oid=?`, m.Buyer).Scan(&igral); err != nil {
			t.Error("query integral err:", err)
		}
		convey.So(igral, convey.ShouldEqual, 990)
	})
	convey.Convey("InsertWithIntegral updateErr", t, func() {
		m.Buyer = 102938
		err := InsertWithIntegral(tx, m, ono)
		convey.So(err, convey.ShouldNotBeNil)

	})
	convey.Convey("InsertWithIntegral err", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*update.*uap_attr.*")
		err := InsertWithIntegral(tx, m, ono)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("InsertOrderItem err", t, func() {
		tx, err = t_db.Begin()
		if err != nil {
			fmt.Printf("start a transaction err,%s\n", err.Error())
		}
		m.Buyer = 41890
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*insert.*ods_order_item.*")
		err := InsertOrderItem(tx, m, ono)
		convey.So(err, convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("InsertWithRMB N", t, func() {
		tx, err := t_db.Begin()
		if err != nil {
			fmt.Printf("start another transaction err,%s\n", err.Error())
		}
		err = InsertWithRMB(tx, m, ono, 0.01, "ALIPAY")
		convey.So(err, convey.ShouldBeNil)
		tx.Rollback()
	})
	convey.Convey("InsertWithRMB N", t, func() {
		err := InsertWithRMB(tx, m, ono, 0.01, "ALIPAY")
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("CheckParas N", t, func() {
		err := CheckParas(m)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("CheckParas err", t, func() {
		m.Return_url = ""
		err := CheckParas(m)
		convey.So(err, convey.ShouldNotBeNil)
		m.Return_url = "http://rcp.dev.jxzy.com/courseDetail.html?id=40040"
		m.Expand = ""
		err = CheckParas(m)
		convey.So(err, convey.ShouldNotBeNil)
		m.Expand = "id=40040&token=4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35"
	})
	convey.Convey("CheckParas err", t, func() {
		m.Buyer = 0
		err := CheckParas(m)
		convey.So(err, convey.ShouldNotBeNil)
		m.Buyer = 41890
		m.Seller = 0
		err = CheckParas(m)
		convey.So(err, convey.ShouldNotBeNil)
		m.Seller = 438982
	})
	convey.Convey("UpdateRecord err", t, func() {
		_, err := UpdateRecord(tx, m.Type, "PAID", m.Buyer, ono)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("UpdateRecord N", t, func() {
		tx, err := t_conn.Begin()
		Type := "PAID"
		sts := "PAID"
		if err != nil {
			fmt.Printf("start another transaction err,%s\n", err.Error())
		}
		_, err = UpdateRecord(tx, Type, sts, m.Buyer, ono)
		convey.So(err, convey.ShouldBeNil)
		tx.Commit()
	})

}

func TestCheckPaidcb(t *testing.T) {

	DEM.Evb.SetErrs(0)
	paid_cb := CheckPaidcb("TEST")
	if paid_cb != nil {
		t.Error("paid_cb err")
	}
}

func TestSyn(t *testing.T) {
	defer deleteData()
	convey.Convey("syn", t, func() {
		err := SynUser(0, 0)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("syn", t, func() {
		SynUser(438982, 267250)
		err := SynUser(438982, 0)
		convey.So(err, convey.ShouldNotBeNil)
	})

}
