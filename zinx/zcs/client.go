package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

// 模拟客户端
func main() {
	fmt.Println("client ping start...")
	time.Sleep(time.Second * 1)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client err :", err)
		return
	}
	for {
		// 发送封包消息包
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinx client ping test message")))
		if err != nil {
			fmt.Println("Pack error: ", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Write error: ", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
			break
		}

		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unPack error: ", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// 第二次数据读取
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.DataLen)
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read MsgData error: ", err)
				break
			}
			fmt.Println("-----> Recv Server Msg: ID = ", msg.GetMsgId(), ", len = ", msg.GetMsgLen(), ", data = ", string(msg.GetMsgData()))
		}

		// cpu 阻塞
			time.Sleep(time.Second * 1)

	}

}
