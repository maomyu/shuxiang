package dber

import (
	"database/sql"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetRandomID() string {

	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)

}

type DbClientInterface interface {
	Connect(user string, password string, url string, dbname string) *DBClient
}
type DBClient struct {
}

func GetClient() *DBClient {
	var db = new(DBClient)
	return db
}

func (db *DBClient) ConnectTry(user string, password string, url string, dbname string) *sql.DB {
	DB_connect_string := user + ":" + password + "@tcp(" + url + ")/" + dbname + "?charset=utf8&loc=Asia%2FShanghai&parseTime=true"
	d, _ := sql.Open("mysql", DB_connect_string)
	return d
}
