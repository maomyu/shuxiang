/*
 * @Description: In User Settings Edit
 * @Author: your name
 * @Date: 2019-08-17 15:52:08
 * @LastEditTime: 2019-10-23 11:23:18
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"errors"
	"fmt"

	"github.com/yuwe1/shuxiang/common/dber"

	"github.com/garyburd/redigo/redis"
)

var (
	username = "root"
	password = "03354ab3"
	urls     = "192.168.10.200:20024"
	dbname   = "comment"
)

var conn redis.Conn

// 连接数据库
func init() {
	var err error

	conn, err = redis.Dial("tcp", "192.168.10.200:20025")
	conn.Send("auth", "947607a8")

	CheckError("连接不到Redis！", err)
}

//添加评论
func (com *Comment) AddCommentSQL(token string) (err error) {

	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	_, err = tx.Exec("insert into comment(commentid, content, parentid, userid, created, bookisbn, likenum, status, isreply) values (?, ?, ?, ?, ?, ?, ?, ?, ?) ", com.CommentID, com.Content, com.ParentID, com.UserID, com.Created, com.BookISBN, com.LikeNum, com.Status, com.IsReply)
	CheckError("添加记录失败！\n", err)
	tx.Commit()
	_, err = conn.Do("DEL", token)
	CheckError("删除token失败！\n", err)
	//defer conn.Close()
	return

}

//添加回复
func (com *Comment) AddReplySQL() (err error) {

	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	_, err = tx.Exec("insert into comment(commentid, content, parentid, userid, created, bookisbn, likenum, status, isreply) values (?, ?, ?, ?, ?, ?, ?, ?, ?) ", com.CommentID, com.Content, com.ParentID, com.UserID, com.Created, com.BookISBN, com.LikeNum, com.Status, com.IsReply)
	CheckError("添加记录失败！\n", err)
	//更新 父评论 的 isreply 为拥有回复
	_, err = tx.Exec("update comment set isreply = 1  where commentid = ?", com.CommentID)
	CheckError("更新父评论拥有回复失败！\n", err)
	tx.Commit()
	return

}

//获取一本书的评论/回复
func (com *Comment) ShowCommentsSQl() (os []Other, userids []string, err error) {

	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	rows, err := tx.Query("select commentid, content, parentid, userid, created, likenum, isreply from comment where bookisbn = ? and status = ?", com.BookISBN, com.Status)
	defer rows.Close()
	CheckError("查询失败！\n", err)
	o := Other{}
	var userid string
	for rows.Next() {

		rows.Scan(&o.CommentID, &o.Content, &o.ParentID, &userid, &o.Created, &o.LikeNum, &o.IsReply)
		os = append(os, o)
		userids = append(userids, userid)
	}
	return

}

// 显示已审核/未审核的书评
func (com *Comment) ShowYesOrNoSQL() (os []OtherYesOrNo, err error) {

	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	rows, err := tx.Query("select commentid, content, userid from comment where status = ?", com.Status)
	defer rows.Close()
	CheckError("查询失败！\n", err)
	o := OtherYesOrNo{}
	for rows.Next() {
		rows.Scan(&o.CommentID, &o.Content, &o.UserID)
		os = append(os, o)
	}

	return

}

//评论通过
func (com *Comment) CheckCommentSQL() (err error) {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	_, err = tx.Exec("update comment set status = ? where commentid = ?", 1, com.CommentID)
	CheckError("修改status失败！\n", err)
	tx.Commit()
	return
}

//删除评论
func (com *Comment) DeleteCommentSQL() (err error) {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	row, err := tx.Exec("delete from comment where commentid = ? and status = 0", com.CommentID)
	CheckError("删除评论失败！\n", err)
	tx.Commit()
	if affect, _ := row.RowsAffected(); affect == 0 {
		return errors.New("该评论已审核通过！（审核通过的评论不能删除）")
	}
	return
}

//点赞评论
func (com *Comment) LikeCommentSQL() (err error) {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	_, err = tx.Exec("update comment set likenum = likenum+1  where commentid = ?", com.CommentID)
	CheckError("点赞评论失败！\n", err)
	_, err = tx.Exec("insert into likenum (userid, commentid) values(?, ?)", com.UserID, com.CommentID)
	CheckError("添加点赞表记录失败！\n", err)
	tx.Commit()
	return
}

//取消点赞评论
func (com *Comment) UnLikeCommentSQL() (err error) {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	_, err = tx.Exec("update comment set likenum = likenum-1  where commentid = ?", com.CommentID)
	CheckError("取消点赞评论失败！\n", err)
	_, err = tx.Exec("delete from likenum where userid = ? and commentid = ?", com.UserID, com.CommentID)
	CheckError("删除点赞表记录失败！\n", err)
	tx.Commit()
	return
}

//用户为该该评论是否点过赞
func (com *Comment) IsExist() int {
	c := dber.GetClient()
	Db := c.ConnectTry(username, password, urls, dbname)
	defer Db.Close()
	tx, err := Db.Begin()
	CheckError("开启事务失败！\n", err)
	var str string
	err = tx.QueryRow("select userid from likenum where userid = ? and commentid = ?", com.UserID, com.CommentID).Scan(&str)
	if err != nil {
		return 0
	}
	return 1
}

//用户是否可以评论
func (whe *Whether) Whether(token string) int {

	//检查是否存在key值
	exists, err := redis.Int(conn.Do("EXISTS", token))
	CheckError("Redis查询token失败！", err)

	fmt.Printf("查询的结果为：", exists)
	return exists
}
