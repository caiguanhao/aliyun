package main

import (
	"fmt"
	"sort"
)

type ECSInstanceType struct {
	CpuCoreCount   int64   `json:"CpuCoreCount"`
	InstanceTypeId string  `json:"InstanceTypeId"`
	MemorySize     float64 `json:"MemorySize"`
}

type DescribeInstanceTypes struct {
	InstanceTypes struct {
		InstanceType []ECSInstanceType `json:"InstanceType"`
	} `json:"InstanceTypes"`
	RequestId string `json:"RequestId"`
}

type byCPUCoreThenMemorySize []ECSInstanceType

func (a byCPUCoreThenMemorySize) Len() int      { return len(a) }
func (a byCPUCoreThenMemorySize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCPUCoreThenMemorySize) Less(i, j int) bool {
	if a[i].CpuCoreCount < a[j].CpuCoreCount {
		return true
	} else if a[i].CpuCoreCount > a[j].CpuCoreCount {
		return false
	} else if a[i].MemorySize < a[j].MemorySize {
		return true
	}
	return false
}

func (types *DescribeInstanceTypes) Do(ecs *ECS) (*DescribeInstanceTypes, error) {
	return types, ecs.Request(map[string]string{
		"Action": "DescribeInstanceTypes",
	}, types)
}

func (types DescribeInstanceTypes) Print() {
	for _, itype := range types.InstanceTypes.InstanceType {
		fmt.Println(itype.InstanceTypeId)
	}
}

func (types DescribeInstanceTypes) PrintTable() {
	sort.Sort(byCPUCoreThenMemorySize(types.InstanceTypes.InstanceType))
	idMaxLength := 4
	for _, itype := range types.InstanceTypes.InstanceType {
		idLength := len(itype.InstanceTypeId)
		if idLength > idMaxLength {
			idMaxLength = idLength
		}
	}
	format := fmt.Sprintf("%%-%ds  %%-8s  %%-6s\n", idMaxLength)
	fmt.Printf(format, "Name", "CPU Core", "Memory")
	for _, itype := range types.InstanceTypes.InstanceType {
		fmt.Printf(
			format,
			itype.InstanceTypeId,
			fmt.Sprintf("%d", itype.CpuCoreCount),
			fmt.Sprintf("%.1f G", itype.MemorySize),
		)
	}
}
