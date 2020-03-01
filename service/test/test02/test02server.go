package main

import (
	"fmt"
	"reflect"
	"shuxiang/common/reflecter"
)

type RouteMethod interface {
	Insert()
	Add()
}

type Router struct {
}

func (v *Router) Insert() {

}

func (v *Router) Add() {

}
func change(r interface{}) {
	vft := reflect.TypeOf(r)
	// vf := reflect.ValueOf(&r)

	// vft := vf.Type()
	//读取方法数量
	mNum := vft.NumMethod()
	fmt.Println("NumMethod:", mNum)
}
func main() {
	var r *Router

	reflecter.GetClient().FuncNameToFuncHandler(r)

	change(r)

}
