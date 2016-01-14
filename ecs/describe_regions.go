package main

import (
	"fmt"
	"sort"
)

type ECSRegion struct {
	RegionID string `json:"RegionId"`
}

type DescribeRegions struct {
	Regions struct {
		Region []ECSRegion `json:"Region"`
	} `json:"Regions"`
	RequestID string `json:"RequestId"`
}

type byRegionID []ECSRegion

func (a byRegionID) Len() int           { return len(a) }
func (a byRegionID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byRegionID) Less(i, j int) bool { return a[i].RegionID < a[j].RegionID }

func (regions *DescribeRegions) Do(ecs *ECS) (*DescribeRegions, error) {
	return regions, ecs.Request(map[string]string{
		"Action": "DescribeRegions",
	}, regions)
}

func (regions DescribeRegions) Print() {
	sort.Sort(byRegionID(regions.Regions.Region))
	for _, region := range regions.Regions.Region {
		fmt.Println(region.RegionID)
	}
}

func (regions DescribeRegions) PrintTable() {
	regions.Print()
}
