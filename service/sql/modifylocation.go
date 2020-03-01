/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-24 09:40:48
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yuwe1/shuxiang/common/dber"
)

//修改图书信息
func (h *BookHub) ModifyLocation(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	itpages, err := strconv.Atoi(r.FormValue("bookpages"))
	itcagid, err := strconv.Atoi(r.FormValue("categoryid"))
	itbk, err := strconv.Atoi(r.FormValue("booknum"))
	bookin := BookInformation{

		Bookname:         r.FormValue("bookname"),
		Bookauthor:       r.FormValue("bookauthor"),
		Bookpublic:       r.FormValue("bookpublic"),
		Bookpages:        itpages,
		Bookprice:        r.FormValue("bookprice"),
		Bookisbn:         r.FormValue("bookisbn"),
		Bookintroduction: r.FormValue("bookname"),
		Categoryid:       itcagid,
		Booknum:          itbk,
	}
	_, err = ChangeStatus(Db, r.FormValue("bookcode"), r.FormValue("status"))
	boolean := alterBookmation(Db, bookin)
	if boolean == true && err == nil {
		cadata := jiagong("200", 1, nil)
		byt, _ := json.Marshal(cadata)
		return byt
	}
	return nil
}
