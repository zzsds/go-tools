package aliyun

import (
	"mime/multipart"
	"testing"

	"github.com/zzsds/go-tools/pkg/storage"
)

var aliyun storage.Storage

func TestMain(t *testing.M) {
	aliyun = NewAliYun(func(o *Options) {
		o.Path = "http://xxxx.oss-cn-hangzhou.aliyuncs.com"
		o.Bucket = "xxxx"
		o.AccessKey = "xxxxx"
		o.AccessKey = "xxxxx"
		o.Endpoint = "http://oss-cn-hangzhou.aliyuncs.com"
	})
	t.Run()
}

func TestUpdate(t *testing.T) {
	resource, err := aliyun.Upload(&multipart.FileHeader{}, "test")
	if err != nil {
		t.Error(err)
	}
	t.Log(resource)
}
