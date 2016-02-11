package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ECSRegion struct {
	RegionID string `json:"RegionId"`
}

var showZonesOnly bool

var DESCRIBE_REGIONS cli.Command = cli.Command{
	Name:      "list-regions",
	Aliases:   []string{"regions", "n"},
	Usage:     "list all available regions and zones",
	ArgsUsage: " ",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:        "zones-only, z",
			Usage:       "show only zones (when --quiet is set)",
			Destination: &showZonesOnly,
		},
	},
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.DescribeRegionsAndZones())
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

func (regions ECSRegions) Print() {
	sort.Sort(regions)
	for _, region := range regions {
		fmt.Println(region.RegionID)
	}
}

func (regions ECSRegions) PrintTable() {
	regions.Print()
}

type ECSZone struct {
	AvailableDiskCategories struct {
		DiskCategories []string `json:"DiskCategories"`
	} `json:"AvailableDiskCategories"`
	AvailableResourceCreation struct {
		ResourceTypes []string `json:"ResourceTypes"`
	} `json:"AvailableResourceCreation"`
	LocalName string `json:"LocalName"`
	ZoneID    string `json:"ZoneId"`
}

type DescribeZones struct {
	RequestID string `json:"RequestId"`
	Zones     struct {
		Zone []ECSZone `json:"Zone"`
	} `json:"Zones"`
}

type ECSZones []ECSZone

func (a ECSZones) Len() int           { return len(a) }
func (a ECSZones) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ECSZones) Less(i, j int) bool { return a[i].ZoneID < a[j].ZoneID }

func (zones ECSZones) Print() {
	sort.Sort(zones)
	for _, zone := range zones {
		fmt.Println(zone.ZoneID)
	}
}

func (zones ECSZones) PrintTable() {
	zones.Print()
}

func (ecs *ECS) DescribeRegions() (_ ECSRegions, resp DescribeRegions, _ error) {
	return resp.Regions.Region, resp, ecs.Request(map[string]string{
		"Action": "DescribeRegions",
	}, &resp)
}

func (ecs *ECS) DescribeZones(region string) (_ ECSZones, resp DescribeZones, _ error) {
	return resp.Zones.Zone, resp, ecs.Request(map[string]string{
		"Action":   "DescribeZones",
		"RegionId": region,
	}, &resp)
}

type ECSRegionsAndZones []struct {
	RegionName string
	Zones      []string
}

func (a ECSRegionsAndZones) Len() int {
	for _, _a := range a {
		sort.Sort(sort.StringSlice(_a.Zones))
	}
	return len(a)
}
func (a ECSRegionsAndZones) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ECSRegionsAndZones) Less(i, j int) bool { return a[i].RegionName < a[j].RegionName }

func (ecs *ECS) DescribeRegionsAndZones() (regionsNzones ECSRegionsAndZones, err error) {
	err = ForAllRegionsDo(func(region string) (err error) {
		var zones ECSZones
		zones, _, err = ecs.DescribeZones(region)
		if err == nil {
			var rNz struct {
				RegionName string
				Zones      []string
			}
			rNz.RegionName = region
			for _, zone := range zones {
				rNz.Zones = append(rNz.Zones, zone.ZoneID)
			}
			regionsNzones = append(regionsNzones, rNz)
		}
		return
	})
	return
}

func (regionsNzones ECSRegionsAndZones) Print() {
	sort.Sort(regionsNzones)
	for _, region := range regionsNzones {
		if showZonesOnly {
			for _, zone := range region.Zones {
				fmt.Println(zone)
			}
		} else {
			fmt.Println(region.RegionName)
		}
	}
}

func (regionsNzones ECSRegionsAndZones) PrintTable() {
	sort.Sort(regionsNzones)
	fields := []interface{}{"Region", "Available Zones"}
	PrintTable(fields, len(regionsNzones), func(i int) []interface{} {
		rNz := regionsNzones[i]
		return []interface{}{
			rNz.RegionName,
			strings.Join(rNz.Zones, ", "),
		}
	})
}
