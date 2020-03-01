package main

import (
	"math/rand"
	"time"
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
