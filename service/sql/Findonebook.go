/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-17 15:59:32
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

func (b *BookHub) FindOnebook(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	finone := Findallonebookss(Db)
	log.Info(finone)
	fmt.Println(finone)
	cadata := jiagong("200", 1, finone)
	byt, _ := json.Marshal(cadata)
	return byt
}
