/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-22 09:16:07
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

func (h *BookHub) Onlyfindbook(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	var bss []Bookshelfs
	loids, err := strconv.Atoi(r.FormValue("locationid"))
	if err != nil {
		log.Error(err)
	}
	//查询出所有licationid
	cd := findallbsbook(Db, loids)
	for _, cid := range cd {

		bs := Bookshelfs{}
		loids, _ := strconv.Atoi(cid)
		cname := Getbookshelfname(Db, loids)
		bs.Locationid = loids
		bs.Locationname = cname
		bss = append(bss, bs)
	}
	cadata := jiagong("200", 1, bss)
	byt, err := json.Marshal(cadata)
	return byt
}
