package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	// ExitBuffChan 告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	// MsgHandler 链接的处理方法
	MsgHandler ziface.IMsgHandler
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	connection := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}
	connection.TcpServer.GetConnMgr().Add(connection)
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
			return
		}
		//拆包，得到msgId 和 dataLen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			conn.ExitBuffChan <- true
			return
		}
		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(conn.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				conn.ExitBuffChan <- true
				return
			}
		}
		msg.SetData(data)
		//得到当前客户端请求的Request数据
		request := Request{
			conn: conn,
			msg:  msg, //将之前的buf 改成 msg
		}
		//根据用户配置决定是否使用工作池机制
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			conn.MsgHandler.SendMsgToTaskQueue(&request)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go conn.MsgHandler.DoMsgHandler(&request)
		}
	}
}

// Start 启动连接
func (conn *Connection) Start() {
	//必须要先开启读写再执行钩子函数，因为钩子函数内可能有网络通信
	go conn.StartReader()
	go conn.StartWriter()
	conn.TcpServer.CallOnConnStart(conn)
	select {
	case <-conn.ExitBuffChan:
		//得到退出消息，不再阻塞
		return
	}
}

// Stop 终止连接
func (conn *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", conn.ConnID)
	if conn.isClosed == true {
		return
	}
	conn.isClosed = true

	//如果用户注册了该链接的关闭回调业务，在此刻调用
	conn.TcpServer.CallOnConnStop(conn)
	// 关闭socket链接
	conn.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	conn.ExitBuffChan <- true
	//将链接从连接管理器中删除
	conn.TcpServer.GetConnMgr().Remove(conn)
	//关闭该链接全部管道
	close(conn.ExitBuffChan)
	close(conn.msgBuffChan)
	close(conn.msgChan)
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
	conn.msgChan <- msg //发送给Channel 供Writer读取

	return nil
}

func (conn *Connection) StartWriter() {

	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(conn.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-conn.msgChan:
			//有数据要写给客户端
			if _, err := conn.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case data, ok := <-conn.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := conn.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				break
				fmt.Println("msgBuffChan is Closed")
			}
		case <-conn.ExitBuffChan:
			//conn已经关闭
			return
		}
	}
}
func (conn *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if conn.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	conn.msgBuffChan <- msg

	return nil
}
