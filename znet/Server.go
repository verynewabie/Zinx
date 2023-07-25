package znet

import (
	"Zinx/ziface"
	"fmt"
	"net"
)

type Serve struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func (s *Serve) Start() {
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
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err.Error())
				continue
			}
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf) // cnt代表读取的长度
					if err != nil {
						fmt.Println("Receive error:", err.Error())
						continue
					}
					if _, err := conn.Write(buf[:cnt]); err != nil { //我们不关心写的长度
						fmt.Println("Write back error:", err.Error())
						continue
					}
				}
			}()
		}
	}()
}
func (s *Serve) Stop() {
	//TODO 释放或回收资源、连接等
}
func (s *Serve) Serve() {
	s.Start()
	//TODO 服务器启动后的一些额外业务
	//这里要阻塞，不然客户端一调用Serve函数就结束了
	select {}
}
func NewServer(name string) ziface.IServer {
	// s 使用tcp/ip版本4
	s := &Serve{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
