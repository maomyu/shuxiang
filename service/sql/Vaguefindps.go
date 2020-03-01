/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-17 16:00:52
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/yuwe1/shuxiang/common/dber"
)

func (h *BookHub) Vaguefindps(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	vaguepsss := Likefind(Db, r.FormValue("requestContent"), r.FormValue("category"))
	fmt.Println(vaguepsss)

	cadata := jiagong("200", 1, vaguepsss)

	fmt.Println(cadata)
	byt, _ := json.Marshal(cadata)
	fmt.Println(byt)
	return byt
}
