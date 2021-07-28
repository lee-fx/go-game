package core

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
	"zinx/ziface"
)

// 玩家对象
type Player struct {
	Pid  int32
	Conn ziface.IConn // 当前玩家的链接（用于和客户端的链接）
	X    float32      // xyz 和角度v
	Y    float32
	Z    float32
	V    float32
}

// Player id 生成器
var PidGen int32 = 1

var IdLock sync.Mutex

// 创建玩家
func NewPlayer(conn ziface.IConn) *Player {
	IdLock.Lock()
	PidGen++
	IdLock.Unlock()
	return &Player{
		Pid:  PidGen,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // 随机生成
		Y:    9,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}
}

// 给客户端消息的方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("proto marshal err ：", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("conn in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player Send Msg err :", err)
		return
	}
	return
}
