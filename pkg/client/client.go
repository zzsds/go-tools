/*
 * @Author: your name
 * @Date: 2021-03-24 11:32:44
 * @LastEditTime: 2021-06-16 14:52:34
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \socloud-servere:\wuhuarou\renlanglang\zzsd\tools\pkg\client\client.go
 */
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	ContentType = "Content-Type"
	ContentJson = "application/json"
	ContentFrom = "application/x-www-form-urlencoded;charset=utf-8"
	User_Agent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586"
	Accept      = "application/x-www-form-urlencoded;charset=utf-8"
)

type Client struct {
	hc  *http.Client
	req *http.Request
}

func NewClient(ctx context.Context, opts ...OptionClient) (client *Client) {
	httpClient := initOptoionClient(opts)
	if httpClient.client == nil {
		httpClient.client = &http.Client{}
	}
	if httpClient.timeout == 0 {
		httpClient.timeout = 5
	}
	return &Client{
		hc: httpClient.client,
	}
}

func initOptoionClient(opts []OptionClient) *HttpClient {
	return &HttpClient{}
}

func (client *Client) Get(ctx context.Context, requestURL string) (*http.Response, error) {
	return client.doRequest(ctx, http.MethodGet, requestURL, "", "")
}

func (client *Client) Post(ctx context.Context, requestURL string, requestBody interface{}) (*http.Response, error) {
	return client.do(ctx, http.MethodPost, requestURL, requestBody)
}

func (Client *Client) PostForm(ctx context.Context, requestURL string, requestBody string) (*http.Response, error) {
	return Client.doRequest(ctx, http.MethodPost, requestURL, ContentFrom, requestBody)
}

func (client *Client) Delete(ctx context.Context, requestURL string) (*http.Response, error) {
	return client.doRequest(ctx, http.MethodDelete, requestURL, "", "")
}

func (client *Client) Put(ctx context.Context, requestURL string, requestBody interface{}) (*http.Response, error) {
	return client.do(ctx, http.MethodPut, requestURL, requestBody)
}

func (client *Client) do(ctx context.Context, method, requestURL string, reqBody interface{}) (*http.Response, error) {
	var reqBodyString, err = marshalBody(reqBody)
	if err != nil {
		return nil, err
	}
	return client.doRequest(ctx, method, requestURL, ContentJson, reqBodyString)
}

func (client *Client) doRequest(ctx context.Context, method, requestURL, contentType, reqBody string) (*http.Response, error) {
	var err error
	if client.req, err = http.NewRequestWithContext(ctx, method, requestURL, bytes.NewBufferString(reqBody)); err != nil {
		return nil, err
	}
	client.req.Header.Add(ContentType, contentType)
	client.req.Header.Add("Accept", Accept)
	client.req.Header.Add("User-Agent", User_Agent)
	return client.hc.Do(client.req)
}

func marshalBody(body interface{}) (string, error) {
	if body == nil {
		return "", errors.New("request params is nil")
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("Marshal request body failed %v", err)
	}
	return string(buf), nil
}
