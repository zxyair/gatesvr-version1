package client

import (
	"gatesvr/cluster"
)

type Options struct {
	Addr         string       // 连接地址
	InsID        string       // 实例ID
	InsKind      cluster.Kind // 实例类型
	CloseHandler func()       // 关闭处理器
}
