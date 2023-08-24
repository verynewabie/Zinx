package znet

import "Zinx/ziface"

type Request struct {
	conn ziface.IConnection
	msg  ziface.IMessage
}

// GetConnection 获取请求连接信息
func (request *Request) GetConnection() ziface.IConnection {
	return request.conn
}

// GetData 获取请求消息的数据
func (request *Request) GetData() []byte {
	return request.msg.GetData()
}

// GetMsgID 获取请求的消息的ID
func (request *Request) GetMsgID() uint32 {
	return request.msg.GetMsgId()
}
