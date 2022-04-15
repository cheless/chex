package ziface

type IMsgHandler interface {
	DoMsgHandler(req IRequest) error // 处理客户端的业务，并将结果通过 channel 发送给 Writer
	AddRouter(msgID uint32, router IRouter) error
	StartWorkerPool()
	SendMsgToTaskQueue(request IRequest)
}
