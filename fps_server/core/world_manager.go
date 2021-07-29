package core

import "sync"

/*
	游戏世界总管理模块
 */

type WorldManager struct {
	AoiMgr *AOIManager

	Players map[int32]*Player
	pLock sync.RWMutex
}

// 对外全局世界管理句柄
var WorldMgrObj *WorldManager

// 当导入core包的时候调用执行
func init() {
	WorldMgrObj = &WorldManager{
		// 创建世界AOI地图
		AoiMgr: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}

// 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	// 全局添加
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	wm.Players[player.Pid] = player

	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Y)
}

// 删除一个玩家
func (wm *WorldManager) RemovePlayer(pid int32) {
	// 格子中删除
	player := wm.Players[pid]
	wm.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Y)

	// 全局管理删除
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	delete(wm.Players, pid)
}

// 通过玩家查询player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pid]
}

// 获取全部在线玩家对象
func (wm *WorldManager) GetAllPlayer() (players []*Player) {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	players = make([]*Player, 0)
	for _, v := range wm.Players {
		players = append(players, v)
	}
	return
}