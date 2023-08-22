package znet

import "Zinx/ziface"

// BaseRouter 路由的基类
type BaseRouter struct{}

func (router *BaseRouter) PreHandle(request ziface.IRequest)  {}
func (router *BaseRouter) Handle(request ziface.IRequest)     {}
func (router *BaseRouter) PostHandle(request ziface.IRequest) {}
