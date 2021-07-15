package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	// 模拟的服务器
	listenner, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}

	// 创建一个go承载处理业务
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error")
			}

			go func(conn net.Conn) {
				// 拆包
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error :", err)
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpacke err:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// 第二次读取
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
						}

						// 完整的信息读取完毕
						fmt.Println("----> Recv MsgID：", msg.Id, ", dataLen：", msg.DataLen, "data：", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial err：", err)
		return
	}

	// 封包
	dp := NewDataPack()

	// 模拟粘包过程 封装两个msg一起发

	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte("zinx"),
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack ms1 error:", err)
		return
	}

	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte("zinx2"),
	}

	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack ms2 error:", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	select {}
}
