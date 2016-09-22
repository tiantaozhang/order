package conf

import (
	"github.com/Centny/gwf/util"
	"strconv"
)

var Cfg *util.Fcfg = util.NewFcfg3()

func LISTEN_ADDR() string {
	if Cfg.Exist("LISTEN_ADDR") {
		return Cfg.Val("LISTEN_ADDR")
	} else {
		return ""
	}
}

func RC_ADDR() string {
	if Cfg.Exist("RC_ADDR") {
		return Cfg.Val("RC_ADDR")
	} else {
		return ""
	}
}

func ORDER_DB_CONN() string {
	if Cfg.Exist("ORDER_DB_CONN") {
		return Cfg.Val("ORDER_DB_CONN")
	} else {
		return ""
	}
}
func ORDER_TEST_DB_CONN() string {
	if Cfg.Exist("ORDER_TEST_DB_CONN") {
		return Cfg.Val("ORDER_TEST_DB_CONN")
	} else {
		return ""
	}
}

//get the UAP system URL root address.
func UrlRoot() string {
	return Cfg.Val("URL_ROOT")
}

//get the SSO login URL.
func SsoLoginUrl() string {
	return Cfg.Val("SSO_LOGIN_URL")
}

//get the SSO login URL.
func SsoLogoutUrl() string {
	return Cfg.Val("SSO_LOGOUT_URL")
}

//get the SSO authorize URL.
func SsoAuthUrl() string {
	return Cfg.Val("SSO_AUTH_URL")
}

//get the SSO login api
func SsoLoginApi() string {
	return Cfg.Val("SSO_AUTH_API")
}

func Rcp_host() string {
	return Cfg.Val("RCP_HOST")
}

func Order_host() string {
	return Cfg.Val("ORDER_HOST")
}

func RcMemoryBegin() int {
	var begin int
	begin,_ = strconv.Atoi(Cfg.Val("RC_MEMORY_BEGIN"))
	return begin
}

func RcMemoryEnd() int {
	var end int
	end,_ = strconv.Atoi(Cfg.Val("RC_MEMORY_END"))
	return end
}