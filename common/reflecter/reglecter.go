package reflecter

import (
	"fmt"
	"reflect"
)

//定义控制器函数Map类型，便于后续快捷使用
type ControllerMapsType map[string]reflect.Value

//声明控制器函数Map类型变量
var ControllerMaps ControllerMapsType

var crMap = make(ControllerMapsType, 100)

type ReflectInterface interface {
	FuncNameToFuncHandler(v interface{}) ControllerMapsType
}
type Reflecter struct {
}

func GetClient() *Reflecter {
	v := new(Reflecter)
	return v
}
func (s *Reflecter) FuncNameToFuncHandler(v interface{}) ControllerMapsType {
	vft := reflect.TypeOf(v)
	vf := reflect.ValueOf(v)

	// vft := vf.Type()
	//读取方法数量
	mNum := vft.NumMethod()
	fmt.Println("NumMethod:", mNum)
	//遍历路由器的方法，
	for i := 0; i < mNum; i++ {
		// vft.Method(i)
		// vft.Method(i)
		mName := vft.Method(i).Name
		crMap[mName] = vf.Method(i)
	}
	return crMap
}
