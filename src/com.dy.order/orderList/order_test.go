package orderList

import (
	"fmt"
	"github.com/Centny/gwf/tutil"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"testing"
	"com.dy.order/common"
	"database/sql"
	"github.com/smartystreets/goconvey/convey"
	"github.com/Centny/DEM"
)

var t_conn *sql.DB

func insertData(){
	sql_ := `INSERT INTO ods_order (tid, ono, buyer, seller, total_price, type, status, time, return_url, expand, wno, add1, add2)VALUES (9466, '2016012618114979877', 200353, 267250, 0.01, 'N', 'NOT_PAY', '2016-01-26 18:11:50', 'http://xx.dev.jxzy.com/pay-result.html?id=40025&eid=0', 'id=40025&token=d76439a5266735ce371e844038aef759-5730e99a-1da7-4ab2-85b9-98106daeab63', NULL, NULL, NULL)`
	_,err :=t_conn.Exec(sql_)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	t_conn.Exec("INSERT INTO `ods_order_item` (`tid`, `ono`, `oid`, `p_name`, `p_id`, `p_type`, `p_img`, `p_count`, `p_from`, `notified`, `price`, `type`, `status`, `time`, `add1`, `add2`)VALUES(11415, '2016012618114979877', 0, '第二个付费课程哈哈哈', 40025, '10', 'http://u.dev.jxzy.com/fLRlKA==', 1, 'RCP', 0, 0.01, 'N', 'N', '2016-01-26 18:11:49', NULL, NULL)")
	t_conn.Exec("INSERT INTO `ods_order` (`tid`, `ono`, `buyer`, `seller`, `total_price`, `type`, `status`, `time`, `return_url`, `expand`, `wno`, `add1`, `add2`)VALUES(9516, '2016012710275443433', 453177, 267250, 0.01, 'N', 'PAID', '2016-01-27 10:43:45', 'http://rcp.dev.jxzy.com/pay-result.html?id=40025&eid=0', 'id=40025&token=8dad86e559c60901330b126861a447ae-3a2aa1fe-778d-4477-8da7-d969b73de7e7', NULL, NULL, NULL)")
}

func deleteData(){
	sql_ := `delete from ods_order`
	_,err :=t_conn.Exec(sql_)
	if err != nil {
		fmt.Println(err.Error())
	}
	t_conn.Exec("delete from ods_order_item")
}

func testOrderList(t *testing.T) {
	convey.Convey("testOrderList", t, func() {
		res, err := GetUsrOrder(200353, 0, "", 10, 1)
		if err != nil {
			t.Error(err.Error())
		}
		convey.So(len(res.List),convey.ShouldEqual,1)
	})
	convey.Convey("testOrderList", t, func() {
		res, err := GetUsrOrder(0, 267250, "", 10, 1)
		if err != nil {
			t.Error(err.Error())
		}
		convey.So(len(res.List),convey.ShouldEqual,2)
	})
	convey.Convey("testOrderList", t, func() {
		res, err := GetUsrOrder(200353, 0, "2016012618114979877", 10, 1)
		if err != nil {
			t.Error(err.Error())
		}
		convey.So(len(res.List),convey.ShouldEqual,1)
	})
	convey.Convey("testOrderList", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*select a.ono,a.buyer,a.seller,a.total_price,a.type,a.status,a.time,b.usr seller_name,c.usr buyer_name from ods_order.*")
		res, err := GetUsrOrder(200353, 0, "2016012618114979877", 10, 1)
		convey.So(len(res.List),convey.ShouldEqual,	0)
		convey.So(err,convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("testOrderList", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*select count.*")
		res, err := GetUsrOrder(200353, 0, "2016012618114979877", 10, 1)
		convey.So(len(res.List),convey.ShouldEqual,	0)
		convey.So(err,convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("testOrderList", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*select p_name,p_id,p_type,p_img,p_count,price from ods_order_item.*")
		res, err := GetUsrOrder(200353, 0, "2016012618114979877", 10, 1)
		convey.So(len(res.List),convey.ShouldEqual,	0)
		convey.So(err,convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
}

func testCancelOrder(t *testing.T){
	convey.Convey("testCancelOrder", t, func() {
		err := DbCancelUpdate("2016012618114979877")
		convey.So(err,convey.ShouldBeNil)
	})
	convey.Convey("testCancelOrder1", t, func() {
		err := DbCancelUpdate("")
		convey.So(err.Error(),convey.ShouldEqual,"订单编号为空")
	})
	convey.Convey("testCancelOrder1", t, func() {
		err := DbCancelUpdate("2016012710275443433")
		convey.So(err.Error(),convey.ShouldEqual,"订单不处于可关闭状态")
	})
	convey.Convey("testCancelOrder1", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*select status from ods_order where.*")
		err := DbCancelUpdate("2016012618114979877")
		convey.So(err,convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
	convey.Convey("testCancelOrder1", t, func() {
		DEM.Evb.SetErrs(0)
		DEM.Evb.AddQErr3(".*UPDATE ods_order,ods_order_item,ods_record SET.*")
		err := DbCancelUpdate("2016012618114979877")
		convey.So(err,convey.ShouldNotBeNil)
		DEM.Evb.ClsQErr()
	})
}

func TestPerformance(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tc := 1
	DEM.G_Dn = "mysql"
	DEM.G_Dsn = "cny:123@tcp(192.168.2.57:3306)/orderv2_test?charset=utf8"
	err := common.Init("DEM", "cny:123@tcp(192.168.2.57:3306)/orderv2_test?charset=utf8&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
//	t_conn = common.DbConn()
	t_conn=DEM.LAST
	insertData()
	defer deleteData()
	used, err := tutil.DoPerf(tc, "", func(i int) {
		testOrderList(t)
		testCancelOrder(t)
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf(`
		-----------------------------
		 Used:%vms,Count:%v,Per:%vms
		-----------------------------
		`, used, tc, used/int64(tc))
}
