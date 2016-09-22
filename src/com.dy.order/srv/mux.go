package srv

import (
	"com.dy.order/conf"
	"com.dy.order/orderalipay"
	"com.dy.order/orderwxpay"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/filter"
	"github.com/tomlovzki/go_swagger"
	"net/http"
	"org.cny.uap/sync"
	"org.cny.uap/uap"
	"org.cny.uas/sso"
)

func NewSrvMux(pre string, www string) *routing.SessionMux {
	sb := routing.NewSrvSessionBuilder("", "/", "fs", 30*60*1000, 10*1000)
	mux := routing.NewSessionMux(pre, sb)
	mux.ShowLog = false
	af := sso.NewAuthFilter(conf.SsoLoginUrl(), conf.SsoAuthUrl(), "")
	af2 := sso.NewAuthFilter(
		conf.SsoLoginUrl(),
		conf.SsoAuthUrl(),
		conf.UrlRoot())

	af2.Optioned = true
	mux.HFunc("^/ok(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		return hs.MsgRes("ok")
	})

	cors := filter.NewCORS()
	cors.AddSite("*")
	mux.HFilter("^/.*$", cors)

	mux.HFilterFunc("^/.*$", filter.NoCacheFilter)
	mux.HFilterFunc("^/logout(\\?.*)?$", uap.ClsUsrStmtFilter)
	mux.HFilterFunc("^/logout(\\?.*)?$", sso.ClsAuthFilter)
	mux.H("^/logout(\\?.*)?$", sso.NewRedirect3(conf.SsoLogoutUrl(), "http://%s"))
	mux.HFilter("^/usr.*$", af)
	mux.HFilter("^/usr.*$", sync.NewSyncUsrStmtFilter())
	//mux.HandleFunc("^/go-pay(\\?.*)?$", orderalipay.AlipayWebRequest)
	mux.HandleFunc("^/alipay-web-notify(\\?.*)?$", orderalipay.AlipayWebNotify)
	mux.HandleFunc("^/alipay-web-return(\\?.*)?$", orderalipay.AlipayWebReturn)
	//
	mux.HandleFunc("^/alipay-mobile-request(\\?.*)?$", orderalipay.MobilePayTest)
	mux.HandleFunc("^/alipay-mobile-notify(\\?.*)?$", orderalipay.AlipayMobileNotify)
	//
	//mux.HandleFunc("^/alipay-web-request(\\?.*)?$", orderalipay.AlipayWebRequest)
	//wechatpay-notify
	mux.HFunc("^/wxpay-mobile-notify(\\?.*)?$", orderwxpay.WxPayMoblieNotify)
	mux.HFunc("^/wxpay-native-notify(\\?.*)?$", orderwxpay.WxPayWebNotify)
	//testWX扫码支付
	//mux.HandleFunc("^/wxpay-native(\\?.*)?$", orderwxpay.WxPayNavite)

	//for api doc
	mux.HFunc("/api-json(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.W.Write([]byte(go_swagger.VSwagger.ToString()))
		return routing.HRES_RETURN
	})
	if s_igtest {
		mux.HFunc("/exit", exit)
	}
	mux.Handler("^/.*$", http.FileServer(http.Dir(www)))
	return mux
}
