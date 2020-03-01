/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-10-23 11:22:29
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	urll "net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"github.com/yuwe1/shuxiang/common/dber"
	"github.com/yuwe1/shuxiang/common/log"
	"github.com/yuwe1/shuxiang/common/mq"

	"github.com/garyburd/redigo/redis"

	_ "github.com/go-sql-driver/mysql"
)

type Statusdata struct {
	Status  string
	Success int
	Data    interface{}
}
type Lendnums struct {
	Lendnum     int `json:"lendnum"`
	Ttlnum      int `json:"ttlnum"`
	Yearreadnum int `json:"yearreadnum"`
	Hubnum      int `json:"hubnum"`
	Total       int `json:"total"`
}

type Storagemess struct {
	Bookcode     string //书的唯一标识
	BookISBN     string //所属的书籍信息ISBN
	Locationid   int    //所属位置
	Cglocationid int    //分类信息的id
	Status       int    //图书的状态
}

type Data struct {
	Bookcode         string
	Bookname         string
	Bookpic          string
	Bookauthor       string
	Bookpublic       string
	Bookpages        int
	Bookprice        string
	Bookisbn         string
	Bookintroduction string
	Booknum          int    //当前书籍可借阅的数量
	Bookcategory     string //当前书籍的类别名字
	Lendstatus       int    //此用户对此书的借阅状态
	LendLocation     []Location
}
type Location struct {
	Bookcode     string
	Locationname string
}
type Data1 struct {
	Bookcode         string
	Bookname         string
	Bookpic          string
	Bookauthor       string
	Bookpublic       string
	Bookpages        int
	Bookprice        string
	Bookisbn         string
	Bookintroduction string
	Booknum          int    //当前书籍可借阅的数量
	Categoryid       int    //当前书籍的分类
	Bookcategory     string //当前书籍的类别名字
	Status           int    //当前书籍的状态
	Locationid       int    //当前书架的位置id
	Locationidname   string //当前书架的名字
}

type BookInformation struct {
	Bookname         string
	Bookauthor       string
	Bookpublic       string
	Bookpages        int
	Bookprice        string
	Bookisbn         string
	Bookpic          string
	Bookintroduction string
	Categoryid       int
	Booknum          int
}

type Bookxinxi struct {
	Bookname   string
	Bookpic    string
	Bookauthor string
	Bookpublic string
	Bookpages  int
	Bookprice  string
	Bookisbn   string
	Booknum    int
}
type Bookallstatus struct {
	Bookcode     string
	Bookname     string
	Bookauthor   string
	Bookpic      string
	Bookpublic   string
	Locationname string
	Status       int
	UserID       string
}
type Bookshelfs struct {
	Locationid   int
	Locationname string
}
type Classification struct {
	Categoryid   int
	Bookcategory string
	Parentid     int
	Isparent     int
}
type Sendmess struct {
	Userid   string `json:"userid"`
	Bookcode string `json:"bookcode"`
	ISBN     string `json:"ISBN"`
	Status   int    `json:"status"`
}
type Vaguepst struct {
	Bookisbn   string
	Bookname   string
	Bookauthor string
	Bookpic    string
	Bookpublic string
}

var (
	username = "root"
	password = "03354ab3"
	url      = "192.168.10.200:20024"
	dbname   = "hub"
)

type BookHub struct{}

var Client mq.MessagingClient

//初始化rabbitmq
func init() {
	Client.Conn = Client.ConnectToRabbitmq("amqp://admin:admin@192.168.10.200:20026")
}

//将数据封装至结构体返回
func jiagong(Status string, Success int, d interface{}) interface{} {
	c := Statusdata{
		Status:  Status,
		Success: Success,
		Data:    d,
	}
	return c

}

//随机生成bookcode
func Rangecode() string {
	str := "ycf"
	times := time.Now()
	t := times.String()
	data := str + t
	bytes := []byte(data)
	c := md5.Sum(bytes)
	H := fmt.Sprintf("%x", c)
	return H
}

//向web端发送消息
func Requestto(requestBody string, urls string) []byte {
	client := &http.Client{}

	urlmap := urll.Values{}

	urlmap.Add("data", requestBody)
	fmt.Println("数据：", urlmap)
	parms := ioutil.NopCloser(strings.NewReader(urlmap.Encode())) //把form数据编下码
	req, err := http.NewRequest("POST", urls, parms)
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	return body
}

//修改图书数量（内部）
func AlterBooknum(Db *sql.DB, isbn string, numdata int) (b bool, err error) {

	_, num := Findabooksnum(Db, isbn)
	fmt.Println(num)
	fmt.Println(num + numdata)
	stmt, err := Db.Prepare("UPDATE  BookInformation SET Booknum=? WHERE Bookisbn=?")
	if err != nil {
		log.Error("Prepare出错啦")
		return
	}

	_, err = stmt.Exec((num + numdata), isbn)
	if err != nil {
		log.Error("fetch last update id failed:", err.Error())
		return
	}
	return
}

//删除书籍信息表中的此类书籍
func Deleteallinfobk(Db *sql.DB, isbn string) error {
	stmt, err := Db.Prepare("DELETE FROM BookInformation WHERE Bookisbn=?")
	if err != nil {
		log.Error("删除失败", err)
	}
	_, err = stmt.Exec(isbn)
	if err != nil {
		log.Error("删除失败", err)
	}
	return nil
}

//查询书籍信息表中书籍的数量(内部)
func Findabooksnum(Db *sql.DB, isbn string) (errs error, num int) {
	var nums int
	err := Db.QueryRow("SELECT Booknum FROM BookInformation WHERE Bookisbn=?", isbn).Scan(&nums)
	if err != nil {
		log.Error("Prepare出错啦")
		return err, -1
	}

	return nil, nums
}

//向存储关系表添加图书(内部)
func AddRtsbooks(Db *sql.DB, s *Storagemess) (b bool, err error) {
	s.Bookcode = Rangecode()
	stmt, err := Db.Prepare("INSERT INTO Storagemess(Bookcode,Bookisbn,Locationid,Cglocationid,Status) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Error("failed", err.Error())
		return
	}
	_, err = stmt.Exec(s.Bookcode, s.BookISBN, s.Locationid, s.Cglocationid, s.Status)
	if err != nil {
		log.Error("插入失败")
		return
	}
	return
}

//向书籍信息表添加数据(内部)
func Addbkufmation(Db *sql.DB, s *BookInformation) error {

	stmt, err := Db.Prepare("INSERT INTO BookInformation(Bookname,Bookauthor,Bookpublic,Bookpages,Bookprice,Bookisbn,Bookpic,Bookintroduction,Categoryid,Booknum) VALUES(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Error("failed", err.Error())
		return err
	}
	_, err = stmt.Exec(s.Bookname, s.Bookauthor, s.Bookpublic, s.Bookpages, s.Bookprice, s.Bookisbn, s.Bookpic, s.Bookintroduction, s.Categoryid, s.Booknum)
	if err != nil {
		log.Error("插入失败")
		return err
	}
	return nil
}

//获取当前书籍可以借阅的数量（内部）
func Borrownum(Db *sql.DB, isbn string) (b bool, err error, num int) {

	rows, err := Db.Query("SELECT * FROM Storagemess WHERE Bookisbn = ?&&Status=0", isbn)
	if err != nil {
		log.Error(err)
	}
	var i int
	i = 0
	for rows.Next() {
		i++
	}
	return true, nil, i
}

//获得所有的未归还数量
func Noreturnbook(Db *sql.DB, status string) (b bool, err error, num int) {

	rows, err := Db.Query("SELECT * FROM Storagemess WHERE Status=1")
	if err != nil {
		log.Error(err)
	}

	var i int
	i = 0
	for rows.Next() {
		i++
	}
	return true, nil, i
}

//根据bookcode查询书的isbn
func Bookcodeisbn(Db *sql.DB, bookcode string) string {
	var isbn string
	err := Db.QueryRow("SELECT Bookisbn FROM Storagemess WHERE Bookcode=?", bookcode).Scan(&isbn)
	if err != nil {
		log.Error("查询失败")
		return ""
	}
	fmt.Println("isbn=", isbn)
	return isbn
}

//根据ISBN查询所有书的cookcode
func Isbnbookcode(Db *sql.DB, ISBN string) []string {
	var isbn []string
	rows, err := Db.Query("SELECT Bookisbn FROM Storagemess WHERE Bookcode=?", ISBN)
	if err != nil {
		log.Error("查询失败")
		// return ""
	}
	for rows.Next() {
		var Bookisbns string
		rows.Scan(&Bookisbns)
		isbn = append(isbn, Bookisbns)
	}
	return isbn
}

//查询location的parentid是否为0

//总共的库存量
func Totalnumbook(Db *sql.DB) (b bool, err error, num int) {

	rows, err := Db.Query("SELECT Booknum FROM BookInformation")
	if err != nil {
		log.Error(err)
	}
	var totalnum int = 0
	var nums int
	for rows.Next() {
		rows.Scan(&nums)
		totalnum += nums
	}
	return true, nil, totalnum
}

//添加一级书架
func Addashelf(Db *sql.DB, Locationname string, Parentid int, Isparent int, Created string, Updated string) (b bool, err error, i int) {
	stmt, err := Db.Prepare("INSERT INTO Bookshelf(Locationname,Parentid,Isparent,Created,Updated) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Error("failed", err.Error())
		return false, err, 0
	}
	r, err := stmt.Exec(Locationname, Parentid, Isparent, Created, Updated)
	if err != nil {
		log.Error("添加失败")
		fmt.Println(err)
		return false, err, 1
	}
	id, err := r.LastInsertId()
	return true, nil, int(id)
}

//查询存储关系表中总共有多少本书
func Totalisbook(Db *sql.DB) (b bool, err error, num int) {

	rows, err := Db.Query("SELECT Bookcode FROM Storagemess ")
	if err != nil {
		log.Error(err)
	}
	var i int
	i = 0
	for rows.Next() {
		i++
	}
	fmt.Println(i)
	return true, nil, i
}

//去除切片中重复的元素
func SliceRemoveDuplicates(s []string) []string {

	sort.Strings(s)
	i := 0
	var j int
	for {
		if i >= len(s)-1 {
			break
		}
		for j = i + 1; j < len(s) && s[i] == s[j]; j++ {

		}
		s = append(s[:i+1], s[j:]...)
		i++
	}
	return s
}

//查询某个分类下的所有书,返回书的isbn
func findallcgbook(Db *sql.DB, categoryid int) []string {

	rows, err := Db.Query("SELECT Bookisbn FROM Storagemess WHERE Cglocationid=?", categoryid)
	if err != nil {
		log.Error("查询失败")
	}
	var s []string
	for rows.Next() {
		var Bookisbns string
		rows.Scan(&Bookisbns)
		s = append(s, Bookisbns)
	}
	sl := SliceRemoveDuplicates(s)
	return sl
}

//查询某个书架下的所有书的code，返回书的bookcode
func findallbssbook(Db *sql.DB, locationid int) []string {
	rows, err := Db.Query("SELECT Bookcode FROM Storagemess WHERE Locationid=?", locationid)
	if err != nil {
		log.Error("查询失败")
	}
	var s []string
	for rows.Next() {
		var Bookcd string
		rows.Scan(&Bookcd)
		s = append(s, Bookcd)
	}
	return s
}

//查询某个书架下的所有子书架,返回书架的切片
func findallbsbook(Db *sql.DB, Locationid int) []string {
	rows, err := Db.Query("SELECT Locationid FROM Bookshelf WHERE Parentid=?", Locationid)
	if err != nil {
		log.Error("查询失败")
	}
	var s []string
	for rows.Next() {
		var Bookid string
		rows.Scan(&Bookid)
		s = append(s, Bookid)
	}
	return s
}

//查询某个分类下的子分类所有信息
func selectallcgbook(Db *sql.DB, categoryid int) []Classification {
	rows, err := Db.Query("SELECT Categoryid,Bookcategory,Parentid,Isparent FROM Classification WHERE Parentid=?", categoryid)
	if err != nil {
		log.Error("查询失败")
	}
	var Classifications []Classification
	for rows.Next() {
		classification := Classification{}
		rows.Scan(&classification.Categoryid, &classification.Bookcategory, &classification.Parentid, &classification.Isparent)
		Classifications = append(Classifications, classification)
	}
	return Classifications
}

//更改图书的状态
func ChangeStatus(Db *sql.DB, bookcode string, status string) (b bool, err error) {

	s, err := strconv.Atoi(status)
	fmt.Println(s)
	stmt, err := Db.Prepare("UPDATE Storagemess SET Status=? WHERE Bookcode=?")
	if err != nil {
		log.Error("Prepare出错啦")
		return false, err
	}
	_, err = stmt.Exec(s, bookcode)
	if err != nil {
		log.Error("fetch last update id failed:", err.Error())
		return false, err
	}
	return true, nil
}

//根据bookcode查询图书的状态
func Findbookstatuss(Db *sql.DB, bookcode string) int {
	var status int
	err := Db.QueryRow("SELECT Status FROM Storagemess WHERE Bookcode=?", bookcode).Scan(&status)
	if err != nil {
		return 0
	}

	fmt.Println("status=", status)
	return status
}

//修改图书信息
func alterBookmation(Db *sql.DB, b BookInformation) bool {
	stmt, err := Db.Prepare("UPDATE BookInformation SET Bookname=?,Bookauthor=?,Bookpublic=?,Bookpages=?,Bookprice=?,Bookisbn=?,Bookintroduction=?,Categoryid=?,Booknum=? WHERE Bookisbn=?")
	if err != nil {
		log.Error("Prepare出错啦")
		return false
	}
	_, err = stmt.Exec(b.Bookname, b.Bookauthor, b.Bookpublic, b.Bookpages, b.Bookprice, b.Bookisbn, b.Bookintroduction, b.Categoryid, b.Booknum, b.Bookisbn)
	if err != nil {
		log.Error("fetch last update id failed:", err.Error())
		return false
	}
	return true
}

//查询某本书所在书架的位置
func Getthislocation(Db *sql.DB, bookcode int) (b bool, err error, s []int) {

	//根据书架的bookcode查询出locationid
	var sv []int
	var loid int
	var ispat int
	err = Db.QueryRow("SELECT Locationid，Isparent FROM Storagemess WHERE Bookcode=?", bookcode).Scan(&loid, &ispat)
	if err != nil {
		log.Error("Prepare出错啦")
		return
	}
	fmt.Println(loid)
	//根据locationid查询出parentid
	if ispat == 0 {
		sv = append(sv, loid)
	} else {
		var parid int
		err = Db.QueryRow("SELECT Parentid FROM Bookshelf WHERE Locationid=?", loid).Scan(&parid)
		if err != nil {
			log.Error("Prepare出错啦")
			return
		}
		if parid != 0 {
			sv = append(sv, parid)
			sv = append(sv, loid)
		} else {
			sv = append(sv, loid)
		}
	}

	return false, nil, sv
}

//查询图书所在书架的名称
func Getbookshelfname(Db *sql.DB, Locationid int) string {
	var Locationname string
	err := Db.QueryRow("SELECT Locationname FROM Bookshelf WHERE Locationid=?", Locationid).Scan(&Locationname)
	if err != nil {
		log.Error("查询失败")
		return ""
	}
	return Locationname
}

//查询所有一级存储库
func Findallonebookss(Db *sql.DB) []Bookshelfs {
	rows, err := Db.Query("SELECT Locationid,Locationname FROM Bookshelf WHERE Parentid=-1")
	if err != nil {
		log.Error("找不到", err)
	}
	var Bookshelfss []Bookshelfs
	for rows.Next() {
		Bookisbns := Bookshelfs{}
		rows.Scan(&Bookisbns.Locationid, &Bookisbns.Locationname)
		Bookshelfss = append(Bookshelfss, Bookisbns)
	}
	return Bookshelfss
}

//删除图书
func Deletethisbook(Db *sql.DB, bookcode string) bool {
	stmt, err := Db.Prepare("DELETE FROM Storagemess WHERE Bookcode=?")
	if err != nil {
		log.Error("DELETE is error", err)
		return false
	}
	_, err = stmt.Exec(bookcode)

	if err != nil {
		log.Error("fetch last delete id failed:", err.Error())
		return false
	}
	return true
}

//删除书架
func Deletethisbs(Db *sql.DB, locationid int) bool {
	stmt, err := Db.Prepare("DELETE FROM Bookshelf WHERE Locationid=?")
	if err != nil {
		log.Error("DELETE is error", err)
		return false
	}
	_, err = stmt.Exec(locationid)

	if err != nil {
		log.Error("fetch last delete id failed:", err.Error())
		return false
	}
	return true
}

//查询图书的Locationid
func Getbookshelfloid(Db *sql.DB, bookcode string) string {
	var Locationid string
	err := Db.QueryRow("SELECT Locationid FROM Storagemess WHERE bookcode=?", bookcode).Scan(&Locationid)
	if err != nil {
		log.Error(err)
		return ""
	}
	return Locationid
}

//查询图书的所有信息
func Findallbooky(Db *sql.DB, isbn string) BookInformation {
	bookIn := BookInformation{}
	err := Db.QueryRow("SELECT * FROM BookInformation WHERE Bookisbn=?", isbn).Scan(&bookIn.Bookname, &bookIn.Bookauthor, &bookIn.Bookpublic, &bookIn.Bookpages, &bookIn.Bookprice, &bookIn.Bookisbn, &bookIn.Bookpic, &bookIn.Bookintroduction, &bookIn.Categoryid, &bookIn.Booknum)
	if err != nil {
		log.Error("找不到", err)

	}
	return bookIn
}

//查询书籍信息表是否存在此书籍(内部)
func Havethisbook(Db *sql.DB, isbn string) string {
	var bookis string
	err := Db.QueryRow("SELECT Bookisbn FROM BookInformation WHERE Bookisbn=?", isbn).Scan(&bookis)
	if err != nil {
		log.Error(err)
		return ""
	}
	return bookis
}

//获取当前书籍的类别id以及类别名字(内部)
func Getcategory(Db *sql.DB, bookcode string) (cgids int, cgnames string) {
	var cgid int
	var cgname string
	err := Db.QueryRow("SELECT Cglocationid FROM Storagemess WHERE Bookcode=?", bookcode).Scan(&cgid)
	if err != nil {
		log.Fatal(err)
	}
	err = Db.QueryRow("SELECT Bookcategory FROM Classification WHERE Categoryid=?", cgid).Scan(&cgname)
	if err != nil {
		log.Fatal(err)
	}
	return cgid, cgname
}

//根据状态获取该状态下的所有图书的bookcode
func Geteverybookcode(Db *sql.DB, status string) []string {
	var str string
	if len(status) <= 0 {
		str = "SELECT Bookcode FROM Storagemess"
		rows, err := Db.Query(str)
		if err != nil {
			log.Error("查询失败")
		}
		var s []string
		for rows.Next() {
			var Bookisbns string
			rows.Scan(&Bookisbns)
			s = append(s, Bookisbns)
		}
		return s
	} else {
		str = "SELECT Bookcode FROM Storagemess WHERE Status=?"
	}
	rows, err := Db.Query(str, status)
	if err != nil {
		log.Error("查询失败")
	}
	var s []string
	for rows.Next() {
		var Bookisbns string
		rows.Scan(&Bookisbns)
		s = append(s, Bookisbns)
	}
	return s
}

func Getlocation(Db *sql.DB, isbn string) (l []Location) {
	str := "select Locationid,Bookcode from Storagemess where Bookisbn = ? and Status = 0"
	rows, err := Db.Query(str, isbn)
	if err != nil {
		log.Error("查询失败")
	}

	for rows.Next() {
		temp := Location{}
		rows.Scan(&temp.Bookcode, &temp.Locationname)
		l = append(l, temp)
	}

	return
}

// 查询

//查询parentid
func isornoparent(Db *sql.DB, locationid int) int {
	var parentid int
	err := Db.QueryRow("SELECT Parentid FROM Bookshelf WHERE Locationid=?", locationid).Scan(&parentid)
	if err != nil {
		log.Error("查询失败")
	}
	return parentid
}

//模糊查询出书的isbn
func Likefind(Db *sql.DB, data string, field string) []Vaguepst {
	rows, err := Db.Query("SELECT Bookisbn,Bookname,Bookauthor,Bookpic,Bookpublic FROM BookInformation WHERE LOCATE(?,"+field+")>0", data)
	if err != nil {
		log.Error("查询失败")
	}
	var Vaguepss []Vaguepst
	for rows.Next() {
		var Booknms Vaguepst
		rows.Scan(&Booknms.Bookisbn, &Booknms.Bookname, &Booknms.Bookauthor, &Booknms.Bookpic, &Booknms.Bookpublic)
		Vaguepss = append(Vaguepss, Booknms)
	}
	return Vaguepss
}

// 更改一级存储库为父母
func UpdateBookShelf(Db *sql.DB, locationid int) {
	tx, _ := Db.Begin()
	tx.Exec("update Bookshelf set Isparent = 1 where Locationid = ?", locationid)
	tx.Commit()
}

// 添加一个token
func AddToken(token string) {

	conn, _ := redis.Dial("tcp", "192.168.10.200:20025")
	conn.Send("auth", "947607a8")
	defer conn.Close()
	conn.Do("set", token, token)
}

// 获得类别名称通过isbn
func GetCategoryName(Db *sql.DB, isbn string) string {
	var cid string
	var cname string
	Db.QueryRow("SELECT Categoryid FROM BookInformation WHERE Bookisbn=?", isbn).Scan(&cid)
	Db.QueryRow("SELECT Bookcategory FROM Classification WHERE Categoryid=?", cid).Scan(&cname)

	return cname
}
func UpdateStorageStatus(delivery amqp.Delivery) {
	type lendbookmq struct {
		Userid   string
		Bookcode string
		ISBN     string
		Status   int
	}

	l := lendbookmq{}
	json.Unmarshal(delivery.Body, &l)
	fmt.Println(l)
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, url, dbname)
	defer Db.Close()
	boolean, _ := ChangeStatus(Db, l.Bookcode, strconv.Itoa(l.Status))
	if boolean == true {
		log.Info("成功")
		delivery.Acknowledger.Ack(delivery.DeliveryTag, true)
	}

}
