/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-10-23 11:39:38
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/yuwe1/shuxiang/common/dingding"
	"github.com/yuwe1/shuxiang/common/mq"

	"github.com/streadway/amqp"
)

var Client mq.MessagingClient

type Message struct {
	Msgtype string      `json:"msgtype"`
	Text    TextMessage `json:"text"`
}
type TextMessage struct {
	Content string `json:"content"`
}
type SendMessageReuslt struct {
	UserID  string `json:"userID"`
	Content string `json:"content"`
}
type SendResult struct {
	UserID string `json:"userID"`
	Url    string `json:"url"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

// 初始化rabbitmq
func init() {
	Client.Conn = Client.ConnectToRabbitmq("amqp://admin:admin@192.168.10.200:20026")

}

// 实现发送消息
func SendMessage(d amqp.Delivery) {
	fmt.Println("*************************")
	// m := make(map[string]string)
	var r SendResult
	json.Unmarshal(d.Body, &r)
	params := url.Values{}
	agentID := "265907067"
	client := &http.Client{}
	accesstoken := dingding.GetAccessToken()
	fmt.Println("accesstoken:" + accesstoken)
	urlPath := "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=" + accesstoken
	message := &dingding.LinkMessageResult{
		Msgtype: "link",
		Link: dingding.LinkResult{
			MessageUrl: r.Url,
			PicUrl:     "test",
			Title:      r.Title,
			Text:       r.Text,
		},
	}
	body, _ := json.Marshal(message)
	params.Set("agent_id", agentID)
	params.Set("userid_list", r.UserID)
	params.Set("msg", string(body))
	fmt.Println(string(body))
	data := ioutil.NopCloser(strings.NewReader(params.Encode())) //把form数据编下码
	req, _ := http.NewRequest("POST", urlPath, data)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	robots, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(robots))
	d.Acknowledger.Ack(d.DeliveryTag, true)
	fmt.Println(r)
}
func SendTextMessage(d amqp.Delivery) {
	fmt.Println("*************************")
	// m := make(map[string]string)
	var r SendMessageReuslt
	json.Unmarshal(d.Body, &r)
	params := url.Values{}
	agentID := "265907067"
	client := &http.Client{}
	accesstoken := dingding.GetAccessToken()
	fmt.Println("accesstoken:" + accesstoken)
	urlPath := "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=" + accesstoken
	message := &Message{
		Msgtype: "text",
		Text: TextMessage{
			Content: r.Content,
		},
	}
	body, _ := json.Marshal(message)
	params.Set("agent_id", agentID)
	params.Set("userid_list", r.UserID)
	params.Set("msg", string(body))
	data := ioutil.NopCloser(strings.NewReader(params.Encode())) //把form数据编下码
	req, _ := http.NewRequest("POST", urlPath, data)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	robots, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(robots))
	d.Acknowledger.Ack(d.DeliveryTag, true)
	fmt.Println(r)
}
func main() {

	go Client.ConsumeFromQueue("SendMessage", "SendMessage", SendMessage)

	go Client.ConsumeFromQueue("SendTextMessage", "SendTextMessage", SendTextMessage)

	http.ListenAndServe("0.0.0.0:9000", nil)
}
