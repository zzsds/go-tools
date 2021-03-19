package shasignrsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
)

func ShaSignRsa(source string, privateKey *rsa.PrivateKey, hash crypto.Hash) (siginres string, err error) {
	if privateKey == nil {
		return "", errors.New("private key error")
	}
	var h = crypto.Hash.New(crypto.SHA256)
	if _, err = h.Write([]byte(source)); err != nil {
		return "", err
	}
	var resBuf []byte
	if resBuf, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, h.Sum(nil)); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(resBuf), nil
}
