package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	APIs           map[uint32]ziface.IRouter // 每一个 MsgID 对应的 Router
	TaskQueue      []chan ziface.IRequest    // worker 取任务的消息队列
	WorkerPoolSize uint32                    // WorkerPool 中 worker 的数量
}

func newMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.Global.WorkerPoolSize),
		WorkerPoolSize: utils.Global.WorkerPoolSize,
	}
}

func (mh *MsgHandler) DoMsgHandler(req ziface.IRequest) error {
	// 根据 MsgId 匹配 handler
	handler, ok := mh.APIs[req.GetMsgID()]
	if !ok {
		return fmt.Errorf("api of msgID=%d is not found", req.GetMsgID())
	}
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
	return nil
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) error {
	if mh.APIs[msgID] != nil {
		return fmt.Errorf("api of msgID=%d is already exists", msgID)
	}
	mh.APIs[msgID] = router
	fmt.Println("Add api MsgID =", msgID, "success.")
	return nil
}

// 启动 WorkerPool，只执行一次，外部可见
func (mh *MsgHandler) StartWorkerPool() {
	// 根据 WorkerPoolSize 分别开启 worker，每个对应一个 goroutine
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 为 worker 分配消息队列
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.Global.MaxWorkerTask)
		// 启动 worker
		go mh.startWorker(i)
	}
}

// 启动一个 worker
func (mh *MsgHandler) startWorker(id int) {
	fmt.Println("worker id=", id, "is started...")
	for req := range mh.TaskQueue[id] {
		mh.DoMsgHandler(req)
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	id := request.GetMsgID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		" request msgID=", request.GetMsgID(), "to workerID=", id)
	mh.TaskQueue[id] <- request
}
