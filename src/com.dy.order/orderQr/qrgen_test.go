package orderQr

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenQr(t *testing.T) {
	convey.Convey("string is nil", t, func() {
		_, err := GenQr("", "")
		convey.So(err, convey.ShouldNotBeNil)
	})
	convey.Convey("normal", t, func() {
		_, err := GenQr("weixin://wxpay/bizpayurl?pr=qMfOOs8", "123")
		convey.So(err, convey.ShouldBeNil)
	})
}
