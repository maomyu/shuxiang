/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-31 09:53:41
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yuwe1/shuxiang/common/dber"
)

//查看在馆图书的书籍
func (h *BookHub) Geteverybook(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	defer Db.Close()
	//声明一个结构体切片
	var Bookallstatuss []Bookallstatus
	allcode := Geteverybookcode(Db, r.FormValue("status"))
	fmt.Println("bookcode:", allcode)
	for _, code := range allcode {
		//根据bookcode查询locationid
		bkloid := Getbookshelfloid(Db, code)
		//根据locationid获取locationname
		itbkid, _ := strconv.Atoi(bkloid)
		bkname := Getbookshelfname(Db, itbkid)
		//根据bookcode查询isbn
		bkisbn := Bookcodeisbn(Db, code)
		//查询图书的所有信息
		bkin := Findallbooky(Db, bkisbn)
		//获取图书的状态
		itstatus, _ := strconv.Atoi(r.FormValue("status"))
		var userid string
		if itstatus == 2 {

			requestBody := fmt.Sprintf(`{"bookcode": "%s","status": "%s"}`, code, r.FormValue("status"))
			urls := "http://127.0.0.1:8100/user/getuserlendstatus"
			//请求web端
			bodybyte := Requestto(requestBody, urls)
			fmt.Println(bodybyte)
			map2 := make(map[string]map[string]string)
			err := json.Unmarshal(bodybyte, &map2)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(map2)
			userids := map2["data"]
			fmt.Println(userids)
			userid = userids["userID"]
			fmt.Println(userid)
		} else {
			userid = ""
		}
		bas := Bookallstatus{}
		bas.Bookcode = code
		bas.Bookname = bkin.Bookname
		bas.Bookauthor = bkin.Bookauthor
		bas.Bookpic = bkin.Bookpic
		bas.Bookpublic = bkin.Bookpublic
		bas.Locationname = bkname
		bas.Status = itstatus
		bas.UserID = userid
		Bookallstatuss = append(Bookallstatuss, bas)
	}
	//fmt.Println(Bookallstatuss)
	cadata := jiagong("200", 1, Bookallstatuss)
	byt, _ := json.Marshal(cadata)
	return byt
}
