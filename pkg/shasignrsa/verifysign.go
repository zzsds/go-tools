/*
 * @Author: your name
 * @Date: 2021-04-26 10:47:24
 * @LastEditTime: 2021-04-26 11:10:09
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tools\pkg\shasignrsa\verifysign.go
 */
package shasignrsa

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
)

func VerifySign(source, sign string, pubKey interface{}, hash crypto.Hash) (bool, error) {
	var h = crypto.Hash.New(hash)
	if _, err := h.Write([]byte(source)); err != nil {
		return false, err
	}
	s, _ := base64.StdEncoding.DecodeString(sign)
	if err := rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, h.Sum(nil), s); err != nil {
		return false, err
	}
	return true, nil
}
