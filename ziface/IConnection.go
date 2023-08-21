package ziface

import "net"

type IConnection interface {
	// Start 启动连接
	Start()
	// Stop 停止连接
	Stop()
	// GetConnID 获取当前连接的ID
	GetConnID() uint32
}

// HandFunc 定义统一处理连接的接口 三个参数：socket原生链接 客户端请求数据 数据长度
type HandFunc func(*net.TCPConn, []byte, int) error
