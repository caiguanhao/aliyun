package main

import (
	"fmt"
	"os"

	"github.com/caiguanhao/aliyun/ecs/errors"
	"github.com/caiguanhao/aliyun/ecs/opts"
)

type CreateInstance struct {
	InstanceId string `json:"InstanceId"`
	RequestId  string `json:"RequestId"`
}

var password = os.Getenv("PASSWORD")

func validate() errors.Errors {
	var errs errors.Errors
	if len(password) < 1 {
		errs.Add("Please provide a password via environment variable PASSWORD.")
	}
	if len(opts.InstanceImage) < 1 {
		errs.Add("Please provide a --image.")
	}
	if len(opts.InstanceType) < 1 {
		errs.Add("Please provide a --type.")
	}
	if len(opts.InstanceGroup) < 1 {
		errs.Add("Please provide a --group.")
	}
	if len(opts.InstanceName) < 1 {
		errs.Add("Please provide a --name.")
	}
	if len(opts.Region) < 1 {
		errs.Add("Please provide a --region.")
	}
	return errs
}

func (create *CreateInstance) Do(ecs *ECS) (*CreateInstance, error) {
	if errs := validate(); errs.HaveError() {
		return nil, errs.Errorify()
	}
	return create, ecs.Request(map[string]string{
		"Action":                        "CreateInstance",
		"ImageId":                       opts.InstanceImage,
		"InstanceType":                  opts.InstanceType,
		"SecurityGroupId":               opts.InstanceGroup,
		"InstanceName":                  opts.InstanceName,
		"HostName":                      opts.InstanceName,
		"RegionId":                      opts.Region,
		"InternetChargeType":            "PayByTraffic",
		"InternetMaxBandwidthIn":        "5",
		"InternetMaxBandwidthOut":       "5",
		"Password":                      password,
		"SystemDisk.Category":           "cloud",
		"DataDisk.1.Size":               "10",
		"DataDisk.1.Category":           "cloud",
		"DataDisk.1.Device":             "/dev/xvdb",
		"DataDisk.1.DeleteWithInstance": "true",
	}, create)
}

func (create CreateInstance) Print() {
	fmt.Println(create.InstanceId)
}

func (create CreateInstance) PrintTable() {
	fmt.Println(create.InstanceId)
}
