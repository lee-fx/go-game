package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

// 封包与拆包模块
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// datalen uint32(4 byte) + id uint32(4 byte)
	return 8
}

// 封包
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	// 将 len id data 写入databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// 拆包
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个ioReader
	reader := bytes.NewReader(binaryData)
	msg := &Message{}

	if err := binary.Read(reader, binary.LittleEndian,  &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian,  &msg.Id); err != nil {
		return nil, err
	}

	// 判断datalen长度是否非法
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}
