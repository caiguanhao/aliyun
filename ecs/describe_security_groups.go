package main

import (
	"errors"
	"fmt"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type DescribeSecurityGroups struct {
	PageNumber     int64  `json:"PageNumber"`
	PageSize       int64  `json:"PageSize"`
	RegionId       string `json:"RegionId"`
	RequestId      string `json:"RequestId"`
	SecurityGroups struct {
		SecurityGroup []struct {
			Description     string `json:"Description"`
			SecurityGroupId string `json:"SecurityGroupId"`
		} `json:"SecurityGroup"`
	} `json:"SecurityGroups"`
	TotalCount int64 `json:"TotalCount"`
}

func (groups *DescribeSecurityGroups) Do(ecs *ECS) (*DescribeSecurityGroups, error) {
	if len(opts.Region) < 1 {
		if opts.IsQuiet {
			return groups, nil
		}
		return nil, errors.New("Please provide a --region.")
	}
	return groups, ecs.Request(map[string]string{
		"Action":   "DescribeSecurityGroups",
		"RegionId": opts.Region,
	}, groups)
}

func (groups DescribeSecurityGroups) Print() {
	for _, group := range groups.SecurityGroups.SecurityGroup {
		fmt.Println(group.SecurityGroupId)
	}
}

func (groups DescribeSecurityGroups) PrintTable() {
	idMaxLength := 2
	for _, group := range groups.SecurityGroups.SecurityGroup {
		idLength := len(group.SecurityGroupId)
		if idLength > idMaxLength {
			idMaxLength = idLength
		}
	}
	format := fmt.Sprintf("%%-%ds  %%s\n", idMaxLength)
	fmt.Printf(format, "ID", "Description")
	for _, group := range groups.SecurityGroups.SecurityGroup {
		fmt.Printf(
			format,
			group.SecurityGroupId,
			group.Description,
		)
	}
}
