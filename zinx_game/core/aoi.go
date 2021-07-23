package core

import "fmt"

/**
aoi 管理模块
*/

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
