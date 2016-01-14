package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type StartInstance struct {
	RequestId string `json:"RequestId"`
}

func (start *StartInstance) Do(ecs *ECS) (*StartInstance, error) {
	if id, err := opts.GetInstanceId(); err == nil {
		return start, ecs.Request(map[string]string{
			"Action":     "StartInstance",
			"InstanceId": id,
		}, start)
	} else {
		return nil, err
	}
}

func (start StartInstance) Print() {
	fmt.Println(start.RequestId)
}

func (start StartInstance) PrintTable() {
	fmt.Println(start.RequestId)
}
