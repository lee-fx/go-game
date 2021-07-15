package znet

import "zinx/ziface"

// 实现router时，先嵌入这个baseRouter基类，然后根据需要对这个基类进行重写
type BaseRouter struct {
}

// 处理conn业务之前的钩子方法hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// 处理Conn业务的主方法hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// 处理Conn业务之后的钩子方法hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
