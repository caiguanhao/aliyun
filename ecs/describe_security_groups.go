package main

import (
	"fmt"
	"sort"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ECSSecurityGroup struct {
	Description     string `json:"Description"`
	SecurityGroupId string `json:"SecurityGroupId"`

	regionId string
}

type ECSSecurityGroups []ECSSecurityGroup

func (a ECSSecurityGroups) Len() int           { return len(a) }
func (a ECSSecurityGroups) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ECSSecurityGroups) Less(i, j int) bool { return a[i].SecurityGroupId < a[j].SecurityGroupId }

var DESCRIBE_SECURITY_GROUPS cli.Command = cli.Command{
	Name:      "list-security-groups",
	Aliases:   []string{"groups", "g"},
	Usage:     "list all security groups",
	ArgsUsage: " ",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.DescribeSecurityGroups())
	},
}

type DescribeSecurityGroups struct {
	PageNumber     int64  `json:"PageNumber"`
	PageSize       int64  `json:"PageSize"`
	RegionId       string `json:"RegionId"`
	RequestId      string `json:"RequestId"`
	SecurityGroups struct {
		SecurityGroup ECSSecurityGroups `json:"SecurityGroup"`
	} `json:"SecurityGroups"`
	TotalCount int64 `json:"TotalCount"`
}

func (ecs *ECS) DescribeSecurityGroups() (groups ECSSecurityGroups, err error) {
	err = ForAllRegionsDo(func(region string) (err error) {
		var resp DescribeSecurityGroups
		err = ecs.Request(map[string]string{
			"Action":   "DescribeSecurityGroups",
			"RegionId": region,
		}, &resp)
		if err == nil {
			for _, group := range resp.SecurityGroups.SecurityGroup {
				group.regionId = region
				groups = append(groups, group)
			}
		}
		return
	})
	sort.Sort(groups)
	return
}

func (groups ECSSecurityGroups) Print() {
	for _, group := range groups {
		fmt.Println(group.SecurityGroupId)
	}
}

func (groups ECSSecurityGroups) PrintTable() {
	PrintTable(
		/* fields     */ []interface{}{"ID", "Description", "Region"},
		/* showFields */ true,
		/* listLength */ len(groups),
		/* filter     */ nil,
		/* getInfo    */ func(i int) map[interface{}]interface{} {
			group := groups[i]
			return map[interface{}]interface{}{
				"ID":          group.SecurityGroupId,
				"Description": group.Description,
				"Region":      group.regionId,
			}
		},
	)
}
