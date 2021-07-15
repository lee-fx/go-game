package ziface

// TLV格式封装与拆包 面向tcp链接的数据流，处理tcp粘包问题
// stream format [dataLen][id][data] <-> [dataLen][id][data] <-> [dataLen][id][data ]
type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
