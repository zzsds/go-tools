package yunpian

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strconv"

	ypclnt "github.com/yunpian/yunpian-go-sdk/sdk"
	"github.com/zzsds/kratos-tools/pkg/sms"
)

// YunPian ...
type YunPian struct {
	opts   Options
	client ypclnt.YunpianClient
}

// NewYunPian ...
func NewYunPian(opts ...Option) sms.Server {
	yun := YunPian{
		opts: newOptions(opts...),
	}
	yun.Init()
	return &yun
}

func (m *YunPian) Init(opts ...Option) {
	if m.opts.Key == "" {
		log.Fatal("云片 ApiKey 不能为空")
	}
	m.client = ypclnt.New(m.opts.Key)
}

// Debug ...
func (m *YunPian) Debug() bool {
	return m.opts.Debug
}

// String ...
func (m *YunPian) String() string {
	return "yunpian"
}

// Send ...
func (m *YunPian) Send(mobile, content string) (string, error) {
	param := ypclnt.NewParam(2)
	param[ypclnt.MOBILE] = mobile
	param[ypclnt.TEXT] = content
	if m.opts.Debug {
		return "0", nil
	}
	r := m.client.Sms().SingleSend(param)
	if r.Code != 0 {
		return "", errors.New(r.Msg + ": " + r.Detail)
	}

	var result struct {
		Count  int
		Fee    float64
		Mobile string
		Sid    int
		Unit   string
	}
	dByte, _ := json.Marshal(r.Data)
	json.Unmarshal(dByte, &result)

	return strconv.Itoa(result.Sid), nil
}

// Assignment ...
func (m *YunPian) Assignment(content string, value ...string) string {
	re := regexp.MustCompile("#([a-z|A-Z]*)#")
	var i int
	return re.ReplaceAllStringFunc(content, func(s string) string {
		for k, v := range value {
			if i == k {
				s = v
			}
		}
		i++
		return s
	})
}
