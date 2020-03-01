package main

import (
	"encoding/json"
	"fmt"
	"shuxiang/common/mq"

	"github.com/streadway/amqp"
)

var Client mq.MessagingClient

// 初始化rabbitmq
func init() {
	Client.Conn = Client.ConnectToRabbitmq("amqp://guest:guest@192.168.10.252:5672")
	// fmt.Println("************")

}

func getBooking(delivery amqp.Delivery) {
	m := make(map[string]string)
	json.Unmarshal(delivery.Body, &m)
	fmt.Println(m)
}
func InsertUser(delivery amqp.Delivery) {
	m := make(map[string]string)
	json.Unmarshal(delivery.Body, &m)
	fmt.Println(m)
	fmt.Println("*********************woshi")

}

type UserLendMQ struct {
	UserID   string `json:"userid"`
	Bookcode string `json:"bookcode"`
	IsBN     string `json:"ISBN"`
	Status   int    `json:"status"`
}

func main() {
	m := &UserLendMQ{
		UserID:   "111",
		IsBN:     "222",
		Bookcode: "222",
		Status:   0,
	}
	body, _ := json.Marshal(m)
	Client.PublishOnQueue(body, "getBooking", "getBooking")
	Client.PublishOnQueue(body, "InsertUser", "InsertUser")
	for {

		Client.ConsumeFromQueue("InsertUser", "InsertUser", InsertUser)
		Client.ConsumeFromQueue("getBooking", "getBooking", getBooking)

	}

}
