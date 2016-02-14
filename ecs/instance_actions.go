package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ActionResponse struct {
	RequestId string `json:"RequestId"`
}

var REMOVE_INSTANCE cli.Command = cli.Command{
	Name:      "remove-instance",
	Aliases:   []string{"remove", "d"},
	Usage:     "remove an instance",
	ArgsUsage: "[instance IDs...]",
	Action: func(c *cli.Context) {
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.RemoveInstanceById(arg))
		})
	},
	BashComplete: describeInstancesForBashComplete,
}

var RESTART_INSTANCE cli.Command = cli.Command{
	Name:      "restart-instance",
	Aliases:   []string{"restart", "r"},
	Usage:     "restart an instance",
	ArgsUsage: "[instance IDs...]",
	Action: func(c *cli.Context) {
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.RestartInstanceById(arg))
		})
	},
	BashComplete: describeInstancesForBashComplete,
}

var START_INSTANCE cli.Command = cli.Command{
	Name:      "start-instance",
	Aliases:   []string{"start", "s"},
	Usage:     "start an instance",
	ArgsUsage: "[instance IDs...]",
	Action: func(c *cli.Context) {
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.StartInstanceById(arg))
		})
	},
	BashComplete: describeInstancesForBashComplete,
}

var STOP_INSTANCE cli.Command = cli.Command{
	Name:      "stop-instance",
	Aliases:   []string{"stop", "k"},
	Usage:     "stop an instance",
	ArgsUsage: "[instance IDs...]",
	Action: func(c *cli.Context) {
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.StopInstanceById(arg))
		})
	},
	BashComplete: describeInstancesForBashComplete,
}

func executeInstanceActionById(ecs *ECS, action, id string) (resp ActionResponse, _ error) {
	return resp, ecs.Request(map[string]string{
		"Action":     action,
		"InstanceId": id,
	}, &resp)
}

func (ecs *ECS) RemoveInstanceById(id string) (ActionResponse, error) {
	return executeInstanceActionById(ecs, "DeleteInstance", id)
}

func (ecs *ECS) RestartInstanceById(id string) (ActionResponse, error) {
	return executeInstanceActionById(ecs, "RebootInstance", id)
}

func (ecs *ECS) StartInstanceById(id string) (ActionResponse, error) {
	return executeInstanceActionById(ecs, "StartInstance", id)
}

func (ecs *ECS) StopInstanceById(id string) (ActionResponse, error) {
	return executeInstanceActionById(ecs, "StopInstance", id)
}

func (resp ActionResponse) Print() {
	fmt.Println(resp.RequestId)
}

func (resp ActionResponse) PrintTable() {
	fmt.Println(resp.RequestId)
}
