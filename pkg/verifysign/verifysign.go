package verifysign

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
)

func VerifySigin(source, sign string, pubKey *rsa.PublicKey, hash crypto.Hash) (bool, error) {
	var h = crypto.Hash.New(hash)
	if _, err := h.Write([]byte(source)); err != nil {
		return false, err
	}
	s, _ := base64.StdEncoding.DecodeString(sign)
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, h.Sum(nil), s); err != nil {
		return false, err
	}
	return true, nil
}
