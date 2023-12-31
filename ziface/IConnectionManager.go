package ziface

type IConnectionManager interface {
	Add(conn IConnection)                   //添加链接
	Remove(conn IConnection)                //删除连接
	Get(connID uint32) (IConnection, error) //利用ConnID获取链接
	Len() int                               //获取当前连接个数
	ClearConnection()                       //删除并停止所有链接
}
