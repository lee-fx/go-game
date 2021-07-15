package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConn
	msg  ziface.IMessage
}

func (r *Request) GetConn() ziface.IConn {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
