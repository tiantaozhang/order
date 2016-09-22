package main

import (
	"com.dy.aliPkgT/alipay"
	"errors"
	"fmt"
	//"io"
	//"log"
	"io/ioutil"
	"net/http"
	"net/url"
	//"regexp"
	"github.com/Centny/gwf/util"
	"strings"
	"sync"
	"time"
)

type Times struct {
	Cnt  int
	Lock sync.RWMutex
}

func (ts *Times) AddOne() {
	ts.Lock.Lock()
	defer ts.Lock.Unlock()
	ts.Cnt++
}

func (ts *Times) Reset() {
	ts.Lock.Lock()
	defer ts.Lock.Unlock()
	ts.Cnt = 0
}

var ts *Times = new(Times)

/*
http://192.168.2.231:8888
/alipay-web-return
?body=test+web+Chinese&
buyer_email=327468120%40qq.com&
buyer_id=2088502994384781&
exterface=create_direct_pay_by_user&
is_success=T&
notify_id=RqPnCoPT3K9%252Fvwbh3InVZ3r3Q8XImQft4yHdn52q1w4C8YioyZQ8YGGBWRFLjAeogS0l&
notify_time=2015-12-04+12%3A16%3A11&notify_type=trade_status_sync&
out_trade_no=2015120412150153694&payment_type=1&
seller_email=itdayang%40gmail.com&
seller_id=2088501949844011&
subject=迟到扣200&
total_fee=0.01&
trade_no=2015120421001004780215834730&
trade_status=TRADE_SUCCESS&
sign=c0a5dbab3a9e707367fc8aad962a5288&
sign_type=MD5

*/

func ImitateAliReturn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request from QC")
	r.ParseForm()
	got, values, err := alipay.GenAndNotify(r, alipay.AWebConfig)
	if err != nil {
		fmt.Printf("\n", err.Error())
		return
	}
	go func(strNotifyUrl string, value url.Values) {
		var n time.Duration
	L:
		for i := 1; i <= 5; i++ {
			timer := time.NewTimer(n * time.Second)
			select {
			case <-timer.C:
				if err := imitateAliNotify(strNotifyUrl, value); err != nil {
					n = n + time.Duration(i)
				} else {
					break L
				}

			}
		}

	}(alipay.AWebConfig.Notify_url, values)

	fmt.Fprintf(w, "%s", string(got))
}

func imitateAliNotify(strNotifyUrl string, value url.Values) error {
	res, err := http.PostForm(strNotifyUrl, value)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}
	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}
	if string(got) != "success" {
		fmt.Printf("-------string got failed-------")
		return errors.New("string not equal success")
	}
	ts.AddOne()
	return nil
}

func TotalSuccessTimes(w http.ResponseWriter, r *http.Request) {
	ts.Lock.RLock()
	defer ts.Lock.RUnlock()
	fmt.Fprintf(w, "%d", ts.Cnt)
}

func ResetTimes(w http.ResponseWriter, r *http.Request) {
	ts.Reset()
}
func ImitaAliVerify(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("true"))
}

func SimulateAlipay() {

}

//
func TestAlipay(w http.ResponseWriter, r *http.Request) {

	args, _ := alipay.GenTestData()

	args_ := ""
	for k, v := range args {
		args_ += k + "=" + v + "&"
	}
	_args := strings.TrimRight(args_, "&")
	fmt.Println("https://www.alipay.com/cooperate/gateway.do?" + _args)
	res, err := util.HPost("https://www.alipay.com/cooperate/gateway.do", args)

	//res, err := http.Get("https://www.alipay.com/cooperate/gateway.do?" + _args)
	//defer res.Close()
	//got, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Printf("%s\n", err.Error())
	// }
	if err != nil {

	} else {
		//fmt.Printf("%s\n", res)
		w.Write([]byte(res))
	}
	return
	//	return common.MsgRes(hs, parse)
}
