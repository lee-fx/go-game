package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	// 初始化
	aoiMgr := NewAOIManager(100, 300, 4, 200, 450, 5)
	// 打印
	fmt.Println(aoiMgr)
}
