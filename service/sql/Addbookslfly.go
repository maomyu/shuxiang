package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
	"strconv"
	"time"
)

//添加一个书层
func (h *BookHub) Addbookslfly(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	nowtime := time.Now().String()
	itloid, _ := strconv.Atoi(r.FormValue("locationid"))
	boolean, err, _ := Addashelf(Db, r.FormValue("locationname"), itloid, 0, nowtime, nowtime)
	if err != nil {
		fmt.Println(err)
	}
	if boolean == false {
		log.Info("添加失败")
		cadata := jiagong("500", 0, nil)
		byt, _ := json.Marshal(cadata)
		return byt
	}
	log.Info("添加成功")
	UpdateBookShelf(Db, itloid)
	cadata := jiagong("200", 1, nil)
	byt, _ := json.Marshal(cadata)
	return byt
}
