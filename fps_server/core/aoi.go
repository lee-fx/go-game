package core

import "fmt"

/**
aoi 管理模块
*/

// 地图宏
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

type AOIManager struct {
	// 左边界坐标
	MinX int
	// 右边界坐标
	MaxX int
	// x方向格子数量
	CntsX int
	// 上边界坐标
	MinY int
	// 下边界坐标
	MaxY int
	// Y方向格子数量
	CntsY int
	// 当前区域有哪些格子id 对应value(格子对象)
	grids map[int]*Grid
}

// 初始化一个aoi管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 给aoi初始化区域的格子进行编号和初始化工作
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子的id (id = idY*cntsX + idX)
			gid := y*cntsX + x
			aoiMgr.grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridHeight(),
			)
		}
	}

	return aoiMgr
}

// 获取格子x轴的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 获取格子y轴的高度
func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印格子信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX: %d, MaxX: %d, CntsX: %d, MinY: %d, MaxY:%d, CntsY: %d", m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据格子gid获取周边格子集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断gid是否在aoi中
	if _, ok := m.grids[gID]; !ok {
		return
	}

	// 初始化grids返回值切片
	grids = append(grids, m.grids[gID])

	// 根据gID, 得到格子所在的坐标
	x, y := gID%m.CntsX, gID/m.CntsX

	// 新建一个临时存储周围格子的数组
	surroundGID := make([]int, 0)

	// 新建8个方向向量: 左上: (-1, -1), 左中: (-1, 0), 左下: (-1,1), 中上: (0,-1), 中下: (0,1), 右上:(1, -1)
	// 右中: (1, 0), 右下: (1, 1), 分别将这8个方向的方向向量按顺序写入x, y的分量数组
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	// 根据8个方向向量, 得到周围点的相对坐标, 挑选出没有越界的坐标, 将坐标转换为gID
	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]

		if newX >= 0 && newX < m.CntsX && newY >= 0 && newY < m.CntsY {
			surroundGID = append(surroundGID, newY*m.CntsX+newX)
		}
	}

	// 根据没有越界的gID, 得到格子信息
	for _, gID := range surroundGID {
		grids = append(grids, m.grids[gID])
	}

	return
}

// 通过坐标得到周边九宫格内全部layerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家的格子id
	gid := m.GetGidByPos(x, y)

	// 通过gid得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gid)

	// 将九宫格信息的全部玩家id加入到playerIDs
	for _, v := range grids {
		// 数组拼接
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}
	return
}

// 通过横纵坐标得到当前所在格子id
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridHeight()
	return idy*m.CntsX + idx
}

// 添加playerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 移除一个格子中的playerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Delete(pID)
}

// 通过gid获取全部playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	return m.grids[gID].GetPlayerIDs()
}

// 通过坐标将player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

// 通过一个坐标把一个player从一个格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Delete(pID)
}
