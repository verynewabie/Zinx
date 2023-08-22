package znet

import "Zinx/ziface"

type Request struct {
	conn ziface.IConnection
	data []byte
}

// GetConnection 获取请求连接信息
func (request *Request) GetConnection() ziface.IConnection {
	return request.conn
}

// GetData 获取请求消息的数据
func (request *Request) GetData() []byte {
	return request.data
}
