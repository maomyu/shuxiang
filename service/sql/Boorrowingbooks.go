/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-24 09:41:10
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
)

func (h *BookHub) Borrowingbooks(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	boolean, _ := ChangeStatus(Db, r.FormValue("bookcode"), "2")
	if boolean == true {
		log.Info("成功")
	}
	// sm := Sendmess{}
	// sm.Bookcode = r.FormValue("bookcode")
	// sm.Userid = r.FormValue("userid")
	// bkisbn := Bookcodeisbn(Db, r.FormValue("bookcode"))
	// bkstatus := Findbookstatuss(Db, r.FormValue("bookcode"))
	// sm.ISBN = bkisbn
	// sm.Status = bkstatus
	// fmt.Println(sm)
	// byt, _ := json.Marshal(sm)

	byt, _ := json.Marshal(struct {
		UserID     string
		Bookcode   string
		Lenddate   time.Time
		Returndate time.Time
		Status     int
		IsBN       string
	}{
		UserID:   r.FormValue("userid"),
		Bookcode: r.FormValue("bookcode"),
		Status:   2,
		IsBN:     Bookcodeisbn(Db, r.FormValue("bookcode")),
	})
	Client.PublishOnQueue(byt, "AddLendDoc", "AddLendDoc")

	cadata := jiagong("200", 1, "成功")
	byts, _ := json.Marshal(cadata)
	return byts
}
