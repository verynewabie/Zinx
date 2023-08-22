package znet

import (
	"Zinx/ziface"
	"fmt"
	"net"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	// ExitBuffChan 告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	// Router 该链接的处理方法
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	connection := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}
	return connection
}

// StartReader 处理conn读数据的协程
func (conn *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(conn.RemoteAddr().String(), " conn reader exit!")
	defer conn.Stop()

	for {
		//读取我们的数据到buf中
		buf := make([]byte, 512)
		_, err := conn.Conn.Read(buf)
		if err != nil {
			fmt.Println("receive buf error ", err)
			conn.ExitBuffChan <- true
			continue
		}
		request := Request{
			conn: conn,
			data: buf,
		}
		//调用当前链接业务
		go func(request ziface.IRequest) {
			//执行注册的路由方法
			conn.Router.PreHandle(request)
			conn.Router.Handle(request)
			conn.Router.PostHandle(request)
		}(&request)
	}
}

// Start 启动连接
func (conn *Connection) Start() {
	go conn.StartReader()

	select {
	case <-conn.ExitBuffChan:
		//得到退出消息，不再阻塞
		return
	}
}

// Stop 终止连接
func (conn *Connection) Stop() {
	if conn.isClosed == true {
		return
	}
	conn.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	// 关闭socket链接
	conn.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	conn.ExitBuffChan <- true

	//关闭该链接全部管道
	close(conn.ExitBuffChan)
}

// GetTCPConnection 获取原始的socket TCPConn
func (conn *Connection) GetTCPConnection() *net.TCPConn {
	return conn.Conn
}

// GetConnID 获取当前连接ID
func (conn *Connection) GetConnID() uint32 {
	return conn.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (conn *Connection) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}
