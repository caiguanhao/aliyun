package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

const instanceDateTimeFormat = "2006-01-02T15:04Z"

type ECSInstances []ECSInstance

func (a ECSInstances) Len() int      { return len(a) }
func (a ECSInstances) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ECSInstances) Less(i, j int) bool {
	b, be := time.Parse(instanceDateTimeFormat, a[i].CreationTime)
	c, ce := time.Parse(instanceDateTimeFormat, a[j].CreationTime)
	if be != nil || ce != nil {
		return a[i].CreationTime > a[j].CreationTime
	}
	return b.After(c)
}

type DescribeInstances struct {
	Instances struct {
		Instance ECSInstances `json:"Instance"`
	} `json:"Instances"`
	PageNumber int64  `json:"PageNumber"`
	PageSize   int64  `json:"PageSize"`
	RequestId  string `json:"RequestId"`
	TotalCount int64  `json:"TotalCount"`
}

var showAll bool
var showHiddenOnly bool
var matchRegexes []*regexp.Regexp
var makeHosts bool
var usePrivateIPAddr bool

var DESCRIBE_INSTANCES cli.Command = cli.Command{
	Name:      "list-instances",
	Aliases:   []string{"list", "ls", "l"},
	Usage:     "list all ECS instances of all regions",
	ArgsUsage: "[instance IDs...]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "also show hidden instances, overrides --hidden-only, but not --regex",
			Destination: &showAll,
		},
		cli.BoolFlag{
			Name:        "hidden-only, H",
			Usage:       "show only hidden instances",
			Destination: &showHiddenOnly,
		},
		cli.StringSliceFlag{
			Name:  "regex",
			Usage: "show instance with the name matching regex",
		},
		cli.BoolFlag{
			Name:        "hosts",
			Usage:       "print host names and IP addresses for /etc/hosts",
			Destination: &makeHosts,
		},
		cli.BoolFlag{
			Name:        "private",
			Usage:       "print private IP address when --hosts",
			Destination: &usePrivateIPAddr,
		},
	},
	Action: func(c *cli.Context) {
		for _, regex := range c.StringSlice("regex") {
			re, err := regexp.Compile(regex)
			if err != nil {
				exit(err)
			}
			matchRegexes = append(matchRegexes, re)
		}
		if c.Args().Present() {
			ForAllArgsDo([]string(c.Args()), func(arg string) {
				Print(ECS_INSTANCE.DescribeInstanceAttributeById(arg))
			})
		} else {
			Print(ECS_INSTANCE.DescribeInstances())
		}
	},
	BashComplete: func(c *cli.Context) {
		printFlagsForCommand(c, "list-instances")
		describeInstancesForBashComplete(nil)(c)
	},
}

func (ecs *ECS) DescribeInstances() (instances ECSInstances, err error) {
	err = ForAllRegionsDo(func(region string) (err error) {
		var resp DescribeInstances
		err = ecs.Request(map[string]string{
			"Action":   "DescribeInstances",
			"RegionId": region,
		}, &resp)
		if err == nil {
			instances = append(instances, resp.Instances.Instance...)
		}
		return
	})
	sort.Sort(instances)
	return
}

func (instances ECSInstances) Print() {
	for _, instance := range instances {
		if !shouldShow(instance) {
			continue
		}
		fmt.Println(instance.InstanceId)
	}
}

func (instances ECSInstances) PrintTable() {
	if makeHosts {
		PrintTable([]interface{}{"# IP Address", "Name"}, len(instances), func(i int) []interface{} {
			instance := instances[i]
			if !shouldShow(instance) {
				return nil
			}
			ipAddr := instance.InnerIpAddress.GetIPAddress(0)
			if !usePrivateIPAddr {
				ipAddr = instance.PublicIpAddress.GetIPAddress(0)
			}
			if ipAddr == "" {
				return nil
			}
			return []interface{}{ipAddr, instance.InstanceName}
		})
		return
	}

	fields := []interface{}{"ID", "Name", "Status", "Public IP", "Private IP", "Type", "Region/Zone", "Created At"}
	PrintTable(fields, len(instances), func(i int) []interface{} {
		instance := instances[i]
		if !shouldShow(instance) {
			return nil
		}
		return []interface{}{
			instance.InstanceId,
			instance.InstanceName,
			instance.Status,
			instance.PublicIpAddress.GetIPAddress(0),
			instance.InnerIpAddress.GetIPAddress(0),
			instance.InstanceType,
			instance.ZoneId,
			dateStr(instance.CreationTime),
		}
	})
}

func dateStr(input string) (output string) {
	createdAt, _ := time.Parse(instanceDateTimeFormat, input)
	output = fmt.Sprintf("%s (%.0f days ago)",
		createdAt.Local().Format(YMD_HMS_FORMAT),
		math.Floor(time.Since(createdAt).Hours()/24))
	return
}

func shouldShow(instance ECSInstance) (shouldShow bool) {
	shouldShow = true

	if !showAll && strings.Contains(instance.Description, "[HIDE]") {
		shouldShow = false
	}

	for _, regex := range matchRegexes {
		if regex.MatchString(instance.InstanceName) {
			continue
		}
		shouldShow = false
	}

	return
}
