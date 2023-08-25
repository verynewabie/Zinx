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
	go func() {
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

			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			//3.3 处理该新连接请求的业务方法
			dealConn := NewConnection(conn, connectionID, s.msgHandler)
			connectionID++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	// TODO 释放或回收资源、连接等
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
	}
	return s
}
