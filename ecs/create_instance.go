package main

import (
	"fmt"
	"strings"

	"github.com/caiguanhao/aliyun/ecs/errors"
	"github.com/caiguanhao/aliyun/vendor/cli"
)

type CreateInstance struct {
	InstanceId string `json:"InstanceId"`
	RequestId  string `json:"RequestId"`
}

var DEFAULT_INCOMING_BANDWIDTH = 200
var DEFAULT_OUTGOING_BANDWIDTH = 5

var flags []cli.Flag = []cli.Flag{
	cli.StringFlag{
		Name:  "image, i",
		Usage: "create using this image",
	},
	cli.StringFlag{
		Name:  "type, t",
		Usage: "type of the new instance",
	},
	cli.StringFlag{
		Name:  "name, n",
		Usage: "name of the new instance",
	},
	cli.StringFlag{
		Name:  "host, H",
		Usage: "host name of the new instance, defaults to value of --name",
	},
	cli.StringFlag{
		Name:  "group, g",
		Usage: "put the new instance in to this group",
	},
	cli.StringFlag{
		Name:  "region, r",
		Usage: "put the new instance in to this region",
	},
	cli.StringFlag{
		Name:  "zone, z",
		Usage: "put the new instance in to this zone, use random zone if not specified",
	},
	cli.StringSliceFlag{
		Name:  "disk, d",
		Usage: "specify data disk size in GB ranges from 5 to 2000 (can be specified more than once; no data disk by default)",
	},
	cli.StringFlag{
		Name:  "incoming-bandwidth, I",
		Value: fmt.Sprintf("%d", DEFAULT_INCOMING_BANDWIDTH),
		Usage: fmt.Sprintf("maximum incoming bandwidth in Mbps ranges from 1 to 200, default is %d (free in charge)", DEFAULT_INCOMING_BANDWIDTH),
	},
	cli.StringFlag{
		Name:  "outgoing-bandwidth, O",
		Value: fmt.Sprintf("%d", DEFAULT_OUTGOING_BANDWIDTH),
		Usage: fmt.Sprintf("maximum outgoing bandwidth in Mbps ranges from 1 to 200, default is %d (pay per use)", DEFAULT_OUTGOING_BANDWIDTH),
	},
	cli.StringFlag{
		Name:   "password, p",
		Usage:  "password of the new instance, can be specified from env var",
		EnvVar: "PASSWORD",
	},
}

var CREATE_INSTANCE cli.Command = cli.Command{
	Name:      "create-instance",
	Aliases:   []string{"create", "c"},
	Usage:     "create an instance",
	ArgsUsage: " ",
	Flags:     flags,
	Action: func(c *cli.Context) {
		if checkValuesForBashComplete(c) {
			return
		}
		host := c.String("host")
		if host == "" {
			host = c.String("name")
		}
		params := map[string]string{
			"ImageId":                 c.String("image"),
			"InstanceType":            c.String("type"),
			"SecurityGroupId":         c.String("group"),
			"InstanceName":            c.String("name"),
			"HostName":                host,
			"RegionId":                c.String("region"),
			"ZoneId":                  c.String("zone"),
			"Password":                c.String("password"),
			"InternetMaxBandwidthIn":  c.String("incoming-bandwidth"),
			"InternetMaxBandwidthOut": c.String("outgoing-bandwidth"),
			"InternetChargeType":      "PayByTraffic",
			"SystemDisk.Category":     "cloud",
		}
		for i, size := range c.StringSlice("disk") {
			params[fmt.Sprintf("DataDisk.%d.Size", i+1)] = size
		}
		Print(ECS_INSTANCE.CreateInstance(params))
	},
	BashComplete: func(c *cli.Context) {
		hintsForBashComplete(c, nil)
	},
}

func checkValuesForBashComplete(c *cli.Context) bool {
	bashCompletionFlag := "--" + cli.BashCompletionFlag.Name
	for _, flag := range flags {
		name := strings.Split(flag.GetName(), ",")[0]
		switch flag.(type) {
		case cli.StringSliceFlag:
			for _, value := range c.StringSlice(name) {
				if value == bashCompletionFlag {
					hintsForBashComplete(c, &name)
					return true
				}
			}
		default:
			value := c.String(name)
			if value == bashCompletionFlag {
				hintsForBashComplete(c, &name)
				return true
			}
		}
	}
	return false
}

func hintsForBashComplete(c *cli.Context, flagName *string) {
	if flagName == nil {
		for _, flag := range flags {
			for _, name := range strings.Split(flag.GetName(), ",") {
				name = strings.TrimSpace(name)
				fmt.Print("-")
				if len(name) > 1 {
					fmt.Print("-")
				}
				fmt.Println(name)
			}
		}
	} else if *flagName == "disk" {
		fmt.Println(5, 10, 100, 200, 500, 1000, 2000)
	} else if *flagName == "group" {
		groups, _ := ECS_INSTANCE.DescribeSecurityGroups()
		for _, group := range groups {
			fmt.Println(group.SecurityGroupId)
		}
	} else if *flagName == "host" || *flagName == "name" {
		instances, _ := ECS_INSTANCE.DescribeInstances()
		for _, instance := range instances {
			fmt.Println(instance.InstanceName)
		}
	} else if *flagName == "image" {
		images, _, _ := ECS_INSTANCE.DescribeImages()
		for _, image := range images {
			fmt.Println(image.ImageId)
		}
	} else if *flagName == "incoming-bandwidth" {
		fmt.Println(DEFAULT_INCOMING_BANDWIDTH)
	} else if *flagName == "outgoing-bandwidth" {
		fmt.Println(DEFAULT_OUTGOING_BANDWIDTH)
	} else if *flagName == "region" {
		regions, _, _ := ECS_INSTANCE.DescribeRegions()
		for _, region := range regions {
			fmt.Println(region.RegionID)
		}
	} else if *flagName == "type" {
		types, _, _ := ECS_INSTANCE.DescribeInstanceTypes()
		for _, _type := range types {
			fmt.Println(_type.InstanceTypeId)
		}
	} else if *flagName == "zone" {
		region := c.String("region")
		if region == "" {
			regions, _ := ECS_INSTANCE.DescribeRegionsAndZones()
			for _, region := range regions {
				for _, zone := range region.Zones {
					fmt.Println(zone)
				}
			}
		} else {
			zones, _, _ := ECS_INSTANCE.DescribeZones(region)
			for _, zone := range zones {
				fmt.Println(zone.ZoneID)
			}
		}
	}
}

func (ecs *ECS) CreateInstance(_params map[string]string) (resp CreateInstance, _ error) {
	params := map[string]string{
		"Action": "CreateInstance",
	}
	for k, v := range _params {
		params[k] = v
	}
	var errs errors.Errors
	for k, v := range map[string]string{
		"Password":        "password",
		"ImageId":         "image",
		"InstanceType":    "type",
		"SecurityGroupId": "group",
		"InstanceName":    "name",
		"RegionId":        "region",
	} {
		if len(params[k]) < 1 {
			errs.Add(fmt.Sprintf("Please provide --%s.", v))
		}
	}
	if errs.HaveError() {
		return resp, errs.Errorify()
	}
	return resp, ecs.Request(params, &resp)
}

func (create CreateInstance) Print() {
	fmt.Println(create.InstanceId)
}

func (create CreateInstance) PrintTable() {
	fmt.Println(create.InstanceId)
}
