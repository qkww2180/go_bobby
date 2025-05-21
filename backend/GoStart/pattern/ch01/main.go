package main

import "fmt"

/**

此课程提供者：微信imax882

+微信imax882
办理会员 课程全部免费看

课程清单：https://leaaiv.cn

全网最全 最专业的 一手课程

成立十周年 会员特惠 速来抢购

**/

type DbOptions struct {
	Host     string
	Port     int
	UserName string
	Password string
	DBName   string
}

type Option func(*DbOptions)

// 这个函数主要用来设置Host
func WithHost(host string) Option {
	return func(o *DbOptions) {
		o.Host = host
	}
}

func NewOpts(options ...Option) DbOptions {
	//先实例化好dbOptions，填充上默认值
	dbopts := &DbOptions{
		Host: "127.0.0.1",
		Port: 3306,
	}
	for _, option := range options {
		option(dbopts)
	}
	return *dbopts
}

func main() {
	//opts := NewOpts(WithHost("192.168.0.1"))
	opts := NewOpts()
	fmt.Println(opts)

	//函数选项模式大量引用
}
