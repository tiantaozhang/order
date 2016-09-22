package orderQr

import (
	"bytes"
	"code.google.com/p/rsc/qr"
	// "fmt"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/Centny/gwf/log"
	"image"
	"image/png"
	//"io/ioutil"
	//"os"
)

func GenQr(data string, ono string) (string, error) {
	// 	//imgName := "QR"
	if data == "" {
		log.E("data is nil")
		return "", errors.New("data is nil")
	}
	//fileName := "QR.png"

	code, err := qr.Encode(data, qr.Q)
	if err != nil {
		log.E("qr encode :%s", err.Error())
		return "", err
	}

	imgByte := code.PNG()
	img, _, _ := image.Decode(bytes.NewReader(imgByte))

	b := bytes.NewBuffer([]byte{})
	//	var b bytes.Buffer
	err = png.Encode(b, img)

	sb64 := base64.StdEncoding.EncodeToString([]byte(b.String()))

	log.D("%s\n", sb64)

	//log.D("QR code generated and saved to " + fileName)

	//{"qrbase64":"...","ono":"..."}
	mcode := map[string]interface{}{}
	mcode["code"] = 0
	mmsg := make(map[string]interface{})

	mmsg["ono"] = ono
	mmsg["qrbase64"] = sb64
	// mm = append(mm, mcode)
	// mm = append(mm, mmsg)
	mcode["data"] = mmsg

	jm, err := json.Marshal(mcode)
	if err != nil {
		return "", err
	}
	log.D("jm:%s", string(jm))

	return string(jm), nil
}

/*
func GenQr(data string, ono string) (string, error) {
	// 	//imgName := "QR"
	if data == "" {
		log.E("data is nil")
		return "", errors.New("data is nil")
	}
	fileName := "QR.png"

	code, err := qr.Encode(data, qr.Q)
	if err != nil {
		log.E("qr encode :%s", err.Error())
		return "", err
	}

	imgByte := code.PNG()
	img, _, _ := image.Decode(bytes.NewReader(imgByte))

	out, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.E("create qr err,%s", err.Error())
		return "", err
	}

	err = png.Encode(out, img)
	if err != nil {
		log.E("png :%s", err.Error())
		return "", err
	}
	got := []byte{}
	fpng, _ := os.OpenFile("QR.png", os.O_RDONLY, 0666)
	if got, err = ioutil.ReadAll(fpng); err != nil {
		log.E("read png :%s", err.Error())
		return "", err
	}
	defer fpng.Close()
	sb64 := base64.StdEncoding.EncodeToString(got)

	log.D("%s\n", sb64)

	log.D("QR code generated and saved to " + fileName)

	//{"qrbase64":"...","ono":"..."}
	mcode := map[string]interface{}{}
	mcode["code"] = 0
	mmsg := make(map[string]interface{})

	mmsg["ono"] = ono
	mmsg["qrbase64"] = sb64
	// mm = append(mm, mcode)
	// mm = append(mm, mmsg)
	mcode["data"] = mmsg

	jm, err := json.Marshal(mcode)
	if err != nil {
		return "", err
	}
	log.D("jm:%s", string(jm))

	return string(jm), nil
}
*/
