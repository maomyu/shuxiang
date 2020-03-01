package dingding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type AppMessage struct {
	AppKey    string `json:"appkey"`
	Appsecret string `json:"appsecret"`
}
type AccessToken struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Access_token string `json:"access_token"`
}

func GetAppMessage(file string) []AppMessage {
	// 打开文件
	jsonFile, _ := os.Open(file)
	defer jsonFile.Close()
	// 读取文件
	jsonData, _ := ioutil.ReadAll(jsonFile)
	var app []AppMessage
	json.Unmarshal(jsonData, &app)
	return app

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
type DindingUserResult struct {
	Errcode int             `jsn:"errcode"`
	Unionid string          `json:"unionid"`
	Roles   []DingdingRoles `json:"roles"`
	UserID  string          `json:"userid"`
	Name    string          `json:"name"`
	Avatar  string          `json:"avatar"`
}
type DingdingRoles struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	GroupName string `json:"groupname"`
	Type      int    `json:"type"`
}

func GetUserInfo(userID string) User {
	// userID = "064365012233516404"
	// GET请求的参数
	params := url.Values{}
	Url, _ := url.Parse("https://oapi.dingtalk.com/user/get")
	accesstk := GetAccessToken()
	params.Set("access_token", accesstk)
	params.Set("userid", userID)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, _ := http.Get(urlPath)
	defer resp.Body.Close()
	robots, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(robots))
	var r DindingUserResult
	json.Unmarshal(robots, &r)
	user := User{
		UserID:     r.UserID,
		Username:   r.Name,
		Userage:    20,
		Usermobile: "",
		Userpic:    r.Avatar,
		Ismanager:  0,
	}
	fmt.Println(user)
	return user
}

type Code struct {
	Tmp_auth_code string `json:"tmp_auth_code"`
}

func GetUserInfoBycode(code string) {
	client := &http.Client{}
	signature, data := GetAppSecretXXX("dingoaazziae2gt36ntphp")
	fmt.Println(signature[10:])
	fmt.Println(strconv.FormatInt(data, 10))
	pa := "=" + signature[10:] + "&timestamp=" + strconv.FormatInt(data, 10) + "&accessKey=" + "dingoaazziae2gt36ntphp"
	var c Code
	c.Tmp_auth_code = code
	jsondata, _ := json.Marshal(&c)
	fmt.Println(string(jsondata))

	var jsonStr = []byte(string(jsondata))
	req, err := http.NewRequest("POST", "https://oapi.dingtalk.com/sns/getuserinfo_bycode?signature"+pa, bytes.NewBuffer(jsonStr))
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
	fmt.Println(string(body))
}

func GetAccessToken() string {

	// GET请求的参数
	params := url.Values{}
	Url, _ := url.Parse("https://oapi.dingtalk.com/gettoken")
	file, _ := os.Getwd()
	file = file + "/../public/app.json"
	fmt.Println(file)
	app := GetAppMessage(file)
	// app := GetAppMessage("F:/GoWeb/src/shuxiang/public/app.json")

	fmt.Println(app)
	params.Set("appkey", app[0].AppKey)
	params.Set("appsecret", app[0].Appsecret)
	fmt.Println(app)
	// 对中文进行处理
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, _ := http.Get(urlPath)
	defer resp.Body.Close()
	robots, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(robots))
	var actoken AccessToken
	json.Unmarshal(robots, &actoken)
	return actoken.Access_token
}

// 获得加密后的内容
func GetAppSecretXXX(appid string) (urlencode string, data int64) {
	file, _ := os.Getwd()
	fmt.Println(file)
	file = file + "/../public/app.json"
	app := GetAppMessage(file)
	fmt.Println(app)
	fmt.Println(appid)
	fmt.Println(app[1].AppKey)
	if appid == app[1].AppKey {
		// 获取当前的时间戳
		now := time.Now()
		nanos := now.UnixNano()
		// 获得毫秒数
		millis := nanos / 1000000
		// millis = 1546084445901
		// h := hmac.New(sha256.New, []byte("testappSecret"))
		h := hmac.New(sha256.New, []byte(app[1].Appsecret))

		h.Write([]byte(strconv.FormatInt(millis, 10)))
		sha := base64.StdEncoding.EncodeToString(h.Sum(nil))
		v := url.Values{}
		v.Add("signature", sha)
		urlEncodeSignature := v.Encode()
		return urlEncodeSignature, millis
	}
	return "", 1

}
