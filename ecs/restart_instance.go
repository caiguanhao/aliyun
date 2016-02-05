package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type RestartInstance struct {
	RequestId string `json:"RequestId"`
}

var RESTART_INSTANCE cli.Command = cli.Command{
	Name:    "restart-instance",
	Aliases: []string{"restart", "r"},
	Usage:   "restart an instance",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.RestartInstanceById(c.Args().Get(0)))
	},
}

func (ecs *ECS) RestartInstanceById(id string) (restart RestartInstance, _ error) {
	return restart, ecs.Request(map[string]string{
		"Action":     "RebootInstance",
		"InstanceId": id,
	}, &restart)
}

func (restart RestartInstance) Print() {
	fmt.Println(restart.RequestId)
}

func (restart RestartInstance) PrintTable() {
	fmt.Println(restart.RequestId)
}
