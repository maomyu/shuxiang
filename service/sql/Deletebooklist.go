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
	"fmt"
	"net/http"
	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
	"strconv"
)

//删除一级或二级书架
func (h *BookHub) Deletebooklist(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()

	itloid, _ := strconv.Atoi(r.FormValue("locationid"))

	//查询出parentid
	isid := isornoparent(Db, itloid)
	//说明是父书架//删除此书架上的所有书籍

	if isid == -1 {

		//查询某个书架下的子书架
		bsslice := findallbsbook(Db, itloid)

		//遍历所有的子书架
		for _, issl := range bsslice {
			itloidss, _ := strconv.Atoi(issl)
			fmt.Println(itloidss)
			//查询此书架下的书籍
			slice := findallbssbook(Db, itloidss)
			for _, sl := range slice {
				boolean, _ := ChangeStatus(Db, sl, "4")
				if boolean == true {
					log.Info("删除成功")
				}
			}
			boolean := Deletethisbs(Db, itloidss)
			if boolean == false {
				log.Info("删除失败")
			}
		}
		slice := findallbssbook(Db, itloid)
		for _, sl := range slice {
			boolean, _ := ChangeStatus(Db, sl, "4")
			if boolean == true {
				log.Info("删除成功")
			}
		}
		boolean := Deletethisbs(Db, itloid)
		if boolean == false {
			log.Info("删除失败")
		}

	} else {
		slice := findallbssbook(Db, itloid)
		for _, sl := range slice {
			boolean, _ := ChangeStatus(Db, sl, "4")
			if boolean == true {
				fmt.Println("删除成功")
			}
		}
		boolean := Deletethisbs(Db, itloid)
		if boolean == false {
			log.Info("删除失败")
		}
	}

	cadata := jiagong("200", 1, nil)
	byt, _ := json.Marshal(cadata)
	return byt

}
