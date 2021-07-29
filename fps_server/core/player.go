package core

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
	"zinx/fps_server/pb"
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

// 告知客户端玩家pid 同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	// 组装proto数据
	data := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 发送给客户端

	p.SendMsg(1, data)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, data)
}

// 广播聊天信息
func (p *Player) Talk(content string) {
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, // 聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	p.SendMsg(200, proto_msg)
}
