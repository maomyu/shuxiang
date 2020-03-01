/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-28 10:53:55
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	urll "net/url"
	"time"

	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"

	"github.com/streadway/amqp"
)

type SendResult struct {
	UserID string `json:"userID"`
	Url    string `json:"url"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}
type WhetherConmentParm struct {
	Userid   string `json:"userid"`
	Bookisbn string `json:"bookisbn"`
	Token    string `json:"token"`
}

func (h *BookHub) Returnabooks(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	fmt.Println(r.FormValue("bookcode"))
	boolean, err := ChangeStatus(Db, r.FormValue("bookcode"), "4")

	if err != nil {
		log.Error(err)
	}
	if boolean == true {
		log.Info("成功")
	}
	sm := Sendmess{}
	sm.Bookcode = r.FormValue("bookcode")
	sm.Userid = r.FormValue("userid")
	bkisbn := Bookcodeisbn(Db, r.FormValue("bookcode"))
	//bkstatus := Findbookstatuss(Db, r.FormValue("bookcode"))
	sm.ISBN = bkisbn
	sm.Status = 2
	byt, err := json.Marshal(sm)
	if err != nil {
		log.Error(err)
	}
	Client.PublishOnQueue(byt, "UpdateLendDoc", "UpdateLendDoc")
	//
	t := time.Now()
	formatNow := t.Format("2006-01-02 15:04:05")
	p := &WhetherConmentParm{
		Userid:   sm.Userid,
		Bookisbn: bkisbn,
		Token:    dber.GetRandomID(),
	}
	bodyp, _ := json.Marshal(p)

	// http://192.168.10.150:8100/comment/whether
	Url, _ := urll.Parse("eapp://page/my/mycomment/index")
	urlPath := Url.String() + "?data=" + string(bodyp)
	sendResult := &SendResult{
		UserID: sm.Userid,
		Url:    urlPath,
		Title:  "归还通知",
		Text:   formatNow + "成功归还\n点此链接进行评论",
	}
	AddToken(p.Token)
	fmt.Println(sendResult)
	var delivery amqp.Delivery
	delivery.Body, _ = json.Marshal(sendResult)
	Client.PublishOnQueue(delivery.Body, "SendMessage", "SendMessage")

	cadata := jiagong("200", 1, nil)
	byts, _ := json.Marshal(cadata)
	return byts
}
