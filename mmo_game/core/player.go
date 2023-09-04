package core

import (
	"Zinx/mmo_game/pb"
	"Zinx/ziface"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
)

type Player struct {
	Pid  int32              //玩家id
	Conn ziface.IConnection //与该玩家客户端的连接
	X    float32            //平面x坐标
	Y    float32            //高度
	Z    float32            //平面y坐标,为什么这么写？因为客户端这么写的
	V    float32            //角度
}

var PidGen int32 = 1  //用来生成玩家ID的计数器
var IdLock sync.Mutex //保护PidGen的互斥机制

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴偏移若干坐标
		Y:    0,
		Z:    float32(134 + rand.Intn(17)), //随机在134坐标点 基于Y轴偏移若干坐标
		V:    0,                            //角度为0，尚未实现
	}
	return p
}

func (player *Player) SendMsg(msgId uint32, data proto.Message) {
	fmt.Printf("before Marshal data = %+v\n", data)
	//将proto Message结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}
	fmt.Printf("after Marshal data = %+v\n", msg)
	if player.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	//调用Zinx框架的SendMsg发包
	if err := player.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}
	return
}

// SyncPid 告知客户端pid,同步已经生成的玩家ID给客户端
func (player *Player) SyncPid() {
	//组建proto数据
	data := &pb.SyncPid{
		Pid: player.Pid,
	}

	//发送数据给客户端
	player.SendMsg(1, data)
}

// BroadCastStartPosition 广播玩家自己的出生地点
func (player *Player) BroadCastStartPosition() {

	msg := &pb.BroadCast{
		Pid: player.Pid,
		Tp:  2, //TP2 代表广播坐标
		Data: &pb.BroadCast_P{ // oneof选中其中的P
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		},
	}

	player.SendMsg(200, msg)
}

// Talk 广播玩家聊天
func (player *Player) Talk(content string) {
	//1. 组建MsgId200 proto数据
	msg := &pb.BroadCast{
		Pid: player.Pid,
		Tp:  1, //TP 1 代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//2. 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//3. 向所有的玩家发送MsgId:200消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

// SyncSurrounding 给当前玩家周边的(九宫格内)玩家广播自己的位置，让他们显示自己
func (player *Player) SyncSurrounding() {
	//1 根据自己的位置，获取周围九宫格内的玩家pid
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(player.X, player.Z)
	//2 根据pid得到所有玩家对象
	players := make([]*Player, 0, len(pids))
	//3 给这些玩家发送MsgID:200消息，让自己出现在对方视野中
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	//3.1 组建MsgId200 proto数据
	msg := &pb.BroadCast{
		Pid: player.Pid,
		Tp:  2, //TP2 代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		},
	}
	//3.2 每个玩家分别给对应的客户端发送200消息，显示人物
	for _, player := range players {
		player.SendMsg(200, msg)
	}
	//4 让周围九宫格内的玩家出现在自己的视野中
	//4.1 制作Message SyncPlayers 数据
	playersData := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersData = append(playersData, p)
	}

	//4.2 封装SyncPlayer protobuf数据
	SyncPlayersMsg := &pb.SyncPlayers{
		Ps: playersData[:],
	}

	//4.3 给当前玩家发送需要显示周围的全部玩家数据
	player.SendMsg(202, SyncPlayersMsg)
}

// UpdatePos 广播玩家位置移动
func (player *Player) UpdatePos(x float32, y float32, z float32, v float32) {
	//更新玩家的位置信息
	player.X = x
	player.Y = y
	player.Z = z
	player.V = v

	//组装protobuf协议，发送位置给周围玩家
	msg := &pb.BroadCast{
		Pid: player.Pid,
		Tp:  4, //4 - 移动之后的坐标信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		},
	}

	//获取当前玩家周边全部玩家
	players := player.GetSurroundingPlayers()
	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

// GetSurroundingPlayers 获得当前玩家的AOI周边玩家信息
func (player *Player) GetSurroundingPlayers() []*Player {
	//得到当前AOI区域的所有pid
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(player.X, player.Z)

	//将所有pid对应的Player放到Player切片中
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

// LostConnection 玩家下线
func (player *Player) LostConnection() {
	//1 获取周围AOI九宫格内的玩家
	players := player.GetSurroundingPlayers()

	//2 封装MsgID:201消息
	msg := &pb.SyncPid{
		Pid: player.Pid,
	}

	//3 向周围玩家发送消息
	for _, player := range players {
		player.SendMsg(201, msg)
	}

	//4 世界管理器将当前玩家从AOI中摘除
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(player.Pid), player.X, player.Z)
	WorldMgrObj.RemovePlayerByPid(player.Pid)
}
