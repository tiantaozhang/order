package util

import (
	"encoding/gob"
)

func init() {
	gob.Register(Order{})
	gob.Register([]Order{})
}

type Order struct {
	Tid        int64       `m2s:"TID" json:"tid"`
	OrderNo    string      `m2s:"ORD_NO" json:"orderNo"`
	AppId      string      `m2s:"RES_APP" json:"appId"`
	Seller     int64       `m2s:"SELLER" json:"seller"`
	SellerName string      `m2s:"SELLER_NAME" json:"sellerName"`
	Buyer      int64       `m2s:"BUYER" json:"buyer"`
	BuyerName  string      `m2s:"BUYER_NAME" json:"buyerName"`
	TotalPrice float64     `m2s:"TOTAL_PRICE" json:"totalPrice"`
	Price      float64     `m2s:"PRICE" json:"price"`
	OrderType  string      `m2s:"ORDER_TYPE" json:"orderType"`
	Channel    string      `m2s:"CHANNEL" json:"channel"`
	Discount   string      `m2s:"DISCOUNT" json:"discount"`
	PayTime    string      `m2s:"PAYTIME" json:"paytime"`
	PayWay     string      `m2s:"PAYWAY" json:"payway"`
	Time       string      `m2s:"TIME" json:"time"`
	Status     string      `m2s:"STATUS" json:"status"`
	Items      []OrderItem `json:"items"`
	UserIncome []UserIncome
}

type OrderItem struct {
	Tid           int64   `m2s:"TID" json:"tid"`
	OrderId       int64   `m2s:"ORDER_ID" json:"orderId"`
	OrderNo       string  `m2s:"ORD_NO" json:"orderNo"`
	GoodsId       string  `m2s:"GOOD_ID" json:"goodId"`
	GoodsDetail   string  `m2s:"GOOD_DETAIL" json:"goodDetail"`
	Price         float64 `m2s:"PRICE" json:"price"`
	TotalPrice    float64 `m2s:"TOTAL_PRICE" json:"totalPrice"`
	DiscountPrice float64 `m2s:"DISCOUNT_PRICE" json:"discountPrice"`
	Buyer         int64   `m2s:"BUYER" json:"buyer"`
	Seller        int64   `m2s:"SELLER" json:"seller"`
	Count         int64   `m2s:"COUNT" json:"count"`
	CommentStatus string  `m2s:"COMMENT_STATUS" json:"commentStatus"`
	Remark        string  `m2s:"REMARK" json:"remark"`
	Time          string  `m2s:"TIME" json:"time"`
	Status        string  `m2s:"STATUS" json:"status"`
}

//discount
type UserIncome struct {
	UserId int64   `json:"userId"`
	Money  float64 `json:"money"`
}
type OrderDiscount struct {
	Order       []Order
	Discount    string `m2s:"DISCOUNT" json:"discount"`
	UseDiscount bool   `json:"useDiscount"`
}

//cart
type OdsCart struct {
	Tid         int64   `json:"tid" m2s:"tid,autoinc"`
	AppId       string  `json:"appId" m2s:"app_id"`
	AppName     string  `json:"appName" m2s:"app_name"`
	GoodsId     string  `json:"goodsId" m2s:"goods_id"`
	GoodsDetail string  `json:"goodsDetail" m2s:"goods_detail"`
	Seller      int64   `json:"seller" m2s:"seller"`
	SellerName  string  `json:"sellerName" m2s:"seller_name"`
	Buyer       int64   `json:"buyer" m2s:"buyer"`
	BuyerName   string  `json:"buyerName" m2s:"buyer_name"`
	Count       int64   `json:"count" m2s:"count"`
	Price       float64 `json:"price" m2s:"price"`
	Time        string  `json:"time" m2s:"time"`
	Status      string  `json:"status" m2s:"status"`
	RedictUrl   string  `json:"redictUrl" m2s:"redict_url"`
}
