package orderRecord

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAddRecord(t *testing.T) {
	convey.Convey("add record", t, func() {
		err := AddOrderRecord("name", "INCOME", 199.90, 200353, "ALIPAY", "20004", "USER", "149590394829", "N")
		convey.So(err, convey.ShouldEqual, nil)
	})
}

func TestGetRecord(t *testing.T) {
	convey.Convey("get record", t, func() {
		_, err := GetOrderRecord(200353)
		convey.So(err, convey.ShouldEqual, nil)
	})
}

func TestCheckUsrPaidItem(t *testing.T) {
	convey.Convey("check paid", t, func() {
		_, err := CheckUsrPaidItem(200353, 23891)
		convey.So(err, convey.ShouldEqual, nil)
	})
}
