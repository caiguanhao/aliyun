package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type StopInstance struct {
	RequestId string `json:"RequestId"`
}

func (stop *StopInstance) Do(ecs *ECS) (*StopInstance, error) {
	if id, err := opts.GetInstanceId(); err == nil {
		return stop, ecs.Request(map[string]string{
			"Action":     "StopInstance",
			"InstanceId": id,
		}, stop)
	} else {
		return nil, err
	}
}

func (stop StopInstance) Print() {
	fmt.Println(stop.RequestId)
}

func (stop StopInstance) PrintTable() {
	fmt.Println(stop.RequestId)
}
