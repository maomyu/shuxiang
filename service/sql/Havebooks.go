/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-08-31 10:17:30
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
)

//查询图书信息
func (b *BookHub) Havebooks(r *http.Request) []byte {

	//连接数据库
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	data := Data{}

	isbn := r.FormValue("ISBN")

	BookIn := Findallbooky(Db, isbn)
	// data.Bookcode = r.FormValue("bookcode")
	data.Bookname = BookIn.Bookname
	data.Bookpic = BookIn.Bookpic
	data.Bookauthor = BookIn.Bookauthor
	data.Bookpublic = BookIn.Bookpublic
	data.Bookpages = BookIn.Bookpages
	data.Bookprice = BookIn.Bookprice
	data.Bookisbn = isbn
	data.Bookintroduction = BookIn.Bookintroduction

	//获得此书的类别名称
	// _, cgname := Getcategory(Db, r.FormValue("bookcode"))
	cgname := GetCategoryName(Db, isbn)
	data.Bookcategory = cgname
	//获取此书可借阅的数量
	_, err, ns := Borrownum(Db, data.Bookisbn)
	if err != nil {
		log.Error(err)
	}

	// 获取所有可以借阅的书籍位置id
	cation := Getlocation(Db, data.Bookisbn)
	for _, v := range cation {
		name, _ := strconv.Atoi(v.Bookcode)
		l := Location{
			Bookcode:     v.Locationname,
			Locationname: Getbookshelfname(Db, name),
		}
		// fmt.Println(l)
		data.LendLocation = append(data.LendLocation, l)
	}

	//将此书可借阅的数量赋值给Booknum这个属性
	data.Booknum = ns
	//请求web端获取  d.Lendstatus
	requestBody := fmt.Sprintf(`{"isbn": "%s","userid": "%s"}`, r.FormValue("ISBN"), r.FormValue("userid"))

	urls := "http://127.0.0.1:8200/getuserlendstatusbyid"
	//请求web端
	bodybyte := Requestto(requestBody, urls)
	//返回字节形式，转换为string
	map2 := make(map[string]map[string]int)
	err = json.Unmarshal(bodybyte, &map2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(map2)
	datas := map2["data"]
	fmt.Println(datas)
	stats := datas["status"]
	fmt.Println(stats)

	//var str string = string(bodybyte[:])
	//将string类型转换为int

	data.Lendstatus = stats
	fmt.Println(data)
	cadata := jiagong("200", 1, data)
	byt, err := json.Marshal(cadata)
	fmt.Println(cadata)
	return byt
}
