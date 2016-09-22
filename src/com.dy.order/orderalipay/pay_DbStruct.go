package orderalipay

import (
	om "com.dy.order/orderModel"
)

type AlipayRemoteReqStruct struct {
	// Ono      string  `json:"ono"`
	// Buyer    int64   `json:"buyer"`
	// Seller   int64   `json:"seller"`
	// Subject  string  `json:"subject"`
	// TotalFee float64 `json:"totalFee"`
	// Body     string  `json:"body"`
	// //	ShowUrl  string  `json:showUrl`
	// Type   string `json:"type"`
	// Status string `json:"status"`
	// // url  extern
	// Return_url string `json:"return_url"`
	// Expand     string `json:"expand"`
	// //integral
	// Integral int64 `json:"integral"`
	Ono      string  `json:"ono" crpy:"F"`
	Buyer    int64   `json:"buyer" crpy:"N"`
	Seller   int64   `json:"seller" crpy:"N"`
	Subject  string  `json:"subject" crpy:"N"`
	TotalFee float64 `json:"totalFee" crpy:"N"`
	Body     string  `json:"body" crpy:"N"`
	//	ShowUrl  string  `json:showUrl`
	Type   string `json:"type" crpy:"N"`
	Status string `json:"status" crpy:"N"`
	// url  extern
	Return_url string `json:"return_url" crpy:"N"`
	Expand     string `json:"expand" crpy:"N"`
	//integral
	Integral int64 `json:"integral" crpy:"N"`
	//order_item
	OrderItem []om.Item `json:"item" crpy:"F"`
	//	OrderEnv    Env      `json:env`
	OrderRefund []om.Refund `json:"refund" crpy:"F"`
}

type Order struct {
	Buyer    string  `json:"buyer"`
	Seller   string  `json:"seller"`
	TotalFee float64 `json:"totalFee"`
	//	ShowUrl  string  `json:showUrl`
	Type string `json:"type"`
}

// type Item struct {
// 	Ono      string
// 	Oid      int
// 	P_name   string
// 	P_id     int
// 	P_type   string
// 	P_count  int
// 	P_from   string
// 	Notified int
// 	Price    float64
// 	Type     string
// 	Status   string
// }

// type Refund struct {
// 	Ono     string
// 	Item    int
// 	Content string
// 	Imgs    string
// 	Status  string
// }

type Item struct {
	Ono      string  `json:"ono"`
	Oid      int64   `json:"oid"`
	P_name   string  `json:"p_name"`
	P_id     int64   `json:"p_id"`
	P_type   string  `json:"p_type"`
	P_img    string  `json:"p_img"`
	P_count  int64   `json:"p_count"`
	P_from   string  `json:"p_from"`
	Notified int64   `json:"notified"`
	Price    float64 `json:"price"`
	Type     string  `json:"type"`
	Status   string  `json:"status"`
}

type Refund struct {
	Ono     string `json:"ono"`
	Item    int64  `json:"item"`
	Content string `json:"content"`
	Imgs    string `json:"imgs"`
	Status  string `json:"status"`
}

type Env struct {
	Akey   string `json:"akey"`
	Aval   string `json:"aval"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type Record struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Money       string `json:"money"`
	Uid         string `json:"uid"`
	Pay_type    string `json:"pay_type"`
	Target_id   string `json:"target_id"`
	Target_type string `json:"target_type"`
	Ono         string `json:"ono"`
	Status      string `json:"status"`
}

type CBStruct struct {
	Code int64  `json:"code"`
	Data string `json:"data"`
}
