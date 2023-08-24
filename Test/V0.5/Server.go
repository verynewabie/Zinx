package main

import (
	"Zinx/ziface"
	"Zinx/znet"
	"fmt"
)

// PingRouter 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle 自定义处理方法
func (router *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("receive from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回写数据
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//创建一个server句柄
	s := znet.NewServer()

	//配置路由
	s.AddRouter(&PingRouter{})

	//开启服务
	s.Serve()
}
