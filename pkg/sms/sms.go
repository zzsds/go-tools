package sms

import (
	"bytes"
)

// Driver ...
type Server interface {
	Assignment(content string, value ...string) string
	Send(mobile, content string) (string, error)
	Debug() bool
	String() string
}

// SignContent ...
func SignContent(sign, content string) string {
	var buff bytes.Buffer
	buff.WriteString("【")
	buff.WriteString(sign)
	buff.WriteString("】")
	buff.WriteString(content)
	return buff.String()
}
