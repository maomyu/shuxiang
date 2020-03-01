/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-22 11:07:25
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/yuwe1/shuxiang/common/dber"
)

//增加图书
func (b *BookHub) Addbooks(r *http.Request) []byte {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()

	//获取随机生成的id
	Bookcode := Rangecode()
	i, err := strconv.Atoi(r.FormValue("categoryid"))
	strmess := &Storagemess{
		Bookcode:     Bookcode,
		BookISBN:     r.FormValue("bookisbn"),
		Locationid:   -1,
		Cglocationid: i,
		Status:       5,
	}

	//生成书籍信息表中添加数据
	ita, err := strconv.Atoi(r.FormValue("bookpages"))
	itb, err := strconv.Atoi(r.FormValue("categoryid"))
	itc, err := strconv.Atoi(r.FormValue("booknum"))
	bl := &BookInformation{
		Bookname:         r.FormValue("bookname"),
		Bookauthor:       r.FormValue("bookauthor"),
		Bookpublic:       r.FormValue("bookpublic"),
		Bookpages:        ita,
		Bookpic:          r.FormValue("bookpic"),
		Bookprice:        r.FormValue("bookprice"),
		Bookintroduction: r.FormValue("bookintroduction"),
		Bookisbn:         r.FormValue("bookisbn"),
		Categoryid:       itb,
		Booknum:          itc,
	}

	mess := Havethisbook(Db, r.FormValue("bookisbn"))
	//如果信息表中有此信息让其图书总数量+1，并且向存储关系表添加一条数据

	if mess != "" {
		//将Cglocationid的string类型转化为int类型
		//将书籍信息表中的总数里昂+1

		AlterBooknum(Db, mess, itc)
		for i := 0; i < itc; i++ {
			//向存储关系表中添加一条数据
			AddRtsbooks(Db, strmess)
		}

	} else {
		err = Addbkufmation(Db, bl)
		if err != nil {
			log.Println("错误")
		}
		for i := 0; i < itc; i++ {
			//向存储关系表中添加一条数据
			AddRtsbooks(Db, strmess)
		}
	}
	d := Data{}
	cadata := jiagong("200", 1, d)
	byt, err := json.Marshal(cadata)
	return byt
}
