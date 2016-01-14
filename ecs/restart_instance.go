package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type RestartInstance struct {
	RequestId string `json:"RequestId"`
}

func (restart *RestartInstance) Do(ecs *ECS) (*RestartInstance, error) {
	if id, err := opts.GetInstanceId(); err == nil {
		return restart, ecs.Request(map[string]string{
			"Action":     "RebootInstance",
			"InstanceId": id,
		}, restart)
	} else {
		return nil, err
	}
}

func (restart RestartInstance) Print() {
	fmt.Println(restart.RequestId)
}

func (restart RestartInstance) PrintTable() {
	fmt.Println(restart.RequestId)
}
