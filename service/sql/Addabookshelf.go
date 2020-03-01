/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-17 16:04:45
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yuwe1/shuxiang/common/dber"
)

func (b *BookHub) Addabookshelf(r *http.Request) []byte {
	// fmt.Println(r.FormValue("locationname"))
	loca := r.FormValue("locationname")
	fmt.Println(loca)
	if len(loca) <= 0 {
		cadata := jiagong("200", 0, "名字为空")
		byt, _ := json.Marshal(cadata)
		return byt
	} else {
		c := dber.GetClient()
		Db := c.ConnectTry(username, password, url, dbname)
		defer Db.Close()
		nowtime := time.Now().String()
		boolean, err, userid := Addashelf(Db, loca, -1, 0, nowtime, nowtime)

		if boolean == false {
			fmt.Println("添加失败")
			fmt.Println(err)
			cadata := jiagong("500", 0, nil)
			byt, _ := json.Marshal(cadata)
			return byt
		}
		cadata := jiagong("200", userid, nil)
		byt, _ := json.Marshal(cadata)
		return byt
	}

}
