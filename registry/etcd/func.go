package etcd

import (
	"fmt"
	"gatesvr/registry"
)

// 构建实例ID
func makeInsID(ins *registry.ServiceInstance) string {
	return fmt.Sprintf("%s-%s", ins.Kind, ins.ID)
}
