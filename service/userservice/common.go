package main

import "encoding/json"

type SuccessErrorResult struct {
	Msg string `json:"msg"`
}

func GetSuccessResult(msg string) []byte {
	result := &Result{
		Status:  200,
		Success: 0,
		Data: &SuccessErrorResult{
			Msg: msg,
		},
	}
	body, _ := json.Marshal(result)
	return body
}
func GetErrorResult(msg string) []byte {
	result := &Result{
		Status:  100,
		Success: 0,
		Data: &SuccessErrorResult{
			Msg: msg,
		},
	}
	body, _ := json.Marshal(result)
	return body
}
