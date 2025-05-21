package main

import (
	"context"
	"fmt"
	"time"

	"errors"
	"golang.org/x/sync/errgroup"
)

/**

此课程提供者：微信imax882

+微信imax882
办理会员 课程全部免费看

课程清单：https://leaaiv.cn

全网最全 最专业的 一手课程

成立十周年 会员特惠 速来抢购

**/

func main() {
	//errgroup的go方法内部会启动一个goroutine
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		fmt.Println("doing task1")
		time.Sleep(5 * time.Second)
		return errors.New("task1 error")
	})

	eg.Go(func() error {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("doing task2")
			case <-ctx.Done():
				fmt.Println("task2 canceled")
				return ctx.Err()
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("doing task3")
			case <-ctx.Done():
				fmt.Println("task3 canceled")
				return ctx.Err()
			}
		}
	})

	err := eg.Wait()
	if err != nil {
		fmt.Println("task failed")
	} else {
		fmt.Println("task success")
	}
}
