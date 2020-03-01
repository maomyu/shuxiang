/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-22 12:54:36
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
)

func (b *BookHub) FindCgbooks(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	itcgid, _ := strconv.Atoi(r.FormValue("categoryid"))
	cgslice := findallcgbook(Db, itcgid)

	var bookxins []Bookxinxi
	for _, book := range cgslice {
		bookxin := Bookxinxi{}
		bookin := Findallbooky(Db, book)
		bookxin.Bookname = bookin.Bookname
		bookxin.Bookpic = bookin.Bookpic
		bookxin.Bookauthor = bookin.Bookauthor
		bookxin.Bookpublic = bookin.Bookpublic
		bookxin.Bookpages = bookin.Bookpages
		bookxin.Bookprice = bookin.Bookprice
		bookxin.Bookisbn = bookin.Bookisbn
		_, err, num := Borrownum(Db, book)
		if err != nil {
			log.Error("查询剩余数量失败")
		}
		bookxin.Booknum = num
		bookxins = append(bookxins, bookxin)
	}

	cadata := jiagong("200", 1, bookxins)
	byt, _ := json.Marshal(cadata)
	return byt
}
