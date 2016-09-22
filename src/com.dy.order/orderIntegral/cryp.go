package orderIntegral

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func RsaEncrypt(origData []byte) ([]byte, error) {
	// block, _ := pem.Decode(privateKey)
	// if block == nil {
	// 	return nil, errors.New("public key error")
	// }
	// fmt.Printf("publicKey:%v\n", string(block.Bytes))
	// priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	// if err != nil {
	// 	return nil, err
	// }
	h := sha1.New()
	h.Write(origData)
	digest := h.Sum(nil)
	fmt.Println("digest:", digest)
	s, err := rsa.SignPKCS1v15(rand.Reader, PubK, crypto.SHA1, digest)
	//	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	return s, err
}

func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	// block, _ := pem.Decode(privateKey)
	// if block == nil {
	// 	return nil, errors.New("private key error!")
	// }
	// priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	// if err != nil {
	// 	return nil, err
	// }
	return rsa.DecryptPKCS1v15(rand.Reader, PubK, ciphertext)

}

func RsaVerify(ciphertext []byte, sign []byte) error {
	// block, _ := pem.Decode(publicKey)
	// if block == nil {
	// 	return errors.New("private key error!")
	// }
	// pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	// if err != nil {
	// 	return err
	// }
	// pub := pubInterface.(*rsa.PublicKey)
	h := sha1.New()
	h.Write(ciphertext)

	err = rsa.VerifyPKCS1v15(PrivK, crypto.SHA1, h.Sum(nil), sign)
	if err != nil {
		fmt.Errorf("VerifyPKCS1v15 fail : %v\n", err)
		return err
	}
	return nil
}
