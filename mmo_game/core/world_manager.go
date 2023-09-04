package core

import "sync"

type WorldManager struct {
	AoiMgr  *AOIManager       //当前世界地图的AOI规划管理器
	Players map[int32]*Player //当前在线的玩家集合
	pLock   sync.RWMutex      //保护Players的互斥读写机制
}

var WorldMgrObj *WorldManager

const (
	AOI_MIN_X int = 0
	AOI_MAX_X int = 250
	AOI_CNT_X int = 5
	AOI_MIN_Y int = 0
	AOI_MAX_Y int = 250
	AOI_CNT_Y int = 5
)

func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNT_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNT_Y),
	}
}

// AddPlayer 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) AddPlayer(player *Player) {
	//将player添加到 世界管理器中
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	//将player 添加到AOI网络规划中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

// RemovePlayerByPid 从玩家信息表中移除一个玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	wm.pLock.Lock()
	delete(wm.Players, pid)
	wm.pLock.Unlock()
}

// GetPlayerByPid 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

// GetAllPlayers 获取所有玩家的信息
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	//创建返回的player集合切片
	players := make([]*Player, 0)

	//添加切片
	for _, v := range wm.Players {
		players = append(players, v)
	}

	//返回
	return players
}
