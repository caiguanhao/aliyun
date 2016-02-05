package main

import (
	"fmt"
	"sort"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ECSInstanceType struct {
	CpuCoreCount   int64   `json:"CpuCoreCount"`
	InstanceTypeId string  `json:"InstanceTypeId"`
	MemorySize     float64 `json:"MemorySize"`
}

var DESCRIBE_INSTANCE_TYPES cli.Command = cli.Command{
	Name:    "list-instance-types",
	Aliases: []string{"types", "t"},
	Usage:   "list all instance types",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.DescribeInstanceTypes())
	},
}

type DescribeInstanceTypes struct {
	InstanceTypes struct {
		InstanceType ECSInstanceTypes `json:"InstanceType"`
	} `json:"InstanceTypes"`
	RequestId string `json:"RequestId"`
}

type ECSInstanceTypes []ECSInstanceType

func (a ECSInstanceTypes) Len() int      { return len(a) }
func (a ECSInstanceTypes) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ECSInstanceTypes) Less(i, j int) bool {
	if a[i].CpuCoreCount < a[j].CpuCoreCount {
		return true
	} else if a[i].CpuCoreCount > a[j].CpuCoreCount {
		return false
	} else if a[i].MemorySize < a[j].MemorySize {
		return true
	}
	return false
}

func (ecs *ECS) DescribeInstanceTypes() (_ ECSInstanceTypes, resp DescribeInstanceTypes, _ error) {
	return resp.InstanceTypes.InstanceType, resp, ecs.Request(map[string]string{
		"Action": "DescribeInstanceTypes",
	}, &resp)
}

func (types ECSInstanceTypes) Print() {
	for _, itype := range types {
		fmt.Println(itype.InstanceTypeId)
	}
}

func (types ECSInstanceTypes) PrintTable() {
	sort.Sort(types)
	fields := []interface{}{"Name", "CPU Core", "Memory"}
	PrintTable(fields, len(types), func(i int) []interface{} {
		itype := types[i]
		return []interface{}{
			itype.InstanceTypeId,
			fmt.Sprintf("%d", itype.CpuCoreCount),
			fmt.Sprintf("%.1f G", itype.MemorySize),
		}
	})
}
