package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type AllocatePublicIpAddress struct {
	RequestId string `json:"RequestId"`
	IpAddress string `json:"IpAddress"`
}

var ALLOCATE_PUBLIC_IP_ADDRESS cli.Command = cli.Command{
	Name:    "allocate-public-ip",
	Aliases: []string{"allocate", "a"},
	Usage:   "allocate an IP address for an instance",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.AllocatePublicIpAddressById(c.Args().Get(0)))
	},
}

func (ecs *ECS) AllocatePublicIpAddressById(id string) (alloc AllocatePublicIpAddress, _ error) {
	return alloc, ecs.Request(map[string]string{
		"Action":     "AllocatePublicIpAddress",
		"InstanceId": id,
	}, &alloc)
}

func (alloc AllocatePublicIpAddress) Print() {
	fmt.Println(alloc.IpAddress)
}

func (alloc AllocatePublicIpAddress) PrintTable() {
	fmt.Println(alloc.IpAddress)
}
