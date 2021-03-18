package client

import (
	"net/http"
	"time"
)

type HttpClient struct {
	client  *http.Client
	timeout time.Duration
}

type withHTTPClient struct{ client *http.Client }

func (w withHTTPClient) Apply(o *HttpClient) {
	w.client = o.client
}

func WithHTTPClient(client *http.Client) OptionClient {
	return withHTTPClient{client}
}

type OptionClient interface {
	Apply(*HttpClient)
}

type withTimeout time.Duration

func (w withTimeout) Apply(o *HttpClient) {
	o.timeout = time.Duration(w)
}

func WithTimeOut(timeout time.Duration) OptionClient {
	return withTimeout(timeout)
}
