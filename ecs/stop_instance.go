package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type StopInstance struct {
	RequestId string `json:"RequestId"`
}

var STOP_INSTANCE cli.Command = cli.Command{
	Name:    "stop-instance",
	Aliases: []string{"stop", "k"},
	Usage:   "stop an instance",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.StopInstanceById(c.Args().Get(0)))
	},
}

func (ecs *ECS) StopInstanceById(id string) (stop StopInstance, _ error) {
	return stop, ecs.Request(map[string]string{
		"Action":     "StopInstance",
		"InstanceId": id,
	}, &stop)
}

func (stop StopInstance) Print() {
	fmt.Println(stop.RequestId)
}

func (stop StopInstance) PrintTable() {
	fmt.Println(stop.RequestId)
}
