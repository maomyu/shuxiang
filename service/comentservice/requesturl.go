package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func RequestURL(data interface{}, urls string) []byte {

	parms := url.Values{}

	parmbody, _ := json.Marshal(data)
	parms.Add("data", string(parmbody))
	paemsbody := ioutil.NopCloser(strings.NewReader(parms.Encode()))

	req, err := http.NewRequest("POST", urls, paemsbody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	CheckError("获取response失败！", err)

	body, err := ioutil.ReadAll(resp.Body)
	CheckError("获取body失败！", err)

	return body
}
