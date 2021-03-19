package shasignrsa

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
	"testing"
)

var priKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAkfDV9Y1PqUYCu85TJqLfmWdawcF7r8wA2CwRWh82m3dg6L8BZo1WPJXXtLw8N0o7+YMf/lXneLZ9zd15ESU95mAALKYyniuNsGPMHm88GVQmmh7zewoyL4sV8rSNLl/D5nk7K7idb1+dotCi1T/J6ZRSzEGZuy7oWqEcCoW2vZdeg5F5DJUxzhqWidEevMTwDU3z7nppoFu+JDTz1ADbGZ0Kgv+gmH9p8iOeDmT51ylqIoO3OALjNFiIGT+Bn0c6u3IAHv8KoUWFQu0YCDfX24ZHa/hkyZO5TGaPLWOepHfcDE5nvQzxW/V99dPlsh1j7UPliA6g63zboE4TMFTRwQIDAQABAoIBAFa3MJrYHXZqSCOJpDS34H6JQA8SxUiewf2wqZrQIyVbWLTTEaT65DvZmTMmCe2caWiHtlHsfz5lyPiy2UYLx+0EK/ZbxoXfQTCHC/klhSNTsiAvteLtGwbO8Pqmt6DPfFqMvFDtQHa17LeamrZ1Uac937jIXe0wIRYA1uWVsBCUYk0MJPTZvyfWYnDaV7OXiWTWB8gfACcO0qRh567PqGQrKRkesZhxoo8Jl1V0PLLe3+L40b/UquN5wXk5CJuRvETzFN+x5fHFOArv+eJ7EBEgh2XLkqj9/us4J7cMEDnkmmPSdVz32QD7jxdtkx68JpdlUNGY8KL5SMLSddpjzPECgYEA2iOLuL1nM8+OLo5fu2MahxQw9bbkWDjqevcTbWjSZqxbG9OfqSESdcM0XFhIC3OGYeHi54SjKcCcc+z62AvEKq5MoxRbbp2oNcckHYEkQTwapZ52SwtuXwlJXDYmajz6EMH4A/nWzQURNnPSBiSGw0B8kCiwsvYHo6+7NL49Dz0CgYEAq0VZP4fCA21exh1LS15Z9Ppn0kTEwvtb0eC3moR02S6JmdNNrEc5IW9GGPBXmNV2dIjb8lyCHybTxaIGmsEdJyKJHNWU2mJg2aY2vyXYwOzsrwjA87ioqcfKf3yCd5JUOXtVt8oKYPJZ4cmlnAW9kEaQ4sSJWX23/H1r1ZkV9NUCgYAGqr1zeP51e+t5ispsPLwr0rcoW12hQKQR/Akw99ouXygtsosXrTYWOVAZXm1dRDugNDouH0SpWwStGloUTk/BijA3b8DXoaPpeNumtzK3d7HMzAoLgx7tcqg2VEVaS+DMsFD3NiSVgYkkI+gQXf9sakUkVsoHvjM/knhjRUkydQKBgG/2ruh8PFX3Oryyy9UighZHWHW6JRL+NUFX8U8fBjAwXx3jZ+SWzv9Pefi1rd4otf5qtbaTyTDKNij9yemDEybRSedCrMOzCnNeWG3PNQqyF+w5AcKSVhhflr6Oy8+VJmBJg3jZqL1F0YJsS0pa6liV+QN1zgBl0lBKQaNqJ1NdAoGBAMuWE6DLktUhUyawt4Qs6a0NCNF+6nh3cgCgmEntxQDbER0dAufDxfw/QnWhs6o8KMaLCcbEOdqNk4g04/DwLkLgHSEa+XumaZwtV5N0zvEY6ybT1fy/hZBdwtMLI7XbPgSPvQxOijeMhjGgunOwdHZoggfyMbpNlQSK4TN4KmJq
-----END RSA PRIVATE KEY-----
`
var pubKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkfDV9Y1PqUYCu85TJqLfmWdawcF7r8wA2CwRWh82m3dg6L8BZo1WPJXXtLw8N0o7+YMf/lXneLZ9zd15ESU95mAALKYyniuNsGPMHm88GVQmmh7zewoyL4sV8rSNLl/D5nk7K7idb1+dotCi1T/J6ZRSzEGZuy7oWqEcCoW2vZdeg5F5DJUxzhqWidEevMTwDU3z7nppoFu+JDTz1ADbGZ0Kgv+gmH9p8iOeDmT51ylqIoO3OALjNFiIGT+Bn0c6u3IAHv8KoUWFQu0YCDfX24ZHa/hkyZO5TGaPLWOepHfcDE5nvQzxW/V99dPlsh1j7UPliA6g63zboE4TMFTRwQIDAQAB
-----END PUBLIC KEY-----
`

func TestSha256RSA(t *testing.T) {
	var data = strings.Join([]string{"a=2", "b=3", "c=4"}, "&")
	block, _ := pem.Decode([]byte(priKey))
	if block == nil {
		t.Error("decode failed")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Errorf("Parse rsa failed %v", err)
	}
	ShaSignRsa(data, privateKey, crypto.SHA256)
	// VerifySigin(data, sign)
	// t.Log(VerifySigin(data, sign))

	// block, _ := pem.Decode([]byte(pubKey))
	// if block == nil {
	// 	return false, errors.New("public key invalid")
	// }
	// pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	// if err != nil {
	// 	return false, err
	// }

}

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
