/*
 * @Author: your name
 * @Date: 2021-03-24 11:32:44
 * @LastEditTime: 2021-06-18 13:29:21
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tools\pkg\shasignrsa\sendsign.go
 */
package shasignrsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
)

func SendSign(source string, privateKey *rsa.PrivateKey, hash crypto.Hash) (siginres string, err error) {
	if privateKey == nil {
		return "", errors.New("private key error")
	}
	var h = crypto.Hash.New(hash)
	if _, err = h.Write([]byte(source)); err != nil {
		return "", err
	}
	var resBuf []byte
	if resBuf, err = rsa.SignPKCS1v15(rand.Reader, privateKey, hash, h.Sum(nil)); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(resBuf), nil
}
