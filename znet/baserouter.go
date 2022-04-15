package znet

import "zinx/ziface"

// 实现 router 时，先嵌入 BaseRouter 基类，根据需求对该基类进行部分或全部重写
type BaseRouter struct{}

// 之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化

func (br *BaseRouter) PreHandle(r ziface.IRequest)  {}
func (br *BaseRouter) Handle(r ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(r ziface.IRequest) {}
