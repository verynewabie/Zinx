package ziface

import "net"

type IConnection interface {
	// Start 启动连接
	Start()
	// Stop 停止连接
	Stop()
	// GetConnID 获取当前连接的ID
	GetConnID() uint32
	// GetTCPConnection 获取原生连接
	GetTCPConnection() *net.TCPConn
	// RemoteAddr 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// SendMsg 直接将Message数据发送数据给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error
}
