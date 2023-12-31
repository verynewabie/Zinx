package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID     int          //格子ID
	MinX    int          //格子左边界坐标
	MaxX    int          //格子右边界坐标
	MinY    int          //格子上边界坐标
	MaxY    int          //格子下边界坐标
	players map[int]bool //当前格子内的玩家或者物体成员ID
	lock    sync.RWMutex //playerIDs的保护map的锁
}

// NewGrid 初始化一个格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:     gID,
		MinX:    minX,
		MaxX:    maxX,
		MinY:    minY,
		MaxY:    maxY,
		players: make(map[int]bool),
	}
}

// Add 向当前格子中添加一个玩家
func (g *Grid) Add(playerID int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.players[playerID] = true
}

// Remove 从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	delete(g.players, playerID)
}

// GetPlayers 得到当前格子中所有的玩家
func (g *Grid) GetPlayers() (playerIDs []int) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	for k, _ := range g.players {
		playerIDs = append(playerIDs, k)
	}

	return
}

// String 打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.players)
}
