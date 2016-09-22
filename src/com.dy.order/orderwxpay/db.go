package orderwxpay

import (
	"com.dy.order/common"
	order "com.dy.order/orderModel"
	"errors"
	"github.com/Centny/gwf/log"
	//"github.com/Centny/gwf/routing"
	"fmt"
	//"org.cny.uap/sync"
)

func DealWXOrder(js WXRemoteReqStruct, ono, prepayId string) (string, error) {

	//sync uap
	if err := order.SynUser(js.Seller, js.Buyer); err != nil {
		return "", err
	}
	if prepayId == "" {
		return "", errors.New("prepayId nil")
	}
	db := common.DbConn()
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("start a transaction err,%s\n", err.Error())
		return "", err
	}
	//isOnoExist, isExistErr := order.CheckIsExist(db, js.Ono, 0, "ods_order")
	isOnoExist, _ := order.CheckIsExist(db, js.Ono, 0, "ods_order")
	//integral 假定1=＝1分
	_, _needpay, _payway := order.PayTypeAndNeedPay(js.Integral, js.TotalFee)
	// _integral := float64(js.Integral) / 100.0
	// _needpay := js.TotalFee - _integral
	// _payway := 0 //0:正常支付 1:有积分有💰 2:全积分
	// js.TotalFee = _needpay
	// if 0 == _needpay {
	// 	_payway = 2
	// } else if _needpay > 0 && js.Integral > 0 {
	// 	_payway = 1
	// }
	log.D("_payway:%v", _payway)
	//======以后会有refund 请勿删除=====
	// OrderType := js.Type
	// if OrderType == "REFUND" || OrderType == "refund" {
	// 	log.D("Type Refund")
	// 	if js.Ono == "" {
	// 		log.E("Ono is nil")
	// 		tx.Rollback()
	// 		return "", errors.New("Ono is nil")
	// 	}
	// 	if isOnoExist != true {
	// 		log.E("Ono is not exist in db")
	// 		tx.Rollback()
	// 		return "", errors.New("ono is not exist:" + isExistErr.Error())
	// 	}
	// 	for i := 0; i < len(js.OrderItem); i++ {
	// 		if bl, _ := order.CheckIsExist(db, js.Ono, js.OrderItem[i].P_id, "ods_order_item"); bl != true {
	// 			log.E("ono or oid is not exist in order_item")
	// 			tx.Rollback()
	// 			return "", errors.New("ono or oid is not exist in order_item")
	// 		}
	// 	}
	// 	//更改order_item
	// 	for i := 0; i < len(js.OrderItem); i++ {
	// 		_sql := `update ods_order_item set type ='REFUNG' where ono=? and p_id=?`
	// 		_, err = tx.Exec(_sql, js.OrderItem[i].Ono, js.OrderItem[i].P_id)
	// 		if err != nil {
	// 			log.E("exec order_item refund err")
	// 			tx.Rollback()
	// 			return "", err
	// 		}
	// 	}

	// 	//insert refund
	// 	for i := 0; i < len(js.OrderRefund); i++ {
	// 		_sql := `insert into ods_order_refund(ono,item,content,imgs,status) value(?,?,?,?,?)`
	// 		_, err = tx.Exec(_sql, js.OrderRefund[i].Ono, js.OrderRefund[i].Item, js.OrderRefund[i].Imgs, js.OrderRefund[i].Status)
	// 		if err != nil {
	// 			log.E("Add OrderItem error %v", err.Error())
	// 			tx.Rollback()
	// 			return "", err
	// 		}
	// 	}
	// 	//record
	// 	//uid-->buyer  target_id-->seller
	// 	if bl, err := order.InsertRecord(tx, js.Subject, "REFUND", js.TotalFee, js.Buyer, "ALIPAY", js.Seller, "USER", ono, "NOT_PAY"); bl != true {
	// 		return "", err
	// 	}
	// 	//uid-->seller  target_id-->buyer
	// 	if bl, err := order.InsertRecord(tx, js.Subject, "REFUND", js.TotalFee, js.Seller, "ALIPAY", js.Buyer, "USER", ono, "NOT_PAY"); bl != true {
	// 		return "", err
	// 	}

	// } else {
	log.D("Type N")
	if js.TotalFee < 0 || (js.TotalFee == 0 && js.Integral == 0) {
		log.E("消费金额<=0:%v", js.TotalFee)
		return "", errors.New("TotalFee error")
	}
	//同一订单
	if !isOnoExist {
		//insert item
		log.D("length orderItem:%d", len(js.OrderItem))
		log.D("%v", js.OrderItem)
		if len(js.OrderItem) < 1 {
			log.E("OrderItem is nil")
			tx.Rollback()
			return "", errors.New("OrderItem is nil")
		}

		//detect paid_cb
		if err := order.CheckPaidcb(js.OrderItem[0].P_from); err != nil {
			log.E("p_from err")
			tx.Rollback()
			return "", errors.New("p_from err")
		}

		if err := order.InsertOrderItem(tx, order.CommonRemoteReqStruct(js), ono); err != nil {
			log.E("InsertOrderItem err:%v", err.Error())
			return "", err
		}

		// for i := 0; i < len(js.OrderItem); i++ {
		// 	_sql := `insert into ods_order_item(ono,oid,p_name,p_id,p_type,p_img,p_count,p_from,notified,price,type,status) value(?,?,?,?,?,?,?,?,?,?,?,?)`
		// 	_, err = tx.Exec(_sql, ono, js.OrderItem[i].Oid, js.OrderItem[i].P_name, js.OrderItem[i].P_id, js.OrderItem[i].P_type, js.OrderItem[i].P_img, js.OrderItem[i].P_count, js.OrderItem[i].P_from, js.OrderItem[i].Notified, js.OrderItem[i].Price, js.Type, js.OrderItem[i].Status)
		// 	if err != nil {
		// 		log.E("Add OrderItem error %v", err.Error())
		// 		tx.Rollback()
		// 		return "", err
		// 	}
		// }
		//_payway 0:正常支付 1:有积分有💰 2:全积分
		if 2 == _payway || 1 == _payway {
			//uid-->buyer  target_id-->seller
			if err := order.InsertWithIntegral(tx, order.CommonRemoteReqStruct(js), ono); err != nil {
				log.E("InsertWithIntegral: %v", err.Error())
				return "", err
			}
		}
		if 0 == _payway || 1 == _payway {
			//record
			//uid-->buyer  target_id-->seller
			// if bl, err := order.InsertRecord(tx, js.Subject, "INCOME", _needpay, js.Buyer, "WXPAY", js.Seller, "USER", ono, "NOT_PAY"); bl != true {
			// 	return "", err
			// }
			// //uid-->seller  target_id-->buyer
			// if bl, err := order.InsertRecord(tx, js.Subject, "PAY", _needpay, js.Seller, "WXPAY", js.Buyer, "USER", ono, "NOT_PAY"); bl != true {
			// 	return "", err
			// }

			if err := order.InsertWithRMB(tx, order.CommonRemoteReqStruct(js), ono, _needpay, "WXPAY"); err != nil {
				log.E("InsertWithRMB: %v", err.Error())
				return "", err
			}
		}
	}
	// }
	//insert order
	//	status := "NOT_PAY"
	if !isOnoExist {
		_sql := `insert into ods_order(ono,buyer,seller,total_price,type,status,return_url,expand,wno) value(?,?,?,?,?,?,?,?,?)`
		_, err = tx.Exec(_sql, ono, js.Buyer, js.Seller, js.TotalFee, js.Type, js.Status, js.Return_url, js.Expand, prepayId)
		if err != nil {
			log.E("Add Order error %v", err.Error())
			tx.Rollback()
			return "", err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.E("AddOrder commit error %v", err.Error())
		tx.Rollback()
		return "", err
	}
	return "", nil
}
