/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-10-23 18:10:50
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/yuwe1/shuxiang/common/mq"
	"github.com/yuwe1/shuxiang/common/reflecter"
	"github.com/yuwe1/shuxiang/common/weber"
)

var Client mq.MessagingClient

// 初始化rabbitmq
func init() {
	Client.Conn = Client.ConnectToRabbitmq("amqp://admin:admin@192.168.10.200:20026")
}

func Selectfunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// file := "./../../public/userconfig.json"
	//file:="../"

	//file := "F:/GoWeb/src/shuxiang/public/userconfig.json"
	file, _ := os.Getwd()
	fmt.Println(file)
	file = file + "/public/userconfig.json"
	parsename := weber.GetParsepath(r.RequestURI)
	fmt.Println(parsename)

	route := weber.ParseRoute(file, parsename)
	fmt.Println("route:", route)
	var ruTest *Routers
	crMap := reflecter.GetClient().FuncNameToFuncHandler(ruTest)
	fmt.Println("路径：", route.Funcname)
	fmt.Println("方法：", crMap[route.Funcname])
	//创建带调用方法时需要传入的参数列表
	parms := []reflect.Value{reflect.ValueOf(r)}
	//使用方法名字符串调用指定方法
	body := crMap[route.Funcname].Call(parms)
	w.Write(body[0].Interface().([]byte))
}

//定义路由器结构类型
type Routers struct {
}
type ResultInterface interface {
}
type Result struct {
	Status  int             `json:"status"`
	Success int             `json:"success"`
	Data    ResultInterface `json:"data"`
}
type LoginResult struct {
	UserID    string `json:"userID"`
	Username  string `json:"username"`
	Userpic   string `json:"userpic"`
	Usergrade int    `json:"usergrade"`
	Ismanager int    `json:"ismanager"`
}

func (this *Routers) Getinfo(r *http.Request) []byte {
	userid := r.FormValue("userid")
	user := GetUserInfo(userid)

	result := &Result{
		Status:  200,
		Success: 1,
		Data:    &user,
	}
	body, _ := json.Marshal(result)
	return body
}

// 用户登录
func (this *Routers) LoginEApp(r *http.Request) []byte {
	userid := r.FormValue("userid")
	username := r.FormValue("username")
	userage := r.FormValue("userage")
	usermobile := r.FormValue("usermobile")
	userpic := r.FormValue("userpic")
	// 检查数据库中是否已经存在该用户
	user := GetUserInfo(userid)
	fmt.Println(user)
	if len(userid) > 0 {
		if len(user.UserID) <= 0 {
			// 说明不存在该用户，保存新的用户信息
			// duser := dingding.GetUserInfo(userid)
			fmt.Println("不存在该用户")
			result := &Result{
				Status:  200,
				Success: 1,
				Data: &LoginResult{
					UserID:   userid,
					Username: username,
					Userpic:  userpic,
					// 新用户默认为1
					Usergrade: 1,
					Ismanager: 0,
				},
			}
			body, _ := json.Marshal(result)
			// 保存新的用户信息
			user.UserID = userid
			age, _ := strconv.Atoi(userage)
			user.Userage = age
			user.Usermobile = usermobile
			user.Username = username
			user.Userpic = userpic
			user.Ismanager = 0
			SaveUser(user)
			grade := Grade{
				Integralid:   userid,
				Integragrade: 100,
				Usergrade:    1,
				UserID:       userid,
			}
			InsertGrade(grade)
			return body
		} else {
			grade := GetUserGrade(user.UserID)
			fmt.Println(grade)
			// 查询用户的等级
			result := &Result{
				Status:  200,
				Success: 1,
				Data: &LoginResult{
					UserID:   user.UserID,
					Username: user.Username,
					Userpic:  user.Userpic,
					// 新用户默认为1
					Usergrade: grade,
					Ismanager: user.Ismanager,
				},
			}
			body, _ := json.Marshal(result)
			return body
		}
	} else {
		return GetErrorResult("用户名出错")
	}

	return GetErrorResult("用户名出错")
}

// 管理员登录
func (this *Routers) Login(r *http.Request) []byte {
	// 获得参数userid
	userid := r.FormValue("userid")

	// // 获得个人信息
	// dingding.GetUserInfoBycode(code)
	user := GetUserInfo(userid)
	if user.Ismanager == 1 {
		grade := GetUserGrade(user.UserID)
		// 查询用户的等级
		result := &Result{
			Status:  200,
			Success: 1,
			Data: &LoginResult{
				UserID:   user.UserID,
				Username: user.Username,
				Userpic:  user.Userpic,
				// 新用户默认为1
				Usergrade: grade,
				Ismanager: user.Ismanager,
			},
		}
		body, _ := json.Marshal(result)
		return body
	} else {
		return GetErrorResult("你不是管理员")
	}
	return GetErrorResult("你不是管理员")
}

func GetCuihuanUserID(userids []string) string {

	var result string
	for _, v := range userids {
		result = result + v + ","
	}
	// fmt.Println(result)
	result = result[:len(result)-1]
	return result
}

type Message struct {
	UserID string `json:"userID"`
	Url    string `json:"url"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

// 一键催还
func (this *Routers) Cuihuan(r *http.Request) []byte {
	userids := r.FormValue("userid")
	// 获得所有未归还图书的userid
	alllendbook := GetTTlBook()
	var ttluserids []string
	for _, book := range alllendbook {
		ttluserids = append(ttluserids, book.UserID)
	}
	fmt.Println(userids)

	// 假如已经分号用户id
	t := time.Now()
	formatNow := t.Format("2006-01-02 15:04:05")
	ids := GetCuihuanUserID(ttluserids)
	fmt.Println(ids)
	var m Message
	m.UserID = ids
	m.Url = "test"
	m.Title = "还书提醒"
	m.Text = formatNow + "\n小书箱提醒您要赶快还书哦"
	fmt.Println(m)
	body, _ := json.Marshal(&m)
	Client.PublishOnQueue(body, "SendMessage", "SendMessage")
	return GetSuccessResult("发送成功")
}

type NumResult struct {
	Lendnum     int `json:"lendnum"`
	Ttlnum      int `json:"ttlnum"`
	Yearreadnum int `json:"yearreadnum"`
}

func (this *Routers) Getreadnum(r *http.Request) []byte {

	result := &Result{
		Status:  200,
		Success: 1,
		Data: &NumResult{
			Lendnum:     GetLendNum(1),
			Ttlnum:      GetLendNum(3),
			Yearreadnum: GetYearReadNum(),
		},
	}
	fmt.Println(result)
	body, _ := json.Marshal(result)
	fmt.Println(string(body))
	return body
}
func (this *Routers) Getuserbycode(r *http.Request) []byte {
	bookcode := r.FormValue("boodcode")
	user := GetUserIDByBookCode(bookcode)
	result := &Result{
		Status:  200,
		Success: 1,
		Data: &User{
			UserID: user.UserID,
		},
	}
	body, _ := json.Marshal(result)
	return body
}
func (this *Routers) Lendbook(r *http.Request) []byte {
	bookcode := r.FormValue("code")
	userid := r.FormValue("userid")
	isbn := r.FormValue("isbn")
	if CheckBooktime(bookcode, userid, isbn) {
		return GetSuccessResult("借阅成功")
	} else {
		return GetErrorResult("借阅失败，已经超时")
	}
}
func (this *Routers) Lendbooking(r *http.Request) []byte {
	// dingding.GetUserInfoBycode("8e5d44d56db030b2b74c8b2f1f9396c8")

	books := GetLendingBook()
	result := &Result{
		Status:  200,
		Success: 1,
		Data:    books,
	}
	body, _ := json.Marshal(result)
	return body
}

type StatusResult struct {
	Status int    `json:"status"`
	UserID string `json:"userID"`
}

func (this *Routers) GetUserLendStatus(r *http.Request) []byte {
	code := r.FormValue("bookcode")
	status := r.FormValue("status")
	statuss, _ := strconv.Atoi(status)
	status1, userid := GetLendstatus(code, statuss)
	result := &Result{
		Status:  200,
		Success: 1,
		Data: &StatusResult{
			Status: status1,
			UserID: userid,
		},
	}
	body, _ := json.Marshal(result)
	return body
}
func (this *Routers) GetUserLendStatusByid(r *http.Request) []byte {
	data := r.PostFormValue("data")
	type Data struct {
		Isbn   string
		Userid string
	}
	d := Data{}
	json.Unmarshal([]byte(data), &d)
	isbn := d.Isbn
	userid := d.Userid
	fmt.Println(userid, "   ", isbn)
	status := GetLendstatusByID(isbn, userid)
	fmt.Println(status)
	result := &Result{
		Status:  200,
		Success: 1,
		Data: &StatusResult{
			Status: status,
			UserID: userid,
		},
	}
	body, _ := json.Marshal(result)
	return body
}

type GradeResulto struct {
	Usergrade       int `json:"usergrade"`
	Userexperience  int `json:"userexperience"`
	Leaveexperience int `json:"leaveexperience"`
	Datenum         int `json:"datenum"`
}

func (this *Routers) Showgrade(r *http.Request) []byte {
	userid := r.FormValue("userid")
	grade := GetUserGrade(userid)
	expr := GetUserIntegrade(userid)
	expradd := Getdayadd(expr)
	var g GradeResulto
	g.Usergrade = grade
	g.Userexperience = expr
	g.Leaveexperience = expradd
	g.Datenum = GetDate(expradd)
	result := &Result{
		Status:  200,
		Success: 1,
		Data:    &g,
	}
	body, _ := json.Marshal(result)
	return body
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
type BookShelfResult struct {
	Bookid     string `json:"bookid"`
	Bookname   string `json:"bookname"`
	Pic        string `json:"pic"`
	Expirydate string `json:"expirydate"`
	Status     int    `json:"status"`
}
type BookHhelfResult struct {
	Status  int   `json:"status"`
	Success int   `json:"success"`
	Data    Data1 `json:"data"`
}
type BookcodeResult struct {
	Bookcode string `json:"bookcode"`
}

// 用户查看我的书架
func (this *Routers) Books(r *http.Request) []byte {
	var booksres []BookShelfResult
	userid := r.FormValue("userid")
	// 获得该用户的所有书
	alllend := GetLendBookById(userid)
	// 根据书的bookcode查看书籍详情
	fmt.Println(alllend)
	for _, book := range alllend {
		bookcode := &BookcodeResult{
			Bookcode: book.Bookcode,
		}
		body1, _ := json.Marshal(bookcode)
		// params := url.Values{}
		Url, err := url.Parse("http://127.0.0.1:8100/hub/haveallbooks")
		if err != nil {
			panic(err.Error())

		}

		urlPath := Url.String() + "?data=" + string(body1)
		resp, err := http.Get(urlPath)
		fmt.Println(urlPath)
		defer resp.Body.Close()
		s, err := ioutil.ReadAll(resp.Body)
		var res BookHhelfResult
		fmt.Println(string(s))
		json.Unmarshal(s, &res)
		fmt.Println(res)
		var booksh BookShelfResult
		booksh.Bookid = book.Bookcode
		booksh.Bookname = res.Data.Bookname
		booksh.Pic = res.Data.Bookpic
		booksh.Expirydate = "4"
		booksh.Status = (res.Data.Status)
		booksres = append(booksres, booksh)
	}
	resss := &Result{
		Status:  200,
		Success: 1,
		Data:    booksres,
	}
	body, _ := json.Marshal(resss)
	return body
}

type PunchResult struct {
	Punchnum   int `json:"punchnum"`
	Iscontinue int `json:"iscontinue"`
}

// 用户签到
func (this *Routers) Punch(r *http.Request) []byte {
	userid := r.FormValue("userid")
	if len(userid) <= 0 {
		return GetErrorResult("请输入正确格式的代码")
	}
	// UpdateGrade(userid)

	// 插入一条签到记录
	punchnum, isc := InsertPunch(userid)
	result := &Result{
		Status:  200,
		Success: 1,
		Data: &PunchResult{
			Punchnum:   punchnum,
			Iscontinue: isc,
		},
	}
	graderesult := &GradeResult{
		UserID: userid,
		Num:    10,
	}
	body1, _ := json.Marshal(&graderesult)
	Client.PublishOnQueue(body1, "AddGradenum", "AddGradenum")
	body, _ := json.Marshal(result)
	return body
}

// 将用户转成管理员
func (this *Routers) AddManager(r *http.Request) []byte {
	manageruserid := r.FormValue("manageruserid")
	if len(manageruserid) <= 0 {
		return GetErrorResult("请输入正确格式的代码")
	}
	user := GetUserInfo(manageruserid)
	if user.Ismanager == 1 {
		userid := r.FormValue("userid")
		AddUserToManger(userid)
		return GetSuccessResult("更新成功")
	} else {
		return GetErrorResult("你不是管理员")
	}
	return GetErrorResult("系统开了小差")
}

// 将管理员转变成用户
func (this *Routers) ManagerToUser(r *http.Request) []byte {
	manageruserid := r.FormValue("manageruserid")
	if len(manageruserid) <= 0 {
		return GetErrorResult("请输入正确格式的代码")
	}

	user := GetUserInfo(manageruserid)
	if user.Ismanager == 1 {
		userid := r.FormValue("userid")
		AddManagerTouser(userid)
		return GetSuccessResult("更新成功")
	} else {
		return GetErrorResult("你不是管理员")
	}
	return GetErrorResult("系统开了小差")

}

// 管理员查看所有的用户
func (this *Routers) SelectUsers(r *http.Request) []byte {
	userid := r.FormValue("userid")
	fmt.Println(userid)
	if len(userid) <= 0 {
		return GetErrorResult("请输入正确格式的代码")
	}
	user := GetUserInfo(userid)
	fmt.Println(user)
	if user.Ismanager == 1 {
		users := GetUsers()
		result := &Result{
			Status:  200,
			Success: 1,
			Data:    &users,
		}
		body, _ := json.Marshal(result)
		return body
	} else {
		return GetErrorResult("你不是管理员")
	}
	return GetErrorResult("你不是管理员")
}
func main() {

	go Client.ConsumeFromQueue("UpdateLendDoc", "UpdateLendDoc", UpdateLendDoc)
	go Client.ConsumeFromQueue("AddLendDoc", "AddLendDoc", AddLendDoc)
	go Client.ConsumeFromQueue("AddGradenum", "AddGradenum", AddGradenum)

	fmt.Println("*****************")
	http.HandleFunc("/", Selectfunc)
	http.ListenAndServe("0.0.0.0:8200", nil)
}
