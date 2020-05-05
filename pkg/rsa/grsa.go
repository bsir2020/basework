package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	cfg "github.com/bsir2020/basework/configs"
	"io/ioutil"
)

var (
	publicKey  []byte
	privateKey []byte
)

func init() {
	var err error
	privateKey, err = ioutil.ReadFile(cfg.EnvConfig.Authkey.PrivateKey)
	if err != nil {
		panic("read private pem fail")
	}

	publicKey, err = ioutil.ReadFile(cfg.EnvConfig.Authkey.Publickey)
	if err != nil {
		panic("read public pem fail")
	}
}

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err

	}

	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
