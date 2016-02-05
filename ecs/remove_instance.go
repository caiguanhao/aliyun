package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type RemoveInstance struct {
	RequestId string `json:"RequestId"`
}

var REMOVE_INSTANCE cli.Command = cli.Command{
	Name:    "remove-instance",
	Aliases: []string{"remove", "d"},
	Usage:   "remove an instance",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.RemoveInstanceById(c.Args().Get(0)))
	},
}

func (ecs *ECS) RemoveInstanceById(id string) (remove RemoveInstance, _ error) {
	return remove, ecs.Request(map[string]string{
		"Action":     "DeleteInstance",
		"InstanceId": id,
	}, &remove)
}

func (remove RemoveInstance) Print() {
	fmt.Println(remove.RequestId)
}

func (remove RemoveInstance) PrintTable() {
	fmt.Println(remove.RequestId)
}
