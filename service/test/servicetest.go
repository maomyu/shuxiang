package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"shuxiang/common/mq"
	"shuxiang/common/reflecter"
	"shuxiang/common/weber"
)

var Client mq.MessagingClient

// 初始化rabbitmq
func init() {
	Client.Conn = Client.ConnectToRabbitmq("amqp://guest:guest@192.168.10.252:5672")
}

//定义路由器结构类型
type Routers struct {
}
type User struct {
	UserID string `json:"userID"`
}

func Selectfunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	file := "./../../public/userconfig.json"
	parsename := weber.GetParsepath(r.RequestURI)
	fmt.Println(parsename)

	route := weber.ParseRoute(file, parsename)

	var ruTest *Routers
	crMap := reflecter.GetClient().FuncNameToFuncHandler(ruTest)
	//创建带调用方法时需要传入的参数列表
	parms := []reflect.Value{reflect.ValueOf(r)}
	//使用方法名字符串调用指定方法
	body := crMap[route.Funcname].Call(parms)
	w.Write(body[0].Interface().([]byte))
}

// 对参数进行处理

func (this *Routers) Login(r *http.Request) []byte {
	var user User
	user.UserID = r.FormValue("userID")
	body, _ := json.Marshal(&user)
	// Client.PublishOnQueue(body, "SendMessage", "message")
	Client.PublishOnQueue(body, "InsertUser", "InsertUser")
	Client.PublishOnQueue(body, "getBooking", "getBooking")
	return body
}
func (this *Routers) SaveUser() {

}
func main() {
	http.HandleFunc("/", Selectfunc)
	http.ListenAndServe("0.0.0.0:8200", nil)
}
