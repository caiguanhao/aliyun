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

var regionAvailabilityChan chan func(string) bool = make(chan func(string) bool)

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
		go func() {
			regionAvailabilityChan <- getRegionAvailability()
		}()
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
		Zone ECSZones `json:"Zone"`
	} `json:"Zones"`
}

type ECSZones []ECSZone

func (a ECSZones) Len() int           { return len(a) }
func (a ECSZones) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ECSZones) Less(i, j int) bool { return a[i].ZoneID < a[j].ZoneID }

func (zones ECSZones) Print() {
	for _, zone := range zones {
		fmt.Println(zone.ZoneID)
	}
}

func (zones ECSZones) PrintTable() {
	zones.Print()
}

func (ecs *ECS) DescribeRegions() (regions ECSRegions, resp DescribeRegions, _ error) {
	defer func() {
		sort.Sort(regions)
	}()
	return resp.Regions.Region, resp, ecs.Request(map[string]string{
		"Action": "DescribeRegions",
	}, &resp)
}

func (ecs *ECS) DescribeZones(region string) (zones ECSZones, resp DescribeZones, _ error) {
	defer func() {
		sort.Sort(zones)
	}()
	return resp.Zones.Zone, resp, ecs.Request(map[string]string{
		"Action":   "DescribeZones",
		"RegionId": region,
	}, &resp)
}

type ECSRegionsAndZones []struct {
	RegionName string
	Zones      []string
}

func (a ECSRegionsAndZones) Len() int           { return len(a) }
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
	sort.Sort(regionsNzones)
	return
}

func (regionsNzones ECSRegionsAndZones) Print() {
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
	isAvailable := <-regionAvailabilityChan
	PrintTable(
		/* fields     */ []interface{}{"Region", "Available", "Zones"},
		/* showFields */ true,
		/* listLength */ len(regionsNzones),
		/* filter     */ nil,
		/* getInfo    */ func(i int) map[interface{}]interface{} {
			rNz := regionsNzones[i]
			available := "No"
			if isAvailable(rNz.RegionName) {
				available = "Yes"
			}
			return map[interface{}]interface{}{
				"Region":    rNz.RegionName,
				"Available": available,
				"Zones":     strings.Join(rNz.Zones, ", "),
			}
		},
	)
}
