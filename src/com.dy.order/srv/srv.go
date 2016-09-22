package srv

import (
	// "com.dy.order/cart"
	"com.dy.order/common"
	"com.dy.order/conf"
	"com.dy.order/orderalipay"
	"com.dy.order/orderwxpay"
	"com.dy.order/rcSrv"
	"encoding/json"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tomlovzki/go_swagger"
	"net/http"
	uap_cf "org.cny.uap/conf"
	"org.cny.uap/uap"
	"org.cny.uas/usr"
	"os"
	"sync"
	// "time"
)

var lock sync.WaitGroup
var s_running bool
var s http.Server
var s_igtest bool = false

var _ = go_swagger.NewH().AddPath(
	"/t",
	map[string]go_swagger.Path{
		"get": go_swagger.Path{
			Tags:        []string{"order"},
			Summary:     "测试接口文档",
			Description: "测试接口文档",
			//	Summary:     "测试接口文档",
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Parameters: []go_swagger.Parameters{
				go_swagger.Parameters{
					In:          "query",
					Name:        "l",
					Description: "测试接口文档",
					Required:    false,
					Type:        "int64",
				},
			},
			Responses: map[string]interface{}{
				"200": map[string]interface{}{
					"description": "成功",
					"schema": map[string]string{
						"$ref": "#/definitions/t-rs",
					},
				},
			},
		}}).AddDefinition("t-rs", go_swagger.Definition{
	Type: "object",
	Properties: func() (m map[string]go_swagger.Properties) {
		json.Unmarshal([]byte(`
{
    "_id": {
        "type": "object",
        "example": "dd",
        "description": "测试"
    },
    "word_history": {
        "type": "object",
        "example": [
            {
                "word": "dd",
                "time": 1446694599105334300
            },
            {
                "word": "qq",
                "time": 1446695490038232800
            }
        ],
        "description": "word ->搜索词 , time -> 搜索时间"
    }
}
			`), &m)
		return
	}(),
	Xml: map[string]string{"name": "courseDetail"},
})

func run(args []string) {
	cfg := "conf/order.properties"
	if len(args) > 1 {
		cfg = args[1]
	}
	err := conf.Cfg.InitWithFilePath(cfg)
	//conf.Cfg.Print()
	if err != nil {
		fmt.Println(cfg)
		panic(err)
	}

	err = common.Init("mysql", conf.ORDER_DB_CONN())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := common.CheckDb(common.DbConn()); err != nil {
		fmt.Println(err.Error())
		return
	}
	//
	orderalipay.InitAlipayConfig()
	orderwxpay.InitWxConfig()
	uap.InitDb(common.DbConn)
	usr.CheckUcs(common.DbConn())
	uap_cf.Cfg = conf.Cfg
	//	fmt.Println(uap_cf.UsrApiRoot(), "---->")

	defer StopSrv()
	go_api()
	mux := http.NewServeMux()

	mux.Handle("/", NewSrvMux("", "www"))
	log.D("running server on %v", conf.LISTEN_ADDR())
	s = http.Server{Addr: conf.LISTEN_ADDR(), Handler: mux}
	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func go_api() {
	sinfo := `{"description": "order接口文档",
    "version": "1.0.0",
    "title": "order接口文档"
    }`
	infoParse := map[string]interface{}{}
	json.Unmarshal([]byte(sinfo), &infoParse)
	go_swagger.NewH().InitSwagger(go_swagger.Swagger{
		SwaggerVersion: "2.0",
		Info:           infoParse,
		Host:           "order.dev.jxzy.com",
		BasePath:       "",
		Schemes:        []string{"http"},
	}).AddTag(
		go_swagger.Tag{
			Name:        "order",
			Description: "订单相关",
		}).AddTag(
		go_swagger.Tag{
			Name:        "app",
			Description: "应用相关",
		}).AddTag(
		go_swagger.Tag{
			Name:        "balance",
			Description: "余额相关",
		})
}

// func timer() {
// 	timer1 := time.NewTicker(1 * time.Hour)
// 	for {
// 		select {
// 		case <-timer1.C:
// 			if time.Now().Hour() == 4 {
// 				cart.DeleteCart()
// 			}
// 		}
// 	}
// }

//run the server.
func RunSrv(args []string) {
	s_running = true
	lock.Add(1)
	go run(args)
	go rcSrv.RunSrv(os.Args)
	lock.Wait()
	s_running = false
}

//stop the server.
func StopSrv() {
	if s_running {
		lock.Done()
	}
}

func exit(hs *routing.HTTPSession) routing.HResult {
	log.D("receiving exit command...")
	StopSrv()
	return hs.MsgRes("SUCCESS")
}
