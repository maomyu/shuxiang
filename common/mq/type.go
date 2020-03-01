package mq

import "log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
	}
}

//mq服务定义的包
type MsgInterface interface{}

type Message struct {
}
