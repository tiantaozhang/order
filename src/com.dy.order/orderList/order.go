package orderList

import (
	"com.dy.order/common"

	//order "com.dy.order/orderalipay"
	"fmt"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
	_ "github.com/go-sql-driver/mysql"
)

type ORDERt struct {
	Tid         int64       `json:"tid"`
	Ono         string      `json:"ono"`
	Buyer       int64       `json:"buyer"`
	BuyerName   string      `json:"buyer_name"`
	Seller      int64       `json:"seller"`
	SellerName  string      `json:"seller_name"`
	Total_price float32     `json:"total_price"`
	Status      string      `json:"status"`
	Time        string      `json:"time"`
	Type        string      `json:"type"`
	Items       []orderItem `json:"items"`
}

type orderItem struct {
	Name  string  `json:"name"`
	Id    int64   `json:"id"`
	Type  string  `json:"type"`
	Img   string  `json:"img"`
	Count int64   `json:"count"`
	Price float32 `json:"price"`
}

type pageParams struct {
	Total int64 `json:"total"`
	Ps    int64 `json:"ps"`
	Pn    int64 `json:"pn"`
}

type OrderList struct {
	List   []ORDERt   `json:"list"`
	Params pageParams `json:"params"`
}

func GetUsrOrder(buyer, seller int64, ono string, ps, pn int64) (OrderList, error) {
	var ord []ORDERt
	var res OrderList
	_sql := "select a.ono,a.buyer,a.seller,a.total_price,a.type,a.status,a.time,b.usr seller_name,c.usr buyer_name from ods_order a left join ucs_usr b on b.tid=a.seller left join ucs_usr c on c.tid=a.buyer where "
	_pageSql := "select count(a.tid) from ods_order a where "
	query := " 1=1 "
	limits := ""
	orderBy := " order by a.tid desc "
	conn := common.DbConn()
	if buyer != 0 {
		query += fmt.Sprintf(" and a.buyer = %d ", buyer)
	}
	if seller != 0 {
		query += fmt.Sprintf(" and a.seller = %d ", seller)
	}
	if ono != "" {
		query += fmt.Sprintf(" and a.ono = '%s' ", ono)
	}
	if ps != 0 && pn != 0 {
		var count int64
		start := (pn - 1) * ps
		limits = fmt.Sprintf(" limit %d,%d ", start, ps)
		err := conn.QueryRow(_pageSql + query).Scan(&count)
		if err != nil {
			return res, err
		}
		log.D("ps %d,pn %d,count %d", ps, pn, count)
		res.Params.Pn = pn
		res.Params.Ps = ps
		res.Params.Total = count
	}

	rows, err := conn.Query(_sql + query + orderBy + limits)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var ono, types, status, time, buyerName, sellerName string
		var buyer, seller int64
		var total float32
		var oneOrder ORDERt
		rows.Scan(&ono, &buyer, &seller, &total, &types, &status, &time, &sellerName, &buyerName)
		oneOrder.Ono = ono
		oneOrder.Buyer = buyer
		oneOrder.BuyerName = buyerName
		oneOrder.SellerName = sellerName
		oneOrder.Seller = seller
		oneOrder.Total_price = total
		oneOrder.Status = status
		oneOrder.Type = types
		oneOrder.Time = time
		oneOrder.Items, err = GetOrderItems(ono)
		if err != nil {
			return res, err
		}
		ord = append(ord, oneOrder)
	}
	rows.Close()
	res.List = ord
	return res, nil
}

func GetOrderItems(ono string) ([]orderItem, error) {
	var res []orderItem
	conn := common.DbConn()
	_sql := "select p_name,p_id,p_type,p_img,p_count,price from ods_order_item where ono= ?"
	rows, err := conn.Query(_sql, ono)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var name, types, img string
		var id, count int64
		var price float32
		item := orderItem{}
		rows.Scan(&name, &id, &types, &img, &count, &price)
		item.Name = name
		item.Count = count
		item.Id = id
		item.Img = img
		item.Type = types
		item.Price = price
		res = append(res, item)
	}
	rows.Close()
	return res, nil
}

func DbCancelUpdate(ono string) error {
	if ono == "" {
		log.E("订单编号为空")
		return fmt.Errorf("订单编号为空")
	}
	var ty string
	db := common.DbConn()
	err := common.DbConn().QueryRow("select status from ods_order where ono = ?", ono).Scan(&ty)
	if err != nil {
		log.E("查询订单状态错误:%v", err.Error())
		return err
	}
	if ty == "NOT_PAY" {
		_, err := dbutil.DbUpdate(db, "UPDATE ods_order,ods_order_item,ods_record SET ods_order.status = 'INVALID',ods_order_item.status = 'INVALID',ods_record.status = 'INVALID' WHERE ods_order.ono = ? AND ods_order_item.ono = ? AND ods_record.ono = ?", ono, ono, ono)
		if err != nil {
			log.E("更新订单状态错误:%v", err.Error())
			return err
		}
	} else {
		log.E("订单不处于可关闭状态")
		return fmt.Errorf("订单不处于可关闭状态")
	}
	return nil
}
