package main

import (
	"fmt"

	"github.com/caiguanhao/aliyun/ecs/errors"
	"github.com/caiguanhao/aliyun/vendor/cli"
)

type CreateInstance struct {
	InstanceId string `json:"InstanceId"`
	RequestId  string `json:"RequestId"`
}

var CREATE_INSTANCE cli.Command = cli.Command{
	Name:    "create-instance",
	Aliases: []string{"create", "c"},
	Usage:   "create an instance",
	Flags: []cli.Flag{
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
			Value: "200",
			Usage: "maximum incoming bandwidth in Mbps ranges from 1 to 200, default is 200 (free in charge)",
		},
		cli.StringFlag{
			Name:  "outgoing-bandwidth, O",
			Value: "5",
			Usage: "maximum outgoing bandwidth in Mbps ranges from 1 to 200, default is 5 (pay per use)",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "password of the new instance, can be specified from env var",
			EnvVar: "PASSWORD",
		},
	},
	Action: func(c *cli.Context) {
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
