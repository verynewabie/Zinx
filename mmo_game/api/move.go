package api

import (
	"Zinx/mmo_game/core"
	"Zinx/mmo_game/pb"
	"Zinx/ziface"
	"Zinx/znet"
	"fmt"
	"google.golang.org/protobuf/proto"
)

// MoveApi 玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (*MoveApi) Handle(request ziface.IRequest) {
	//1. 将客户端传来的proto协议解码
	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error ", err)
		return
	}

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pid中获取
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		request.GetConnection().Stop()
		return
	}

	fmt.Printf("user pid = %d , move(%f,%f,%f,%f)", pid, msg.X, msg.Y, msg.Z, msg.V)

	//3. 根据pid得到player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4. 让player对象发起移动位置信息广播
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
