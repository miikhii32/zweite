package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

var SharedClient *http.Client

func init() {
	jar, _ := cookiejar.New(nil)
	SharedClient = &http.Client{
		Jar: jar,
	}
}

func Request(requestedURL string, cookie string, addHeaders bool) (*http.Response, error) {
	req, err := http.NewRequest("GET", requestedURL, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, errors.New("failed to create request")
	}

	if cookie != "" {
		headerReader := http.Header{}
		headerReader.Add("Cookie", cookie)
		dummyReq := &http.Request{Header: headerReader}
		for _, c := range dummyReq.Cookies() {
			req.AddCookie(c)
		}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	if addHeaders {
		req.Header.Add("Referer", requestedURL)
		req.Header.Add("Origin", requestedURL)
		req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	}

	res, err := SharedClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to fetch URL: %v\n", err)
		return nil, errors.New("failed to fetch URL")
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}