package tools

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"
)

const (
	five = 5
)

// CallerHeader define struct fir api call header
type CallerHeader struct {
	Key   string
	Value string
}

// CallAPI is a function for api call
func CallAPI(httpMethod, url string, headers []CallerHeader, payload interface{}, queryParams map[string]string) (*http.Response, error) {
	var req *http.Request
	var err error
	if payload != nil && httpMethod != "GET" {
		payloadBuf := new(bytes.Buffer)
		if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
			return nil, err
		}
		req, err = http.NewRequest(httpMethod, url, payloadBuf)
	} else {
		req, err = http.NewRequest(httpMethod, url, nil)
	}

	if err != nil {
		return nil, err
	}

	if httpMethod == http.MethodGet {
		q := req.URL.Query()
		for i, v := range queryParams {
			q.Add(i, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Content-Type", "application/json")

	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}

	// ignore expired SSL certificates
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint
	}

	client := &http.Client{
		Timeout:   time.Minute * five,
		Transport: transCfg,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
