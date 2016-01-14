package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type RemoveInstance struct {
	RequestId string `json:"RequestId"`
}

func (remove *RemoveInstance) Do(ecs *ECS) (*RemoveInstance, error) {
	if id, err := opts.GetInstanceId(); err == nil {
		return remove, ecs.Request(map[string]string{
			"Action":     "DeleteInstance",
			"InstanceId": id,
		}, remove)
	} else {
		return nil, err
	}
}

func (remove RemoveInstance) Print() {
	fmt.Println(remove.RequestId)
}

func (remove RemoveInstance) PrintTable() {
	fmt.Println(remove.RequestId)
}
