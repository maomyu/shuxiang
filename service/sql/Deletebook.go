/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-17 15:52:08
 * @LastEditors: your name
 */
package main

import (
	"encoding/json"
	"net/http"
	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
)

func (h *BookHub) Deletebook(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	//根据bookcode查询出isbn
	bkisbn := Bookcodeisbn(Db, r.FormValue("bookcode"))
	//删除此图书
	boolean := Deletethisbook(Db, r.FormValue("bookcode"))

	if boolean == true {
		err, num := Findabooksnum(Db, bkisbn)
		if err != nil {
			log.Error("查询booknum失败")
			cadata := jiagong("500", 0, nil)
			byt, _ := json.Marshal(cadata)
			return byt
		}
		if num == 1 {
			err := Deleteallinfobk(Db, bkisbn)
			if err != nil {
				log.Error("删除此书籍失败")
			}
		} else {
			AlterBooknum(Db, bkisbn, -1)
			cadata := jiagong("200", 1, nil)
			byt, _ := json.Marshal(cadata)
			return byt
		}

	} else {
		cadata := jiagong("500", 0, nil)
		byt, _ := json.Marshal(cadata)
		return byt
	}
	cadata := jiagong("200", 1, nil)
	byt, _ := json.Marshal(cadata)
	return byt
}
