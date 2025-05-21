package main

import "fmt"

/*
CGO_ENABLED:
	cgo表示go中的工具 这个表示是否禁用cgo
GOOS： 目标操作系统
	mac: darwin
	linux: linux
	windows: windows
GOARCH： 目标操作系统的架构(386, amd64, arm), amd64, mac m1 arm架构的

*/
//有没有办法解决 1. 可以在隔离的go的容器环境中编译，同时还能减少 空间
//多阶段构建
// 执行go的build 我们希望在一个隔离的go环境环境中去执行 ， 这个docker需要是一个有go的环境的docker
func main() {
	fmt.Println("Hello, go")
}
