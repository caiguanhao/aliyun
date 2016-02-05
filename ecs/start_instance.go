package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type StartInstance struct {
	RequestId string `json:"RequestId"`
}

var START_INSTANCE cli.Command = cli.Command{
	Name:    "start-instance",
	Aliases: []string{"start", "s"},
	Usage:   "start an instance",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.StartInstanceById(c.Args().Get(0)))
	},
}

func (ecs *ECS) StartInstanceById(id string) (start StartInstance, _ error) {
	return start, ecs.Request(map[string]string{
		"Action":     "StartInstance",
		"InstanceId": id,
	}, &start)
}

func (start StartInstance) Print() {
	fmt.Println(start.RequestId)
}

func (start StartInstance) PrintTable() {
	fmt.Println(start.RequestId)
}
