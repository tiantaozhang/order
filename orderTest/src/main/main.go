package main

import (
	"demo"
	"fmt"
	//"github.com/Centny/gwf/routing"
	"net/http"
)

func main() {
	demo.InitAlipayConfig()
	port := "7491"
	http.HandleFunc("/alipay-web-notify", demo.AlipayWebNotify)
	http.HandleFunc("/alipay-web-return", demo.AlipayWebReturn)
	http.HandleFunc("/alipay-test", demo.AlipayWebTest)
	http.HandleFunc("/notify_query.do", ImitaAliVerify)
	http.HandleFunc("/AliPayRequest", ImitateAliReturn)
	http.HandleFunc("/AliPayVerify", ImitaAliVerify)
	http.HandleFunc("/TotalCnt", TotalSuccessTimes)
	http.HandleFunc("/ResetCnt", ResetTimes)
	//mux.HFunc("^/testPay(\\?.*)?$", orderalipay.TestAlipay)
	http.HandleFunc("/testPay", TestAlipay)
	s := &http.Server{Addr: ":" + port}
	fmt.Println("==========ready to service===========")
	fmt.Println(s.ListenAndServe())
}

/*
http://192.168.2.231:8888/alipay-web-return?body=test+web+Chinese&buyer_email=327468120%40qq.com&buyer_id=2088502994384781&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%252Fvwbh3InVZ3r3Q8XImQft4yHdn52q1w4C8YioyZQ8YGGBWRFLjAeogS0l&notify_time=2015-12-04+12%3A16%3A11&notify_type=trade_status_sync&out_trade_no=2015120412150153694&payment_type=1&seller_email=itdayang%40gmail.com&seller_id=2088501949844011&subject=迟到扣200&total_fee=0.01&trade_no=2015120421001004780215834730&trade_status=TRADE_SUCCESS&sign=c0a5dbab3a9e707367fc8aad962a5288&sign_type=MD5

*/

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

/*
http://order.dev.jxzy.com/alipay-web-return?buyer_email=327468120%40qq.com&buyer_id=2088502994384781&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%252Fvwbh3InVZ3%252B1VEgf7tEq2H%252FnYoPlP8Xrd%252Fw1o3XZM3kMWtRtHIjue8%252FU&notify_time=2015-12-10+12%3A02%3A48&notify_type=trade_status_sync&out_trade_no=2015121011563647474&payment_type=1&seller_email=itdayang%40gmail.com&seller_id=2088501949844011&subject=%E6%B5%8B%E8%AF%95%E8%AF%BE%E7%A8%8B&total_fee=0.02&trade_no=2015121021001004780086872099&trade_status=TRADE_SUCCESS&sign=3da88dec24fe572c74a0278e90c833b0&sign_type=MD5
*/

/*

http://192.168.2.231:8888/alipay-web-return?body=test+web+Chinese&buyer_email=327468120%40qq.com&buyer_id=2088502994384781&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%252Fvwbh3InVZ3%252B1UCLNXOpgg3GYBagm%252Flw0hPdyfyt%252BR9m%252BCoqeffCJPZ1V&notify_time=2015-12-10+13%3A20%3A26&notify_type=trade_status_sync&out_trade_no=2015121013184489812&payment_type=1&seller_email=itdayang%40gmail.com&seller_id=2088501949844011&subject=%E8%BF%9F%E5%88%B0%E6%89%A3200&total_fee=0.01&trade_no=2015121021001004780087375931&trade_status=TRADE_SUCCESS&sign=04bacf398b2835f4e098b7c39d464f80&sign_type=MD5
*/

/*

http://order.dev.jxzy.com/alipay-web-return?buyer_email=327468120%40qq.com&buyer_id=2088502994384781&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%252Fvwbh3InVZ3%252B1VEgf7tEq2H%252FnYoPlP8Xrd%252Fw1o3XZM3kMWtRtHIjue8%252FU&notify_time=2015-12-10+12%3A02%3A48&notify_type=trade_status_sync&out_trade_no=2015121011563647474&payment_type=1&seller_email=itdayang%40gmail.com&seller_id=2088501949844011&subject=%E6%B5%8B%E8%AF%95%E8%AF%BE%E7%A8%8B&total_fee=0.02&trade_no=2015121021001004780086872099&trade_status=TRADE_SUCCESS&sign=3da88dec24fe572c74a0278e90c833b0&sign_type=MD5

*/

/*
http://rcp.dev.jxzy.com/courseDetail.html?id=39940&body=id%3D39940&buyer_email=327468120%40qq.com&buyer_id=2088502994384781&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%252Fvwbh3InVZ3%252B1U9KaCyJMg6MmOxu6LVxkM2zLWSD2TplSk4jVaxwbp5F4&notify_time=2015-12-10+14%3A09%3A22&notify_type=trade_status_sync&out_trade_no=2015121014074845527&payment_type=1&seller_email=itdayang%40gmail.com&seller_id=2088501949844011&subject=%E8%BF%9F%E5%88%B0%E6%89%A3200&total_fee=0.01&trade_no=2015121021001004780092665865&trade_status=TRADE_SUCCESS&sign=d7b5f9d7dfec5ff4e771a1dd64c55db4&sign_type=MD5#/
*/

/*
MD5:f7505bd3c662b0b7b812786169eb755a
sign:818c5ea4f4dc66d344915c5d0039171a
/alipay-web-return?
_input_charset=utf-8&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%252Fvwbh3InVZ3r3Q8XImQft4yHdn52q1w4C8YioyZQ8YGGBWRFLjAeogS0l&notify_time=2015-12-14+12%3A16%3A11&out_trade_no=2015121118292420911&partner=2088501949844011&payment_type=1&seller_email=itdayang@gmail.com&service=create_direct_pay_by_user&subject=订单测试第一单&total_fee=0.02&sign=818c5ea4f4dc66d344915c5d0039171a&sign_type=MD5

_input_charset=utf-8&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9%2Fvwbh3InVZ3r3Q8XImQft4yHdn52q1w4C8YioyZQ8YGGBWRFLjAeogS0l&notify_time=2015-12-14 12:16:11&out_trade_no=2015121118292420911&partner=2088501949844011&payment_type=1&seller_email=itdayang@gmail.com&service=create_direct_pay_by_user&subject=订单测试第一单&total_fee=0.02


_input_charset=utf-8&exterface=create_direct_pay_by_user&is_success=T&notify_id=RqPnCoPT3K9252Fvwbh3InVZ3r3Q8XImQft4yHdn52q1w4C8YioyZQ8YGGBWRFLjAeogS0l&notify_time=2015-12-14 12:16:11&out_trade_no=2015121118292420911&partner=2088501949844011&payment_type=1&seller_email=itdayang@gmail.com&service=create_direct_pay_by_user&subject=订单测试第一单&total_fee=0.02&sign=57176f1211bc76902f27c24a736ae587&sign_type=MD5
sign=e345b1d5452bf07ebc36345b11228084


_input_charset=utf-8&out_trade_no=2015121118292420911&partner=2088501949844011&payment_type=1&seller_email=itdayang@gmail.com&service=create_direct_pay_by_user&subject=订单测试第一单&total_fee=0.02


*/
