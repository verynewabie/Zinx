package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"fmt"
	"net"
)

type Server struct {
	Name string
	//tcp4 or other
	IPVersion  string
	IP         string
	Port       int
	msgHandler ziface.IMsgHandler
	//当前Server的链接管理器
	ConnMgr ziface.IConnectionManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection) //可以看做类内函数，该函数不和该结构体绑定
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("Add router success! msgId = ", msgId)
}

func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	fmt.Println("[Start] Server listener at IP:", s.IP, "Port:", s.Port)
	//开启一个go去做服务端Listener业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()
		//1.获取Tcp的address
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err.Error())
			return
		}
		//2.监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "error:", err.Error())
			return
		}
		fmt.Println("start zinx server", s.Name, "success, listening...")
		//3.阻塞等待客户端连接，处理业务
		var connectionID uint32 = 0
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			//3.2 如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			//3.3 处理该新连接请求的业务方法
			dealConn := NewConnection(s, conn, connectionID, s.msgHandler)
			connectionID++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	s.ConnMgr.ClearConnection()
}
func (s *Server) Serve() {
	s.Start()
	// TODO 服务器启动后的一些额外业务
	//这里要阻塞，不然客户端一调用Serve函数就结束了
	select {}
}
func NewServer() ziface.IServer {
	// s 使用tcp/ip版本4
	//utils.GlobalObject.Reload()
	s := &Server{
		Name:       utils.GlobalObject.Name, //从全局参数获取
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,    //从全局参数获取
		Port:       utils.GlobalObject.TcpPort, //从全局参数获取
		msgHandler: NewMsgHandler(),
		ConnMgr:    NewConnectionManager(), //创建ConnManager
	}
	return s
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnectionManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
