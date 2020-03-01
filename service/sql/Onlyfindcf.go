/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-17 16:00:37
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
	"strconv"
)

//查看图书的分类
func (b *BookHub) Onlyfindcf(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	locgids, err := strconv.Atoi(r.FormValue("categoryid"))
	if err != nil {
		log.Error(err)
	}
	cgid := selectallcgbook(Db, locgids)
	fmt.Println(cgid)
	cadata := jiagong("200", 1, cgid)
	byt, err := json.Marshal(cadata)
	return byt
}
