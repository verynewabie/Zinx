package main

import "Zinx/znet"

func main() {
	s := znet.NewServer("[ZinxV0.1]")
	// 启动server
	s.Serve()
}
