package ziface

type IRouter interface {
	PreHandle(r IRequest)  // 处理业务之前的 hook 方法
	Handle(r IRequest)     // 处理业务的主方法
	PostHandle(r IRequest) // 处理业务之后的 hook 方法
}
