package orderIntegral

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDBl55tDNr0IkllAKr6gWstlPNu75x47Vkx/p67y8M3lLHKcRwS
DA7BBtyVKcuQlSne3vm1XVuX6xkgVOnyFVNH+MIWJhjQHsd75kV0FFLU1fwBDBPj
zdQZ6aaVM2kePeyGtVfot36Oyfrxk6dl+zzaixmRuZURMZcbbVelTV/VJQIDAQAB
AoGALpJhBG7xRYXyDiBJAZacyAxrO6bdB6JhsMtGOHtebUKSOtdXH2hTLFCQRDoX
xKJ9viX6AI2C+VsPYl3LIffLXru2D+Urr35GQh0P8lt7jDxf7lZRutu8SoK6C3VT
IjqzFhDFaHIps2RnMf88S7e1MwTTayalZ23r9cCfR5+wCckCQQDucm+C4NCWs177
UQ8fv96KfWoMs44zqzJIzvSrqq7aDg7Quwhrjd0aWgfSFi6JMJDDvLsMpmo2bQKy
8UcGe7qTAkEAz9fjxJPLR69N7vsHVuCU769jW7vlBS21BAtS/YVZ9Zhf/mgNWKsC
lFcO/3Wx4EDJ07p/J4xyYfUin5+uu++sZwJAAWHwe5XKH9WSa2qg59I4/ByWDNTN
skb/16Q7jvNCaElElLlA5z6/VXPIL9OpGWqKrFffzcb5Pq+LIHZ9ru/wuQJBAJfs
FnD6FvyvByhIFXVLc5I/gUDsdtryLf5myKLHdpouZvxu0lKdraUAfdX9Eaf5s40w
JQGjh3hS1pwW/IIjDsECQQCQKuqRPJDF/5jJmvUCmwMz77IrrhcLjZGihSVoKdV6
EGV1NuuO0SHpEpC+Nu1vgW95CCak09ObrgxM67htchAu
-----END RSA PRIVATE KEY-----
`)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDBl55tDNr0IkllAKr6gWstlPNu
75x47Vkx/p67y8M3lLHKcRwSDA7BBtyVKcuQlSne3vm1XVuX6xkgVOnyFVNH+MIW
JhjQHsd75kV0FFLU1fwBDBPjzdQZ6aaVM2kePeyGtVfot36Oyfrxk6dl+zzaixmR
uZURMZcbbVelTV/VJQIDAQAB
-----END PUBLIC KEY-----
`)

var PrivK *rsa.PrivateKey
var PubK *rsa.PublicKey

func InitIntegralConfig() error {
	fmt.Println("init integral key")
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	fmt.Printf("privatekey:%v\n", string(block.Bytes))
	PrivK, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	block, _ = pem.Decode(publicKey)
	if block == nil {
		return errors.New("private key error!")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	PubK := pubInterface.(*rsa.PublicKey)
}
