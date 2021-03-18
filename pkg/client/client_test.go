package client

import (
	"context"
	"testing"
)

func TestClient(t *testing.T) {
	var ctx = context.Background()
	var cli = NewClient(ctx)
	var testClinets = []struct {
		Host string
		Code int
	}{
		{"https://www.baidu.com", 200},
	}

	for _, obj := range testClinets{
		var rsp, err = cli.Get(ctx, obj.Host)
		if err != nil {
			t.Errorf("request failed %v", err)
		}
		if rsp.StatusCode != obj.Code {
			t.Errorf("wanf %v got %v", obj.Code, rsp.StatusCode)
		}
	}
}
