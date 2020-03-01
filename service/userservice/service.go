package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yuwe1/shuxiang/common/dber"

	"github.com/streadway/amqp"

	_ "github.com/go-sql-driver/mysql"
)

type UserLend struct {
	UserID     string    `json:"userID"`
	Bookcode   string    `json:"bookcode"`
	Lenddate   time.Time `json:"lenddate"`
	Returndate time.Time `json:"returndate"`
	Status     int       `json:"status"`
	Bookisbn   string    `json:"bookisbn"`
}
type UserLendMQ struct {
	UserID   string `json:"userid"`
	Bookcode string `json:"bookcode"`
	IsBN     string `json:"ISBN"`
	Status   int    `json:"status"`
}

var (
	username = "root"
	password = "03354ab3"
	urls     = "192.168.10.200:20024"
	dbname   = "users"
)

// 根据userid查看用户所有的书
func GetLendBookById(userid string) (alllendbook []UserLend) {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	rows, _ := d.Query("select * from lendbook where userid = ?", userid)
	for rows.Next() == true {
		var lendbook UserLend
		rows.Scan(&lendbook.UserID,
			&lendbook.Bookcode,
			&lendbook.Lenddate,
			&lendbook.Returndate,
			&lendbook.Status,
			&lendbook.Bookisbn,
		)
		lendbook.Lenddate.Format("2006-01-02 15:04:05.999999999 -0700 MST")
		lendbook.Returndate.Format("2006-01-02 15:04:05.999999999 -0700 MST")
		alllendbook = append(alllendbook, lendbook)

	}
	return alllendbook
}

// 删除预约记录
func DelAppointmentDoc(userID string, bookcode string) bool {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("delete from lendbook where userid = ? and bookcode = ?", userID, bookcode)
	err := tx.Commit()
	if err != nil {
		return false
	}
	return true
}

// 新增预约记录
func AddLendDoc(delivery amqp.Delivery) {

	// 解析内容
	var userlendmq UserLendMQ
	json.Unmarshal(delivery.Body, &userlendmq)
	fmt.Println(userlendmq)
	lenddate := time.Now()
	r, _ := time.ParseDuration("240h")
	returndate := lenddate.Add(r)
	userlend := UserLend{
		UserID:     userlendmq.UserID,
		Bookcode:   userlendmq.Bookcode,
		Lenddate:   lenddate,
		Returndate: returndate,
		Status:     userlendmq.Status,
		Bookisbn:   userlendmq.IsBN,
	}

	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("insert into lendbook(userid,bookcode,lenddate,returndate,status,bookisbn) values(?,?,?,?,?,?)",
		userlend.UserID,
		userlend.Bookcode,
		userlend.Lenddate,
		userlend.Returndate,
		userlend.Status,
		userlend.Bookisbn,
	)
	err := tx.Commit()
	if err != nil {

	}

	// delivery.Acknowledger.Ack(delivery.DeliveryTag, false)
	delivery.Acknowledger.Ack(delivery.DeliveryTag, true)
	// fmt.Println("应答成功")

}

// 根据bookcode查询正在借阅的书返回userid
func GetUserIDByBookCode(bookcode string) User {
	var userid string
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	err := d.QueryRow("select userid from lendbook where status =1 and bookcode=?", bookcode).Scan(&userid)
	if err != nil {
		return User{}
	}
	user := GetUserInfo(userid)
	return user
}

// 更新预约记录，用户成功借阅书籍后通知用户
func UpdateLendDoc(delivery amqp.Delivery) {
	// 解析内容
	var userlendmq UserLendMQ
	json.Unmarshal(delivery.Body, &userlendmq)

	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()

	tx.Exec("update lendbook set status = ? where userid = ? and bookcode = ?", userlendmq.Status,
		userlendmq.UserID,
		userlendmq.Bookcode,
	)
	if tx.Commit() != nil {

	}
	delivery.Acknowledger.Ack(delivery.DeliveryTag, true)
}

//验证当前用户是否已经过了预约时间
func CheckBooktime(code string, userid string, isbn string) bool {
	var oldlenddate time.Time
	// 获取当前时间
	lenddate := time.Now()
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	d.QueryRow("select lenddate from lendbook where userid=? and bookcode = ?", userid, code).Scan(&oldlenddate)
	if lenddate.Hour() == oldlenddate.Hour() {
		if lenddate.Minute()-oldlenddate.Minute() <= 15 {
			// 说明在有效时间
			lendbookmq := &UserLendMQ{
				UserID:   userid,
				Bookcode: code,
				IsBN:     isbn,
				Status:   1,
			}
			body, _ := json.Marshal(lendbookmq)
			Client.PublishOnQueue(body, "UpdateLendDoc", "UpdateLendDoc")
			Client.PublishOnQueue(body, "UpdateStorageStatus", "UpdateStorageStatus")
			return true
		}
	} else if lenddate.Hour()-oldlenddate.Hour() == 1 {
		if (lenddate.Minute() + (60 - oldlenddate.Minute())) <= 15 {
			// 说明在有效时间
			lendbookmq := &UserLendMQ{
				UserID:   userid,
				Bookcode: code,
				IsBN:     isbn,
				Status:   1,
			}
			body, _ := json.Marshal(lendbookmq)
			Client.PublishOnQueue(body, "UpdateLendDoc", "UpdateLendDoc")
			Client.PublishOnQueue(body, "UpdateStorageStatus", "UpdateStorageStatus")
			return true
		}
	}
	// 删除预约记录
	DelAppointmentDoc(userid, code)
	lendbookmq := &UserLendMQ{
		UserID:   userid,
		Bookcode: code,
		IsBN:     isbn,
		Status:   0,
	}
	body, _ := json.Marshal(lendbookmq)
	Client.PublishOnQueue(body, "UpdateStorageStatus", "UpdateStorageStatus")
	return false
}

// 用户结构体
type User struct {
	UserID     string `json:"userID"`
	Username   string `jsson:"username"`
	Userage    int    `json:"userage"`
	Usermobile string `json:"usermobile"`
	Userpic    string `json:"userpic"`
	Ismanager  int    `json:"ismanager"`
}

// 获得用户头像
func GetUserIcon(userid string) string {
	var name string
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	err := d.QueryRow("select userpic from user where userid=?", userid).Scan(&name)
	if err != nil {
		return ""
	}
	return name
}

// 获得用户信息
func GetUserInfo(userid string) User {
	var user User
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	err := d.QueryRow("select * from user where userid = ?", userid).Scan(&user.UserID,
		&user.Username,
		&user.Userage,
		&user.Usermobile,
		&user.Userpic,
		&user.Ismanager,
	)
	if err != nil {
		fmt.Println(err)
		return User{}
	}
	return user
}
func GetUsers() (users []User) {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	rows, _ := d.Query("select * from user")
	for rows.Next() == true {
		var user User
		rows.Scan(&user.UserID,
			&user.Username,
			&user.Userage,
			&user.Usermobile,
			&user.Userpic,
			&user.Ismanager,
		)
		users = append(users, user)
	}
	return users
}
func AddUserToManger(userid string) {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("update user set ismanager = 1 where userid  = ?", userid)
	tx.Commit()
}
func AddManagerTouser(userid string) {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("update user set ismanager = 0 where userid  = ?", userid)
	tx.Commit()
}

// 新填一个以用户
func SaveUser(user User) bool {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("insert into user(userid,username,userage,usermobile,userpic,ismanager) values(?,?,?,?,?,?)", &user.UserID,
		&user.Username,
		&user.Userage,
		&user.Usermobile,
		&user.Userpic,
		&user.Ismanager,
	)
	err := tx.Commit()
	if err != nil {
		return false
	}
	return true
}

type GradeResult struct {
	UserID string `json:"userid"`
	Num    int    `json:"num"`
	Typel  int    `json:"typel"`
}

// 积分结构体
type Grade struct {
	Integralid   string `json:"integralid"`
	Integragrade int    `json:"integragrade"`
	Usergrade    int    `json:"usergrade"`
	UserID       string `json:"userid"`
}

// 更新积分
// 更新积分为其他服务通知mq时增加
// 当进行登录，签到等操作时
func AddGradenum(delivery amqp.Delivery) {
	var g GradeResult
	json.Unmarshal(delivery.Body, &g)
	// 更新之前的积分
	num1 := GetUserIntegrade(g.UserID)
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("update grade set integragrade = integragrade + ? where userid = ?", g.Num, g.UserID)
	tx.Commit()
	// 更新之后的积分
	num2 := GetUserIntegrade(g.UserID)
	if (num2/100)-(num1/100) == 1 {
		// 说明达到升级的需要需要更新等级
		UpdateGrade(g.UserID)
	}
	delivery.Acknowledger.Ack(delivery.DeliveryTag, true)
}

// 插入一条积分记录
// 当有新的用户进来的时候即新添加一个用户时
func InsertGrade(grade Grade) {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("insert into grade(integralid,integragrade,usergrade,userid) values(?,?,?,?)", grade.Integralid, grade.Integragrade, grade.Usergrade, grade.UserID)
	tx.Commit()
}

// 获得当前用户的积分
func GetUserIntegrade(userid string) int {
	var num int
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	err := d.QueryRow("select integragrade from grade where userid=?", userid).Scan(&num)
	if err != nil {

	}
	return num
}

// 获取当前用户的等级
func GetUserGrade(userid string) int {
	var num int
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	err := d.QueryRow("select usergrade from grade where userid=?", userid).Scan(&num)
	if err != nil {
		fmt.Println(err)
	}
	return num
}
func InsertPunch(userid string) (punchnum int, iscontinue int) {
	// 查询是否有该用户
	n := time.Now()
	t, num := SelectPunchById(userid)
	if num != 0 {

		cdt := t.Day()

		cdn := n.Day()

		if cdn-cdt >= 2 {
			// 更新记录签到天数为1，
			fmt.Println("没有连续签到：", cdn-cdt)
			updatePunch(userid, n, 1)
			return 1, 0
		} else if cdn-cdt == 1 {
			updatePunch(userid, n, num+1)
			return num + 1, 1
		} else {
			updatePunch(userid, n, num+1)
			return num + 1, 1
		}
	}
	fmt.Println("保存新的签到记录")
	fmt.Println(savePunch(userid, n, 1))
	return 1, 0
}
func savePunch(userid string, t time.Time, num int) bool {
	fmt.Println(userid)
	// num = 0
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("insert into punch(userid,date,num) values(?,?,?)", userid, t, num)
	err := tx.Commit()
	if err != nil {
		return false
	}
	return true
}

// 更新签到记录
func updatePunch(userid string, t time.Time, num int) bool {
	// num = 0
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("update punch set date = ?,num = ?  where userid = ?", t, num, userid)
	err := tx.Commit()
	if err != nil {
		return false
	}
	return true
}

func SelectPunchById(userid string) (t time.Time, num int) {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()

	err := d.QueryRow("select date,num from punch where userid=?", userid).Scan(&t, &num)
	if err != nil {

	}
	t.Format("2006-01-02 15:04:05.999999999 -0700 MST")
	return t, num
}

// 更新用户等级
//在每次积分更新的时候需要检验积分的情况
func UpdateGrade(userid string) bool {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	tx, _ := d.Begin()
	tx.Exec("update grade set usergrade = usergrade + 1 where userid = ?", userid)
	err := tx.Commit()
	if err != nil {
		return false
	}
	return true
}

// 插入一条等级记录
// func InsertGrade(userid string){
// 	c := dber.GetClient()
// 	d := c.ConnectTry(username, password, urls, dbname)
// 	defer d.Close()
// 	tx, _ := d.Begin()
// 	tx.Exec("insert into grade(integralid,integragrade,usergrade,userid) values(?,1,100,?)", userid,userid)
// 	err := tx.Commit()
// 	if err != nil {

// 	}
// }

// 馆中当前被借阅的数量或者所有超时未归还的数量
func GetLendNum(status int) int {
	var num int
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	err := d.QueryRow("select count(*) from lendbook where status = ?", status).Scan(&num)
	if err != nil {

	}

	return num
}

// 全年阅读量
func GetYearReadNum() int {
	var num int
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	err := d.QueryRow("select count(*) from lendbook").Scan(&num)
	if err != nil {

	}
	return num
}

// 查询正在借阅的所有书籍
func GetLendingBook() (alllendbook []UserLend) {

	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	rows, _ := d.Query("select * from lendbook where status = 1")
	for rows.Next() == true {
		var lendbook UserLend
		rows.Scan(&lendbook.UserID,
			&lendbook.Bookcode,
			&lendbook.Lenddate,
			&lendbook.Returndate,
			&lendbook.Status,
			&lendbook.Bookisbn,
		)
		lendbook.Lenddate.Format("2006-01-02 15:04:05.999999999 -0700 MST")
		lendbook.Returndate.Format("2006-01-02 15:04:05.999999999 -0700 MST")
		alllendbook = append(alllendbook, lendbook)

	}
	return alllendbook
}

// 查询超时的所有的图书
func GetTTlBook() (alllendbook []UserLend) {
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	rows, _ := d.Query("select * from lendbook where status = 3")
	for rows.Next() == true {
		var lendbook UserLend
		rows.Scan(&lendbook.UserID,
			&lendbook.Bookcode,
			&lendbook.Lenddate,
			&lendbook.Returndate,
			&lendbook.Status,
			&lendbook.Bookisbn,
		)
		lendbook.Lenddate.Format("2006-01-02 15:04:05.999999999 -0700 MST")
		lendbook.Returndate.Format("2006-01-02 15:04:05.999999999 -0700 MST")
		alllendbook = append(alllendbook, lendbook)

	}
	return alllendbook
}

// 获得书籍的借阅时间
func GetLendtime(code string, userid string) time.Time {
	var oldlenddate time.Time
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()

	d.QueryRow("select lenddate from lendbook where userid=? and bookcode = ?", userid, code).Scan(&oldlenddate)
	fmt.Println(oldlenddate)
	oldlenddate.Format("2006-01-02 15:04:05.999999999 -0700 MST")
	return oldlenddate
}

// 查看某用户对某本书的借阅状态
func GetLendstatus(code string, status int) (status1 int, userid string) {

	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	d.QueryRow("select status,userid from lendbook where status=? and bookcode = ?", status, code).Scan(&status1, &userid)
	return status1, userid
}

// 根据userid和bookcode获取状态
func GetLendstatusByID(code string, userid string) int {
	var status int
	c := dber.GetClient()
	d := c.ConnectTry(username, password, urls, dbname)
	defer d.Close()
	d.QueryRow("select status from lendbook where userid=? and bookisbn = ?", userid, code).Scan(&status)
	return status
}

// 计算下一等级所需要的经验
func Getdayadd(data int) int {
	a := data / 10000 % 10
	b := data / 1000 % 10
	c := data / 100 % 10
	d := data / 10 % 10
	e := data / 1 % 10
	fmt.Println(a, b, c, d, e)
	return 100 - (d*10 + e)
}
func GetDate(data int) int {
	return data / 10
}
