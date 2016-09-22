package testconf

import (
	"com.dy.order/common"
	"com.dy.order/conf"
	"fmt"
	"github.com/Centny/gwf/log"
	uap_cf "org.cny.uap/conf"
	"org.cny.uap/uap"
	"org.cny.uas/usr"
)

func init() {
	var cfg string = "../../../conf/order.properties"
	err := conf.Cfg.InitWithFilePath(cfg)
	if err != nil {
		fmt.Println(cfg)
		panic(err)
	}

	log.E(conf.ORDER_DB_CONN())
	err = common.Init("mysql", conf.ORDER_DB_CONN())

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := common.CheckDb(common.DbConn()); err != nil {
		fmt.Println(err.Error())
		return
	}

	uap.InitDb(common.DbConn)
	usr.CheckUcs(common.DbConn())
	uap_cf.Cfg = conf.Cfg
}
