/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-10-23 13:10:15
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/yuwe1/shuxiang/common/reflecter"
	"github.com/yuwe1/shuxiang/common/weber"
)

func Selectfunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// file := "./../../public/userconfig.json"
	file, _ := os.Getwd()
	file = file + "/public/hub.json"
	parsename := weber.GetParsepath(r.RequestURI)
	fmt.Println(parsename)

	route := weber.ParseRoute(file, parsename)
	fmt.Println("route:", route)
	var ruTest *BookHub
	crMap := reflecter.GetClient().FuncNameToFuncHandler(ruTest)
	fmt.Println("路径：", route.Funcname)
	fmt.Println("方法：", crMap[route.Funcname])
	//创建带调用方法时需要传入的参数列表
	parms := []reflect.Value{reflect.ValueOf(r)}
	//使用方法名字符串调用指定方法
	body := crMap[route.Funcname].Call(parms)

	w.Write(body[0].Interface().([]byte))
}

func main() {

	go Client.ConsumeFromQueue("UpdateStorageStatus", "UpdateStorageStatus", UpdateStorageStatus)
	http.HandleFunc("/", Selectfunc)
	http.ListenAndServe("0.0.0.0:8300", nil)
	// http.ListenAndServe("192.169.10.150:8300", nil)

}
