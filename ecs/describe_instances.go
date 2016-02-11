package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ECSInstances []ECSInstance

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

var DESCRIBE_INSTANCES cli.Command = cli.Command{
	Name:      "list-instances",
	Aliases:   []string{"list", "l"},
	Usage:     "list all ECS instances of all regions",
	ArgsUsage: "[instance IDs...]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "also show hidden instances, overrides --hidden-only",
			Destination: &showAll,
		},
		cli.BoolFlag{
			Name:        "hidden-only, H",
			Usage:       "show only hidden instances",
			Destination: &showHiddenOnly,
		},
	},
	Action: func(c *cli.Context) {
		if c.Args().Present() {
			ForAllArgsDo([]string(c.Args()), func(arg string) {
				Print(ECS_INSTANCE.DescribeInstanceAttributeById(arg))
			})
		} else {
			Print(ECS_INSTANCE.DescribeInstances())
		}
	},
	BashComplete: DescribeInstancesForBashComplete,
}

func DescribeInstancesForBashComplete(c *cli.Context) {
	instances, _ := ECS_INSTANCE.DescribeInstances()
	for _, instance := range instances {
		fmt.Printf("%s@%s\n", instance.InstanceId, instance.InstanceName)
	}
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
	createdAt, _ := time.Parse("2006-01-02T15:04Z", input)
	output = fmt.Sprintf("%s (%.0f days ago)",
		createdAt.Local().Format("2006-01-02 15:04:05"),
		math.Floor(time.Since(createdAt).Hours()/24))
	return
}

func shouldShow(instance ECSInstance) bool {
	if showAll {
		return true
	} else {
		isHidden := strings.Contains(instance.Description, "[HIDE]")
		if showHiddenOnly {
			return isHidden
		} else {
			return !isHidden
		}
	}
}
