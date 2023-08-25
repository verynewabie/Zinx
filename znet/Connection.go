package znet

import (
	"Zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	// ExitBuffChan 告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	// MsgHandler 链接的处理方法
	MsgHandler ziface.IMsgHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	connection := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
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
		// 创建拆包解包的对象
		dp := NewDataPack()
		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			conn.ExitBuffChan <- true
			continue
		}
		//拆包，得到msgId 和 dataLen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			conn.ExitBuffChan <- true
			continue
		}
		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(conn.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				conn.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)
		//得到当前客户端请求的Request数据
		request := Request{
			conn: conn,
			msg:  msg, //将之前的buf 改成 msg
		}
		//调用当前链接业务
		go conn.MsgHandler.DoMsgHandler(&request)
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

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (conn *Connection) SendMsg(msgId uint32, data []byte) error {
	if conn.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("pack error msg ")
	}

	//写回客户端
	if _, err := conn.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		conn.ExitBuffChan <- true
		return errors.New("conn Write error")
	}

	return nil
}
