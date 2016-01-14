package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type AllocatePublicIP struct {
	RequestId string `json:"RequestId"`
	IpAddress string `json:"IpAddress"`
}

func (alloc *AllocatePublicIP) Do(ecs *ECS) (*AllocatePublicIP, error) {
	if id, err := opts.GetInstanceId(); err == nil {
		return alloc, ecs.Request(map[string]string{
			"Action":     "AllocatePublicIpAddress",
			"InstanceId": id,
		}, alloc)
	} else {
		return nil, err
	}
}

func (alloc AllocatePublicIP) Print() {
	fmt.Println(alloc.IpAddress)
}

func (alloc AllocatePublicIP) PrintTable() {
	fmt.Println(alloc.IpAddress)
}
