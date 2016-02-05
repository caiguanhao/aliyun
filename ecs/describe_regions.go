package main

import (
	"fmt"
	"sort"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ECSRegion struct {
	RegionID string `json:"RegionId"`
}

var DESCRIBE_REGIONS cli.Command = cli.Command{
	Name:    "list-regions",
	Aliases: []string{"regions", "n"},
	Usage:   "list all available regions",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.DescribeRegions())
	},
}

type DescribeRegions struct {
	Regions struct {
		Region ECSRegions `json:"Region"`
	} `json:"Regions"`
	RequestID string `json:"RequestId"`
}

type ECSRegions []ECSRegion

func (a ECSRegions) Len() int           { return len(a) }
func (a ECSRegions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ECSRegions) Less(i, j int) bool { return a[i].RegionID < a[j].RegionID }

func (ecs *ECS) DescribeRegions() (_ ECSRegions, resp DescribeRegions, _ error) {
	return resp.Regions.Region, resp, ecs.Request(map[string]string{
		"Action": "DescribeRegions",
	}, &resp)
}

func (regions ECSRegions) Print() {
	sort.Sort(regions)
	for _, region := range regions {
		fmt.Println(region.RegionID)
	}
}

func (regions ECSRegions) PrintTable() {
	regions.Print()
}
