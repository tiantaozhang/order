package orderModel

import (
	//"com.dy.order/common"
	//"com.dy.order/conf"
	//"database/sql"
	"fmt"
	//"github.com/Centny/DEM"
	//"github.com/Centny/gwf/tutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartystreets/goconvey/convey"
	//"runtime"
	"github.com/Centny/DEM"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

var sono string = NewOrderNo()

func TestCallback(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// initDB()
	insertData()
	convey.Convey("Callback ", t, func() {
		_, err := Callback(sono)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("Callback ", t, func() {
		_, err := Callback(ono)
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("Callback ", t, func() {
		t_conn.Exec("update `ods_order_env` set `aval` = 'http://rcp.dev.jxzy.com/usr/purchase-course?' where `akey`='TEST'")
		_, err := Callback(ono)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey(" ", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*select expand from.*")
		isOnoExist, _ := Callback(ono)
		convey.So(isOnoExist, convey.ShouldBeFalse)
		DEM.Evb.ClsQErr()
	})
}

func TestAfterSuccess(t *testing.T) {
	convey.Convey("...", t, func() {
		m := AssemblePara()
		var bl bool
		var err error
		if bl, err = AlipayPaySuccess(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", ono, "PAID"); bl != true {
			fmt.Println("err:", err)
			IfAlipaySuccessFail(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", ono, "PAID")
		}
		//test paid fail
		IfAlipaySuccessFail(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", ono, "PAID")
		//	CleanData(strono)
		convey.So(err, convey.ShouldBeNil)
		convey.So(bl, convey.ShouldBeTrue)
	})
	convey.Convey("update err", t, func() {
		m := AssemblePara()
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*update ods_record set.*")
		bl, err := AlipayPaySuccess(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", ono, "PAID")
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(bl, convey.ShouldBeFalse)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("update err", t, func() {
		m := AssemblePara()
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*update ods_order set status.*")
		bl, err := AlipayPaySuccess(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", ono, "PAID")
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(bl, convey.ShouldBeFalse)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("update err", t, func() {
		m := AssemblePara()
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*update ods_order set wno.*")
		bl, err := AlipayPaySuccess(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", ono, "PAID")
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(bl, convey.ShouldBeFalse)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("commit err", t, func() {
		m := AssemblePara()
		DEM.Evb.SetErrs(DEM.TX_COMMIT_ERR)
		//	DEM.Evb.AddQErr3(".* .*")
		bl, err := AlipayPaySuccess(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", ono, "PAID")
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(bl, convey.ShouldBeFalse)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("...", t, func() {
		m := AssemblePara()
		var bl bool
		var err error
		if bl, err = AlipayPaySuccess(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", sono, "PAID"); bl != true {
			fmt.Println("err:", err)
			IfAlipaySuccessFail(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", sono, "PAID")
		}
		IfAlipaySuccessFail(m.Subject, "N", 0.01, "ALIPAY", string(m.Seller), "USER", sono, "PAID")
		//test paid fail
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(bl, convey.ShouldBeFalse)
	})

	timer := time.NewTicker(5 * time.Second)
	select {
	case <-timer.C:
		timer.Stop()
	}

}

func TestDealAliReturn(t *testing.T) {
	r, err := http.NewRequest("GET", "", strings.NewReader(""))
	if err != nil {
		t.Error("NewRequest: %s", err.Error())
	}
	r.Form = make(url.Values)
	r.Form.Set("trade_status", "TRADE_SUCCESS")
	DealAliReturn(r)
	r.Form.Set("out_trade_no", ono)
	DealAliReturn(r)
	r.Form.Set("trade_status", "TRADE_FINISHED")
	DealAliReturn(r)
}

func TestDealAliNotify(t *testing.T) {
	r, err := http.NewRequest("POST", "", strings.NewReader(""))
	if err != nil {
		t.Error("NewRequest: %s", err.Error())
	}
	r.Form = make(url.Values)
	r.Form.Set("trade_status", "TRADE_SUCCESS")
	DealAliNotify(r, `ALIPAY`, "web")
	r.Form.Set("trade_status", "TRADE_FINISHED")
	DealAliNotify(r, `ALIPAY`, "web")

	//delete data
	deleteData()
}
