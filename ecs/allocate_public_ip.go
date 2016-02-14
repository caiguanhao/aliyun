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
	Name:      "allocate-public-ip",
	Aliases:   []string{"allocate", "a"},
	Usage:     "allocate an IP address for an instance",
	ArgsUsage: "[instance IDs...]",
	Action: func(c *cli.Context) {
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.AllocatePublicIpAddressById(arg))
		})
	},
	BashComplete: describeInstancesForBashComplete(func(instance ECSInstance) bool {
		return instance.PublicIpAddress.GetIPAddress(0) == ""
	}),
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
