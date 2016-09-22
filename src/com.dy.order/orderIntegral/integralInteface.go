package orderIntegral

import (
	"com.dy.order/orderModel"
	"crypto"
	"encoding/json"
	"errors"
	"github.com/Centny/gwf/log"
	"reflect"
	"time"
)

func checkParas(paras CommonRemoteReqStruct) error {
	return nil
}

//filter adn sort
func parasFilterAndSort(paras CommonRemoteReqStruct, nonstr string) (kvpairs, error) {
	pkv := kvpairs{}
	rva := reflect.ValueOf(paras)
	rty := reflect.TypeOf(paras)

	for i := 0; i < rty.NumField(); i++ {
		if tag := rty.Field(i).Tag.Get("crpy"); tag != "N" {
			continue
		}
		var value string
		switch rva.Field(i).Kind() {
		case reflect.Float64:
			if 0 == rva.Field(i).Float() {
				fmt.Println("0.0,", rva.Field(i).Float())
			}
			value = fmt.Sprintf("%f", rva.Field(i).Float())
			//与零值做比较
			if value == "0" {
				continue
			}

		case reflect.Int64:
			value = fmt.Sprintf("%d", rva.Field(i).Int())
			if value == "0" {
				continue
			}

		case reflect.String:
			value = rva.Field(i).String()
			if value == "" {
				continue
			}
		default:
			fmt.Println("type not support")
			return nil, errors.New("type not support")
		}
		pkv = append(pkv, kvpair{rty.Field(i).Tag.Get("json"), value})
	}
	pkv = append(pkv, kvpair{"nonstr", nonstr})
	fmt.Println("len pkv:", pkv.Len())
	pkv.Sort()
	return pkv, nil

}

func createString(kvs kvpairs) (string, error) {
	var str []string
	for _, v := range kvs {
		str = append(str, v.k+"=\""+v.v+"\"")
	}

	return strings.Join(str, "&"), nil
}

/*
加密用	Buyer:      267250,
		Seller:     438982,
		Subject:    "寻龙诀",
		TotalFee:   0.01,
		Body:       "迟到扣200",
		Type:       "N",
		Status:     "NOT_PAY",
		Return_url: "http://rcp.dev.jxzy.com/courseDetail.html?id=40040",
		Expand:     "id=40040&token=4d42bf9c18cb04139f918ff0ae68f8a0-fd724b48-caf7-4151-932b-dab86282ab35",
		Integral:   0,
和	nonstr随机字符串
用sha1做ras加密
1 先过滤空字符串
2 排序
3 sign
*/
func Encry(paras, nonstr string) (string, []byte, error) {
	var js CommonRemoteReqStruct
	if err := json.Unmarshal([]byte(paras), &js); err != nil {
		panic(err)
	}
	if err := checkParas(js); err != nil {
		fmt.Println(err)
		return "", nil, err
	}
	pkv, _ := parasFilterAndSort(js, nonstr)
	s, _ := createString(pkv)
	vb, err := RsaEncrypt([]byte(s))
	if err != nil {
		fmt.Println(err)
		return "", nil, err
	}
	fmt.Println("RsaEncrypt:", vb)
	return s, vb, nil
}

func Vefy(ciphertext []byte, sign []byte) error {
	return RsaVerify(ciphertext, sign)
}

/*
remote call
*/
func IComsumn(paras, nonstr, sign string) (string, error) {
	//verify sign
	s, vb, err := Encry(paras, nonstr)
	if err != nil {
		lgo.E("%s", err.Error())
		return "", nil
	}
	if err = Vefy(s, vb); err != nil {
		lgo.E("%s", err.Error())
		return "", nil
	}
	//make order
	//扣积分
	//直接写数据库

	return "", nil

}

/*
积分查询
*/
func IQuery(uid int64) string {

}

/*
single order query
*/

func SingleOrderQuery(ono string) string {
	return ""
}

/*
增加减少积分
*/
func IChange(itg int) {

}
