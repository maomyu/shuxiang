/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-10-23 14:32:01
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Route struct {
	Server string `json:"server"`
	Addr   string `json:"addr"`
}

func prase(file string, server string) Route {
	jsonFile, _ := os.Open(file)
	defer jsonFile.Close()
	jsonData, _ := ioutil.ReadAll(jsonFile)
	var r []Route
	json.Unmarshal(jsonData, &r)
	for _, v := range r {
		if v.Server == server {
			return v
		}
	}
	return Route{}
}

// 是json时的转换
func GetPath(ustr string) (server []string) {
	path1 := strings.Split(ustr, "/")
	path2 := strings.Split(path1[2], "?")
	server = append(server, path1[1])
	server = append(server, path2[0])
	return server
}
func SelectServer(w http.ResponseWriter, r *http.Request) {
	// 跨域解决
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	// 获得参数对其进行认证,判断token是否还有效

	u_r_l := r.RequestURI
	service := GetPath(u_r_l)
	fmt.Println(service)
	file, _ := os.Getwd()
	file = file + "/public/webconfig.json"
	fmt.Println(file)
	route := prase(file, service[0])
	fmt.Println(route)
	// 获得json数据
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	params := url.Values{}
	var mp map[string]string

	json.Unmarshal(data, &mp)
	fmt.Println(mp)
	for i, v := range mp {
		params.Set(i, v)
	}
	Url, _ := url.Parse(route.Addr + service[1])
	fmt.Println(route.Addr + service[1])

	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println("请求地址：", urlPath)
	resp, _ := http.Get(urlPath)

	defer resp.Body.Close()

	robots, _ := ioutil.ReadAll(resp.Body)
	w.Write(robots)
}
func main() {
	http.HandleFunc("/", SelectServer)
	// http.ListenAndServe(":8100", nil)
	http.ListenAndServe("0.0.0.0:8100", nil)
	// http.ListenAndServe("192.169.10.150:8100", nil)
}
