//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

func initEvent() Event {
	//newxxx， provider  - java 构造函数
	panic(wire.Build(NewEvent, NewGreeter, NewMessage))
}
