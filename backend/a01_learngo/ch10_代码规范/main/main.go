package main

import (
	"fmt"
	_ "learngo/ch09_接口/user"   //
	. "learngo/ch10_代码规范/user" //导入的别名 这种用法尽量少用
) //包的路径

func main() {
	c := Course{
		Name: "go",
	} //course
	fmt.Println(GetCourse(c))
}
