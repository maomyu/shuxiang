/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-26 15:13:13
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
)

//查询某个图书的所有详情
func (b *BookHub) Haveallbooks(r *http.Request) []byte {

	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	//查询出此bookcode的isbn
	isbn := Bookcodeisbn(Db, r.FormValue("bookcode"))
	fmt.Println("**************" + isbn)
	if isbn != "" {
		//根据isbn查询出此书的全部信息
		BookIn := Findallbooky(Db, isbn)
		data1 := Data1{}
		data1.Bookcode = r.FormValue("bookcode")
		data1.Bookname = BookIn.Bookname
		data1.Bookpic = BookIn.Bookpic
		data1.Bookauthor = BookIn.Bookauthor
		data1.Bookpublic = BookIn.Bookpublic
		data1.Bookpages = BookIn.Bookpages
		data1.Bookprice = BookIn.Bookprice
		data1.Bookisbn = BookIn.Bookisbn
		data1.Booknum = BookIn.Booknum
		data1.Categoryid = BookIn.Categoryid
		data1.Categoryid, data1.Bookcategory = Getcategory(Db, data1.Bookcode)
		status := Findbookstatuss(Db, data1.Bookcode)
		data1.Status = status
		Locationid, err := strconv.Atoi(Getbookshelfloid(Db, r.FormValue("bookcode")))
		data1.Locationid = Locationid
		if err != nil {
			data1.Locationidname = "查无id"
		}
		data1.Locationidname = Getbookshelfname(Db, Locationid)
		if data1.Locationidname == "" {
			data1.Locationidname = "查无此名"
		}

		cadata := jiagong("200", 1, data1)
		byt, err := json.Marshal(cadata)
		return byt

	} else {
		log.Info("查无此书")

	}
	return nil
}
