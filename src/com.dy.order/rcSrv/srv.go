package rcSrv

import (
	"com.dy.order/common"
	"com.dy.order/conf"
	"com.dy.order/orderList"
	"com.dy.order/orderRecord"
	"com.dy.order/orderalipay"
	"com.dy.order/orderwxpay"
	"encoding/json"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
)

type resVal struct {
	Status bool `json:"status"`
}

func CheckPurches(rc *impl.RCM_Cmd) (interface{}, error) {
	log.D("rc CheckPurches")
	var uid, pid int64
	res := resVal{}
	err := rc.ValidF(`
		uid,R|I,R:0;
		pid,R|I,R:0;
		`, &uid, &pid)
	if err != nil {
		log.E("CheckPurches arg err:%v", err)
		return res, err
	}
	dm, err := orderRecord.CheckUsrPaidItem(uid, pid)
	if err != nil {
		return res, err
	} else {
		res.Status = dm
		return res, nil
	}
}

type imss struct {
	Htm string `json:"html"`
}

func PayMakeOrder(rc *impl.RCM_Cmd) (interface{}, error) {
	log.D("rc PayMakeOrder")
	var rs imss
	var params string
	var dm string
	var payType string //AN WN
	err := rc.ValidF(`
		params,R|S,L:0;
		payType,R|S,L:0;
		`, &params, &payType)
	if err != nil {
		log.E("PayMakeOrder arg err:%v", err)
		return rs, err
	}
	if "AN" == payType {
		dm, err = orderalipay.AlipayRemoteRequest(params)
	} else if "WN" == payType {
		dm, err = orderwxpay.WxNativeRemoteCall(params)
	} else {
		err = fmt.Errorf("native payType err")
	}
	if err != nil {
		return rs, err
	} else {
		rs.Htm = dm
		return rs, nil
	}
}

func MobilePayOrder(rc *impl.RCM_Cmd) (interface{}, error) {
	log.D("rc MobilePayOrder")
	var rs imss
	var params string
	var payType string
	var dm string
	err := rc.ValidF(`
		params,R|S,L:0;
		payType,R|S,L:0;
		`, &params, &payType)
	if err != nil {
		log.E("MobilePayOrder arg err:%v", err)
		return rs, err
	}
	if "AM" == payType {
		dm, err = orderalipay.GetRsaSignJson(params)
	} else if "WM" == payType {
		dm, err = orderwxpay.WxMoblieRemoteCall(params)
	} else {
		err = fmt.Errorf("mobile payType err")
	}

	if err != nil {
		return rs, err
	} else {
		rs.Htm = dm
		return rs, nil
	}
}

func ConfirmOrderPay(rc *impl.RCM_Cmd) (interface{}, error) {
	log.D("rc CorfirmOrderPay")
	var rs imss
	var ono, payType string
	var expand string
	err := rc.ValidF(`
		ono,R|S,L:0;
		payType,R|S,L:0;
		expand,O|S,L:0;
		`, &ono, &payType, &expand)
	if err != nil {
		log.E("CorfirmOrderPay arg err:%v", err)
		return rs, err
	}
	dm, err := orderalipay.ConfirmOrderPay(ono, payType, expand)
	if err != nil {
		return rs, err
	} else {
		rs.Htm = dm
		return rs, nil
	}
}

type sisis struct {
	Htm   []orderList.ORDERt `json:"list"`
	Param paimss             `json:"param"`
}

type paimss struct {
	Total int64 `json:"total"`
	Ps    int64 `json:"ps"`
	Pn    int64 `json:"pn"`
}

func OrderList(rc *impl.RCM_Cmd) (interface{}, error) {
	log.D("rc OrderList")
	vb := imss{}
	//var pa paimss
	var ono, types string
	var uid, ps, pn, buyer, seller int64
	err := rc.ValidF(`
        ono,O|S,L:0;
        uid,R|I,R:0;
        ps,O|I,R:0;
		pn,O|I,R:0;
		type,O|S,L:0;
    	`, &ono, &uid, &ps, &pn, &types)
	if err != nil {
		log.E("OrderList arg err:%v", err)
		return vb, err
	}
	if types == "" {
		buyer = uid
		seller = 0
	} else if types == "SELL" {
		buyer = 0
		seller = uid
	}
	res, err := orderList.GetUsrOrder(buyer, seller, ono, ps, pn)
	if err != nil {
		return vb, err
	} else {
		fd, err := json.Marshal(res)
		if err != nil {
			return vb, err
		}
		vb.Htm = string(fd)
		return vb, nil
	}
}

type resres struct {
	Val int64 `json:"val"`
}

func CancelOrder(rc *impl.RCM_Cmd) (interface{}, error) {
	log.D("rc CancelOrder")
	rs := resres{}
	var ono string

	err := rc.ValidF(`
         ono,R|S,L:0;
		`, &ono)
	if err != nil {
		log.E("OrderCancel arg err :%v", err)
		return rs, err
	}
	err = orderList.DbCancelUpdate(ono)
	if err != nil {
		return rs, err
	} else {
		rs.Val = 0
		return rs, nil
	}
}

func RunSrv(args []string) {
	// netw.ShowLog = true
	// impl.ShowLog = true

	cfile := "conf/order.properties"
	if len(args) > 1 {
		cfile = args[1]
	}
	fmt.Println("Using config file:", cfile)
	err := conf.Cfg.InitWithFilePath(cfile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(common.Init("mysql", conf.ORDER_DB_CONN()))

	begin := conf.RcMemoryBegin()
	end := conf.RcMemoryEnd()
	p := pool.NewBytePool(begin, end) //memory pool.
	l, cc, cms := impl.NewChanExecListener_m_j(p, conf.RC_ADDR(), netw.NewCWH(true))
	cms.AddHFunc("check-purchse", CheckPurches)
	cms.AddHFunc("make-order", PayMakeOrder)
	cms.AddHFunc("mobile-order", MobilePayOrder)
	cms.AddHFunc("order-list", OrderList)
	cms.AddHFunc("cancel-order", CancelOrder)
	cms.AddHFunc("order-confirm", ConfirmOrderPay) //确认支付

	// cms.AddHFunc("list-bank-question", ListBankQuestion)
	// cms.AddHFunc("list-bank-paper", ListBankPaper)
	// cms.AddHFunc("update-bank-paper", UpdateBankPaper)
	// cms.AddHFunc("get-bank-paper-score", GetBankPaperScore)
	// cms.AddHFunc("get-bank-item-cnt", GetBankItemCnt)
	// cms.AddHFunc("get-usr-paper-record", GetUsrPaperRecord)
	// cms.AddHFunc("attended-bank-info", AttendBankInfo)
	cc.Run(runtime.NumCPU() - 1) //start the chan distribution, if not start, sub handler will not receive message
	err = l.Run()                //run the listen server
	if err != nil {
		panic(err.Error())
	}
	l.Wait()
}
