package utils

import (
	"Zinx/ziface"
	"encoding/json"
	"os"
)

type GlobalObj struct {
	/*
		Server配置
	*/
	TcpServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称

	/*
		Zinx配置
	*/
	Version          string //当前Zinx版本号
	MaxPacketSize    uint32 //数据包的大小限制
	MaxConn          int    //服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32
	/*
		config file path
	*/
	ConfFilePath string
}

var GlobalObject *GlobalObj

// Reload 读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() { // 导入包的时候触发
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.8",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     "conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
	}

	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
