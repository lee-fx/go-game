package ziface

type IRequest interface {
	// 得到当前连接
	GetConn() IConn
	// 得到请求的数据
	GetData() []byte
	GetMsgID() uint32
}
