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
}
