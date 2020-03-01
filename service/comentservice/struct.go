package main

//获取数据的结构体
type Comment struct {
	CommentID string `json:"commentid"` //评论ID
	Content   string `json:"content"`   //评论内容
	ParentID  string `json:"parentid"`  //父评论ID (0代表一级评论)
	UserID    string `json:"userid"`    //用户ID
	Created   int64  `json:"created"`   //创建时间
	BookISBN  string `json:"bookid"`    //评论的书本ISBN
	LikeNum   int    `json:"likenum"`   //点赞数量
	Status    int    `json:"status"`    //是否通过审核(1通过审核/0未审核)
	IsReply   int    `json:"isreply"`   //是否拥有回复(1有回复/0无回复)
}

type Other struct {
	CommentID    string `json:"commentid"`    //评论ID
	Content      string `json:"content"`      //评论内容
	ParentID     string `json:"parentid"`     //父评论ID (0代表一级评论)
	UserPic      string `json:"userpic"`      //图片地址 url\
	IsReply      int    `json:"isreply"`      //是否拥有回复(1有回复/0无回复)
	LikeNum      int    `json:"likenum"`      //点赞数量
	Created      int64  `json:"created"`      //创建时间
	LikeOrUnlike int    `json:"likeorunlike"` //1代表用户已点赞，0代表未点赞
}

type OtherYesOrNo struct {
	UserID    string `json:"userid"`    //用户ID
	CommentID string `json:"commentid"` //评论ID
	UserName  string `json:"username"`  //用户名
	Content   string `json:"content"`   //评论内容
}

//返回的结构体  结果
type Result struct {
	Status  int         `json:"status"`
	Success int         `json:"success"`
	Data    interface{} `json:"data"`
}

//用户数据结构体
type User struct {
	UserID   string `json:"userid"`
	UserPic  string `json:"userpic"`
	UserName string `json:"username"`
}

//返回的结构体  结果
type Results struct {
	Data User `json:"data"`
}

//调用方法的空结构体
type EmptyComment struct{}

type ResultData struct {
	Msg string `json:"msg"`
}

//获取数据的结构体
type Whether struct {
	UserID   string `json:"userid"` //用户ID
	BookISBN string `json:"bookid"` //评论的书本ISBN
}
