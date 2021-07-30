package core

import (
	"fmt"
	"sync"
)

/*
	一个AOI地图中的格子类型
*/

type Grid struct {
	GID       int          // id
	MinX      int          // 左边界坐标
	MaxX      int          // 右边界坐标
	MinY      int          // 上边界坐标
	MaxY      int          // 下边界坐标
	playerIDs map[int]bool // 格子内玩家集合
	pIDLock   sync.RWMutex
}

// 初始化格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 添加玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID] = true
}

// 删除玩家
func (g *Grid) Delete(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 获取格子所有玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

// 调试信息
func (g *Grid) string() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}