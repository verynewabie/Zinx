package znet

import (
	"Zinx/ziface"
	"fmt"
	"strconv"
)

type MsgHandler struct {
	MsgIDToHandler map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法的map属性
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		MsgIDToHandler: make(map[uint32]ziface.IRouter),
	}
}
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.MsgIDToHandler[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.MsgIDToHandler[msgId]; ok {

		panic("repeated api , msgId = " + strconv.FormatUint(uint64(msgId), 10))
	}
	//2 添加msg与api的绑定关系
	mh.MsgIDToHandler[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}