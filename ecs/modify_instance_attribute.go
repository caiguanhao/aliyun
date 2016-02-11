package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ModifyInstanceAttribute struct {
	RequestId string `json:"RequestId"`
}

var UPDATE_INSTANCE cli.Command = cli.Command{
	Name:      "update-instance",
	Aliases:   []string{"update", "e"},
	Usage:     "update attributes of an instance",
	ArgsUsage: "[instance IDs...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Usage: "new name of the instance",
		},
		cli.StringFlag{
			Name:  "description, d",
			Usage: "new description of the instance",
		},
	},
	Action: func(c *cli.Context) {
		params := map[string]string{}
		if c.IsSet("name") {
			params["InstanceName"] = c.String("name")
		}
		if c.IsSet("description") {
			params["Description"] = c.String("description")
		}
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.ModifyInstanceAttributeById(arg, params))
		})
	},
	BashComplete: DescribeInstancesForBashComplete,
}

var HIDE_INSTANCE cli.Command = cli.Command{
	Name:      "hide-instance",
	Aliases:   []string{"hide", "h"},
	Usage:     "hide instance from instance list",
	ArgsUsage: "[instance IDs...]",
	Action: func(c *cli.Context) {
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.HideInstanceById(arg, true))
		})
	},
	BashComplete: DescribeInstancesForBashComplete,
}

var UNHIDE_INSTANCE cli.Command = cli.Command{
	Name:      "unhide-instance",
	Aliases:   []string{"unhide", "u"},
	Usage:     "un-hide instance from instance list",
	ArgsUsage: "[instance IDs...]",
	Action: func(c *cli.Context) {
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.HideInstanceById(arg, false))
		})
	},
	BashComplete: DescribeInstancesForBashComplete,
}

func (ecs *ECS) ModifyInstanceAttributeById(id string, _params map[string]string) (modify ModifyInstanceAttribute, _ error) {
	params := map[string]string{
		"Action":     "ModifyInstanceAttribute",
		"InstanceId": id,
	}
	for k, v := range _params {
		params[k] = v
	}
	if len(params) > 2 {
		return modify, ecs.Request(params, &modify)
	}
	return modify, errors.New("Please provide at least one: --name, --description.")
}

func (ecs *ECS) HideInstanceById(id string, hide bool) (modify ModifyInstanceAttribute, _ error) {
	instance, err := ecs.DescribeInstanceAttributeById(id)
	if err != nil {
		return modify, err
	}
	description := strings.Replace(instance.Description, "[HIDE]", "", -1)
	if hide {
		description = "[HIDE] " + description
	}
	description = strings.TrimSpace(description)
	return ecs.ModifyInstanceAttributeById(id, map[string]string{
		"Description": description,
	})
}

func (modify ModifyInstanceAttribute) Print() {
	fmt.Println(modify.RequestId)
}

func (modify ModifyInstanceAttribute) PrintTable() {
	fmt.Println(modify.RequestId)
}
