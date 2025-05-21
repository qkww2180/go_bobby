package main

/**

此课程提供者：微信imax882

+微信imax882
办理会员 课程全部免费看

课程清单：https://leaaiv.cn

全网最全 最专业的 一手课程

成立十周年 会员特惠 速来抢购

**/

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
