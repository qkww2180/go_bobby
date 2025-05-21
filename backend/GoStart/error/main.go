package main

import (
	"fmt"
	perrors "github.com/pkg/errors"
)

/*
go的error和其他语言的try catch不一样， go语言将错误和异常分开，其他语言 异常
go中认为error是一种值
*/

/**

此课程提供者：微信imax882

+微信imax882
办理会员 课程全部免费看

课程清单：https://leaaiv.cn

全网最全 最专业的 一手课程

成立十周年 会员特惠 速来抢购

**/

func divFunc(a, b int) (int, error) {
	if b == 0 {
		return 0, perrors.New("b can't be zero")
	}
	return a / b, nil
}

func main() {
	var a, b = 1, 0
	ret, err := divFunc(a, b)
	if err != nil {
		fmt.Printf("ret is %d, err is %+v\n", ret, err)
	}
}
