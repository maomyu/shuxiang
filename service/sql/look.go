/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-31 10:29:37
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
)

//查看在馆图书的书籍
func (h *BookHub) Look(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	//声明一个结构体切片

	allcode := Geteverybookcode(Db, "0")

	totalnums := len(allcode)
	tn := Lendnums{}
	tn.Hubnum = totalnums
	_, _, tm := Totalnumbook(Db)
	tn.Total = tm

	urls := "http://127.0.0.1:8100/user/getreadnum"
	//请求web端
	bodybyte := Requestto("", urls)
	log.Info(bodybyte)
	map2 := make(map[string]map[string]int)

	err := json.Unmarshal(bodybyte, &map2)
	if err != nil {
		log.Info(err)
	}
	log.Info(map2)
	userids := map2["data"]
	log.Info(userids)
	fmt.Println(map2)
	lendnum, _ := userids["lendnum"]
	ttlnum, _ := userids["ttlnum"]
	yearreadnum, _ := userids["yearreadnum"]

	tn.Lendnum = lendnum
	tn.Ttlnum = ttlnum
	tn.Yearreadnum = yearreadnum
	fmt.Println(tn)
	cadata := jiagong("200", 1, tn)
	byt, _ := json.Marshal(cadata)
	return byt
}
