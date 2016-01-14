package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type ModifyInstanceAttribute struct {
	RequestId string `json:"RequestId"`
}

func (modify ModifyInstanceAttribute) HideInstance(ecs *ECS, hide bool) (*ModifyInstanceAttribute, error) {
	id, err := opts.GetInstanceId()
	if err != nil {
		return nil, err
	}
	var instance DescribeInstanceAttribute
	err = instance.DescribeInstanceAttributeById(ecs, id)
	if err != nil {
		return nil, err
	}
	description := strings.Replace(instance.Description, "[HIDE]", "", -1)
	if hide {
		description = "[HIDE] " + description
	}
	description = strings.TrimSpace(description)
	return &modify, ecs.Request(map[string]string{
		"Action":      "ModifyInstanceAttribute",
		"InstanceId":  id,
		"Description": description,
	}, &modify)
}

func (modify *ModifyInstanceAttribute) Do(ecs *ECS) (*ModifyInstanceAttribute, error) {
	id, err := opts.GetInstanceId()
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"Action":     "ModifyInstanceAttribute",
		"InstanceId": id,
	}
	if opts.Description != "\x00" {
		params["Description"] = opts.Description
	}
	if opts.InstanceName != "" {
		params["InstanceName"] = opts.InstanceName
	}
	if len(params) > 2 {
		return modify, ecs.Request(params, modify)
	}
	return nil, errors.New("Please provide at least one: --name, --description.")
}

func (modify ModifyInstanceAttribute) Print() {
	fmt.Println(modify.RequestId)
}

func (modify ModifyInstanceAttribute) PrintTable() {
	fmt.Println(modify.RequestId)
}
