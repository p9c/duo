// Package rpc is an RPC interface to a Parallelcoin full node wallet RPC
package rpc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// NewClient creates a new RPC client
func NewClient(host string, port int, user, passwd string, useTLS bool) *Client {
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	if port == 0 || port < 1024 {
		port = 11048
	}
	var URL string
	var httpClient *http.Client
	if useTLS {
		URL = "https://"
		t := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: t}
	} else {
		URL = "http://"
		httpClient = &http.Client{}
	}
	return &Client{URL: fmt.Sprintf("%s%s:%d", URL, host, port), Username: user, Password: passwd, httpClient: httpClient}
}

// DoTimeoutRequest process a HTTP request with timeout
func (c *Client) DoTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("Timeout reading data from server")
	}
}

// Call prepare & exec the request
func (c *Client) Call(method string, params interface{}) (rr Response, err error) {
	connectTimer := time.NewTimer(RPCClientTimeout * time.Second)
	rpcR := Request{method, params, time.Now().UnixNano(), "1.0"}
	payloadBuffer := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(payloadBuffer)
	err = jsonEncoder.Encode(rpcR)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", c.URL, payloadBuffer)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	// Auth ?
	if len(c.Username) > 0 || len(c.Password) > 0 {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.DoTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("HTTP error: " + resp.Status)
		return
	}
	err = json.Unmarshal(data, &rr)
	return
}
