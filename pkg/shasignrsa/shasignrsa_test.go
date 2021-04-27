/*
 * @Author: your name
 * @Date: 2021-03-24 11:32:44
 * @LastEditTime: 2021-04-26 11:10:27
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tools\pkg\shasignrsa\shasignrsa_test.go
 */
package shasignrsa

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"
)

var priKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAimFB73SETSDJWBOLXlFj8073mJTqHBHgdquOyScd2MJNdbUD+DWWtJ6u6LO18oahCVjTmWGUO8xAz9BwzTZ6TM1OrLt8Y1D9ybJ0PSljqYavb1u0qhMUTJOknInHKsZZ+f49De7MkKna9LUWKInqm12QrX2ge2aAxfDThPvQ1x7IctIY4g1yCCC0UIYfvMaQIGN1p9SWD6VAVA9hUDljOX7dhZD3lsyLT6vZvZl6YfoWEomA34yTmR5J9zXZNJrXTHcelCWr1Y/v5EBUzJ8ATzRLJxLLocDXSabw+c6ITDFF326eckx74Giqn+qo6QVf8gy9t7orYswcETpbximPRwIDAQABAoIBAGk9glnMcXn+/2G+q3W1zKAUZHVke4+RgPZ/jv4og6iATUzBuB0jFFSVgkxzsGKcRQjIx1SVQ5kexAPIcKGBVw3l7nmrtPQLepfU4lZJjgQ17GJyijn2fK+ocb6jghdj9rYLxv87p9Q1edI1jn6SNRyn0go/yrdOw4zGlPdEVBXlphSJssmfV4AqyEmJK5+o41AGr5BBopxAYx5z3a+/6HkWOVozB0hX1TOsCxXyXLPqpb+dMVEoRnZ21sNa3aVh/WFzIkAyA9vkB73zmsMWakZbAxapRWy3oFM98ulmkrSegIOBuNRl7ddx3ouOu7q5P6+zjpi3Lq0SWU+h/+mt4gECgYEA/j61OT2Oqf7WhG9znlNv5XRNBX5GvHqbcdgB5wH5X4XBSs+jbrCrKwZJ29FyvTjl1YmhYEjeFjtcl7lFp1nf7lHgY2/qLkr23RVi982tqRAU1MzP8Emq1hNLlMb64FRW29tenQkRn3PI1zulqwqla3JjxiA2X5Zo4MSyn4if7McCgYEAi1XMHW9Vp6uX3XOKMbvn7LBKAgwGX0ykpYjD0C93L8vy+jQaQOHHhnYvUv/BBybNTUq2Yh1lkUU9ZrUhrAtHlmbD87HjZLvgIUWTnRULA9SKHSXxESlchRNQ7Cug4Pll6ylEpnABf5XEWsvCCsDLnl6WKsHw9FJq3LMvu8T9yYECgYARvpKrYg548t5J8/Vf0Xb3zrwpa/zH3s7GjUrkspCTCCTLcd54NUBdCl1RSDb32ebAlpB6xdsqNg5qUHX6Dh3A5loA1qjDflvoZju4C9TY/dRWXc0NejbAJiyaP6D20ywUwCTEVQOz20LjMriHTYDqFNu90jW5Sigbt963n7N5tQKBgEsoeZS7FHIAHkfm5flTyZOjuBgIkntfZUShVDZ9FAZlNeh+qFatMyo1n+teZ6nK5V022tBr2PiWZ7t6IvKhlvjq7/II14bjzM9Fr41A55MmV2XHrJQ8QlrKA5GRKxOPk8lYll5M9pHyoFr1o/KW8n63uLrRqH6x4lCwGyRm6xqBAoGBAPhXAINL9xbxms6+1kohQm8193eVt1ikv5DBT5+LHxdIkQPObmnliOaDuOMk0o0mhnrUo46pEslsNjBGNb/Cbc+wVH7s94jNd2mFWfr6A9kWX1cC5g164gx8/d+/tXPaVHdkIQ7piAOUg2xLq9Zr2xAEXKJXBNH6rnhVnFza/Vt8
-----END RSA PRIVATE KEY-----
`
var pubKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAimFB73SETSDJWBOLXlFj
8073mJTqHBHgdquOyScd2MJNdbUD+DWWtJ6u6LO18oahCVjTmWGUO8xAz9BwzTZ6
TM1OrLt8Y1D9ybJ0PSljqYavb1u0qhMUTJOknInHKsZZ+f49De7MkKna9LUWKInq
m12QrX2ge2aAxfDThPvQ1x7IctIY4g1yCCC0UIYfvMaQIGN1p9SWD6VAVA9hUDlj
OX7dhZD3lsyLT6vZvZl6YfoWEomA34yTmR5J9zXZNJrXTHcelCWr1Y/v5EBUzJ8A
TzRLJxLLocDXSabw+c6ITDFF326eckx74Giqn+qo6QVf8gy9t7orYswcETpbximP
RwIDAQAB
-----END PUBLIC KEY-----
`

func TestSha256RSA(t *testing.T) {
	var data = strings.Join([]string{"a=2", "b=3", "c=4"}, "&")
	block, _ := pem.Decode([]byte(priKey))
	if block == nil {
		t.Error("decode failed")
	}
	fmt.Println(data)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Errorf("Parse rsa failed %v", err)
	}
	fmt.Println(SendSign(data, privateKey, crypto.SHA256))
	var _, _ = SendSign(data, privateKey, crypto.SHA256)
	block, _ = pem.Decode([]byte(pubKey))
	if block == nil {
	}
	fmt.Println(len(block.Bytes))
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {

	}
	t.Log(VerifySign(data, "C9s1PRX+RdSOtr2C77cJ9TLIqLCXPeWXkniL6BHVUwlNRdgJXlE33DayytD6r/WqksYqGAIT0JWBuLQxh9VdrtgTlGYqabLrg3orUD5vV9cPbAXk7GTJ6kAgBr3wL6Ub+IaKqgXaYRxhUx8hR56yv8a4WIy4JZZW1Thj3DJDCRcsAtPjYqdKOxEhfyLmoQkkr5q9kwNXea13O+ykO2tFmIuekz/4cjrpLmNa3sVdj6qJmijbzvKkJ89+/SUlDI6iKVtbRLDDafwiDIDPk7pG0hb/D+A9MQow38+4eSYJiKpK71M5zr5IMGpoHXXSVHomVoZhhI/LWK8fi3uRaKkubw==", pub, crypto.SHA256))
}
