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
	"sync"
)

//子goroutine如何通知到主的goroutine自己结束了， 主的goroutine如何知道子的goroutine已经结束了

func main(){
	var wg sync.WaitGroup

	//我要监控多少个goroutine执行结束
	wg.Add(100)
	for i := 0; i<100; i++ {
		go func(i int) {
			defer wg.Done()

			fmt.Println(i)
		}(i)
	}

	//等到
	wg.Wait()
	fmt.Println("all done")

	//waitgroup主要用于goroutine的执行等到， Add方法要和Done方法配套
}
