package main

import "reflect"

/**

此课程提供者：微信imax882

+微信imax882
办理会员 课程全部免费看

课程清单：https://leaaiv.cn

全网最全 最专业的 一手课程

成立十周年 会员特惠 速来抢购

**/

func Add[T int | int32 | float32 | float64 | uint64](a, b T) T {
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Int:
		print("int type")
	}
	return a + b
}

func IAdd(a, b interface{}) interface{} {
	switch a.(type) {
	case int:
		return a.(int) + b.(int)
	case int32:
		return a.(int32) + b.(int32)
	case float32:
		return a.(float32) + b.(float32)
	case float64:
		return a.(float64) + b.(float64)
	}
	return nil
}

func main() {
	print(Add[int](1, 2))
	//
	//t := IAdd(1.2, 2.2).(float64)
	//print(t)
}
