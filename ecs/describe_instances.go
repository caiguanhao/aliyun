package main

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type DescribeInstances struct {
	Instances struct {
		Instance []DescribeInstanceAttribute `json:"Instance"`
	} `json:"Instances"`
	PageNumber int64  `json:"PageNumber"`
	PageSize   int64  `json:"PageSize"`
	RequestId  string `json:"RequestId"`
	TotalCount int64  `json:"TotalCount"`
}

func (instances *DescribeInstances) Do(ecs *ECS) (*DescribeInstances, error) {
	if len(opts.Region) < 1 {
		if opts.IsQuiet {
			return instances, nil
		}
		return nil, errors.New("Please provide a --region.")
	}
	return instances, ecs.Request(map[string]string{
		"Action":   "DescribeInstances",
		"RegionId": opts.Region,
	}, instances)
}

func (instances DescribeInstances) Print() {
	for _, instance := range instances.Instances.Instance {
		if instance.ShouldHide() {
			continue
		}
		if opts.PrintNameAndId {
			fmt.Printf("%s:%s\n", instance.InstanceName, instance.InstanceId)
		} else if opts.PrintName {
			fmt.Println(instance.InstanceName)
		} else {
			fmt.Println(instance.InstanceId)
		}
	}
}

func (instances DescribeInstances) PrintTable() {
	idMaxLength := 2
	nameMaxLength := 4
	statusMaxLength := 6
	for _, instance := range instances.Instances.Instance {
		if instance.ShouldHide() {
			continue
		}
		idLength := len(instance.InstanceId)
		nameLength := len(instance.InstanceName)
		statusLength := len(instance.Status)
		if idLength > idMaxLength {
			idMaxLength = idLength
		}
		if nameLength > nameMaxLength {
			nameMaxLength = nameLength
		}
		if statusLength > statusMaxLength {
			statusMaxLength = statusLength
		}
	}
	format := fmt.Sprintf(
		"%%-%ds  %%-%ds  %%-%ds  %%-15s  %%-15s  %%-15s  %%-35s\n",
		idMaxLength,
		nameMaxLength,
		statusMaxLength,
	)
	fmt.Printf(format, "ID", "Name", "Status", "Public IP", "Private IP", "Type", "Created At")
	for _, instance := range instances.Instances.Instance {
		if instance.ShouldHide() {
			continue
		}
		createdAt, _ := time.Parse("2006-01-02T15:04Z", instance.CreationTime)
		duration := time.Since(createdAt)
		createdAtStr := fmt.Sprintf("%s (%.0f days ago)",
			createdAt.Local().Format("2006-01-02 15:04:05"),
			math.Floor(duration.Hours()/24))
		fmt.Printf(
			format,
			instance.InstanceId,
			instance.InstanceName,
			instance.Status,
			instance.PublicIpAddress.GetIPAddress(0),
			instance.InnerIpAddress.GetIPAddress(0),
			instance.InstanceType,
			createdAtStr,
		)
	}
}
