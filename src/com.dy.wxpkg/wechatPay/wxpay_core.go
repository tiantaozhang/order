package wechatPay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Centny/gwf/log"
	"math/rand"
	"strings"
	"time"
)

func GenTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func NewOrderNo() string {
	return fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), RandInt(10000, 99999))
}

func RandInt(min int, max int) int {
	if max-min <= 0 {
		return min
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

//md5
func unifiedorderSign(kv *Kvpairs, wxPayConfig *WxPayConfig) string {

	stringA := createLinkStringNoUrl(kv)
	stringSignTemp := stringA + "&key=" + wxPayConfig.KEY
	mysign := md5String(stringSignTemp)
	return strings.ToUpper(mysign)

}

func ParseUResponse(result []byte) (*UnifiedorderResp, error) {

	uResp := &UnifiedorderResp{}
	var err error

	err = xml.Unmarshal([]byte(result), uResp)
	if err != nil {
		log.D("xml Unmarshal err is %v", err)
		return uResp, err
	}

	if uResp.Result_code == "FAIL" || uResp.Prepay_id == "" {
		return uResp, errors.New(uResp.Return_msg)
	}

	log.D("Appid is:%v,return_code:%v,prepay_id:%v", uResp.Appid, uResp.Result_code, uResp.Prepay_id)

	return uResp, nil
}

func CallbackUResponse(result []byte) (*UnifiedorderResp, error) {
	u, err := ParseUResponse(result)
	if err != nil {
		log.D("解析返回参数错误")
		return u, err
	}
	return u, err
}
