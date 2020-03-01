/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-10-23 19:59:08
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/yuwe1/shuxiang/common/reflecter"

	"github.com/yuwe1/shuxiang/common/weber"
)

func Selectfunc(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	// file := "./../../public/userconfig.json"
	// file := "D:\\Go_WorkSpace\\book\\public\\commentconfig.json"
	file, _ := os.Getwd()
	file = file + "/public/commentconfig.json"
	parsename := weber.GetParsepath(r.RequestURI)
	fmt.Println(parsename)

	route := weber.ParseRoute(file, parsename)
	fmt.Println("route:", route)
	//Routers 调用方法的空结构体
	var ruTest *EmptyComment
	crMap := reflecter.GetClient().FuncNameToFuncHandler(ruTest)
	fmt.Println("路径：", crMap[route.Funcname])
	fmt.Println("方法：", crMap[route.Funcname])
	//创建带调用方法时需要传入的参数列表
	parms := []reflect.Value{reflect.ValueOf(r)}
	//使用方法名字符串调用指定方法
	body := crMap[route.Funcname].Call(parms)
	w.Write(body[0].Interface().([]byte))
}

func main() {

	http.HandleFunc("/", Selectfunc)
	http.ListenAndServe("0.0.0.0:8400", nil)

}

// 打印错误信息！
func CheckError(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

//编辑返回结果  成功：以一个参数为 1
func GetResult(i int, any interface{}) []byte {

	res := Result{}
	res.Success = 200
	res.Status = i
	res.Data = any
	body, _ := json.Marshal(&res)
	return body

}

//添加评论
func (emptyComment *EmptyComment) AddComment(r *http.Request) []byte {

	comment := Comment{}

	//获取 远程服务器 获取数据代码
	comment.UserID = r.FormValue("userid")
	comment.Content = r.FormValue("content")
	comment.BookISBN = r.FormValue("bookisbn")
	token := r.FormValue("token")

	comment.CommentID = GetRandomID()
	comment.ParentID = "0"
	comment.Created = time.Now().Unix()
	comment.LikeNum = 0
	comment.Status = 0
	comment.IsReply = 0
	err := comment.AddCommentSQL(token)
	if err != nil {
		data := ResultData{
			Msg: "添加评论失败！",
		}
		return GetResult(0, data)
	} else {
		data := ResultData{
			Msg: "添加评论成功！",
		}
		return GetResult(1, data)
	}

}

//添加回复
func (emptyComment *EmptyComment) AddReply(r *http.Request) []byte {

	comment := Comment{}

	//获取 远程服务器 获取数据代码
	comment.UserID = r.FormValue("userid")
	comment.Content = r.FormValue("content")
	comment.ParentID = r.FormValue("commentid")
	comment.BookISBN = r.FormValue("bookisbn")

	comment.CommentID = GetRandomID()
	comment.Created = time.Now().Unix()
	comment.LikeNum = 0
	comment.Status = 0
	comment.IsReply = 0
	err := comment.AddReplySQL()
	if err != nil {
		data := ResultData{
			Msg: "添加回复失败！",
		}
		return GetResult(0, data)
	} else {
		data := ResultData{
			Msg: "添加回复成功！",
		}
		return GetResult(1, data)
	}

}

//显示某本书的评论 -----------------------------------------------------------------------
func (emptyComment *EmptyComment) ShowComments(r *http.Request) []byte {

	comment := Comment{}

	//获取 远程服务器 获取数据代码
	comment.UserID = r.FormValue("userid")
	comment.BookISBN = r.FormValue("bookisbn")
	var str string
	str = r.FormValue("status")
	comment.Status, _ = strconv.Atoi(str)

	//userids用来获取用户头像------------------------------------------
	os, userids, err := comment.ShowCommentsSQl()
	results := Results{}
	user := User{}
	for i := 0; i < len(os); i++ {

		//遍历评论，获取目前登陆的用户是否为该评论点过赞
		comment.CommentID = os[i].CommentID
		os[i].LikeOrUnlike = comment.IsExist()

		// //模拟HTTP请求 *************************** 获取user头像
		// //  http.Get()

		user.UserID = userids[i]
		body := RequestURL(user, "http://127.0.0.1:8100/user/getinfo") //getUserPic
		// fmt.Println(string(body))

		json.Unmarshal(body, &results)
		user = results.Data
		os[i].UserPic = user.UserPic

	}

	if err != nil {
		return GetResult(0, os)
	} else {
		return GetResult(1, os)
	}

}

//显示已审核或未审核的书评（所有书）
func (emptyComment *EmptyComment) ShowYesOrNo(r *http.Request) []byte {

	comment := Comment{}

	//获取 远程服务器 获取数据代码
	var str string
	str = r.FormValue("status")
	comment.Status, _ = strconv.Atoi(str)

	os, err := comment.ShowYesOrNoSQL()

	results := Results{}
	user := User{}
	for i := 0; i < len(os); i++ {

		//模拟HTTP请求 *************************** 获取 username
		//  http.Get()
		//fmt.Println(comment.UserID)

		user.UserID = os[i].UserID

		body := RequestURL(user, "http://127.0.0.1:8100/user/getinfo")
		// fmt.Println(string(body))

		json.Unmarshal(body, &results)
		user = results.Data
		os[i].UserName = user.UserName
	}

	if err != nil {
		return GetResult(0, os)
	} else {
		return GetResult(1, os)
	}

}

//评论通过
func (emptyComment *EmptyComment) CheckComment(r *http.Request) []byte {

	comment := Comment{}

	//获取 远程服务器 获取数据代码
	comment.CommentID = r.FormValue("commentid")

	err := comment.CheckCommentSQL()
	if err != nil {
		data := ResultData{
			Msg: "评论通过失败！",
		}
		return GetResult(0, data)
	} else {
		data := ResultData{
			Msg: "评论通过成功！",
		}
		return GetResult(1, data)
	}

}

//删除评论
func (emptyComment *EmptyComment) DeleteComment(r *http.Request) []byte {

	comment := Comment{}

	//获取 远程服务器 获取数据代码
	comment.CommentID = r.FormValue("commentid")

	err := comment.DeleteCommentSQL()
	if err != nil {
		data := ResultData{
			Msg: err.Error(),
		}
		return GetResult(0, data)
	} else {
		data := ResultData{
			Msg: "删除评论成功！",
		}
		return GetResult(1, data)
	}

}

//点赞评论
func (emptyComment *EmptyComment) LikeComment(r *http.Request) []byte {

	comment := Comment{}

	//获取 远程服务器 获取数据代码
	comment.CommentID = r.FormValue("commentid")
	comment.UserID = r.FormValue("userid")

	err := comment.LikeCommentSQL()
	if err != nil {
		data := ResultData{
			Msg: "点赞评论失败！",
		}
		return GetResult(0, data)
	} else {
		data := ResultData{
			Msg: "点赞评论成功！",
		}
		return GetResult(1, data)
	}

}

//取消点赞评论
func (emptyComment *EmptyComment) UnLikeComment(r *http.Request) []byte {

	comment := Comment{}

	//获取 远程服务器 获取数据代码
	comment.CommentID = r.FormValue("commentid")
	comment.UserID = r.FormValue("userid")

	err := comment.UnLikeCommentSQL()
	if err != nil {
		data := ResultData{
			Msg: "取消点赞评论失败！",
		}
		return GetResult(0, data)
	} else {
		data := ResultData{
			Msg: "取消点赞评论成功！",
		}
		return GetResult(1, data)
	}

}

//用户是否可以评论
func (emptyComment *EmptyComment) Whether(r *http.Request) []byte {

	whe := Whether{}

	//获取 远程服务器 获取数据代码
	whe.UserID = r.FormValue("userid")
	whe.BookISBN = r.FormValue("bookisbn")
	token := r.FormValue("token")

	exist := whe.Whether(token)
	if exist == 0 {
		data := ResultData{
			Msg: "您已评论过了！不能继续评论！",
		}
		return GetResult(0, data)
	} else {
		data := whe
		return GetResult(1, data)
	}

}
