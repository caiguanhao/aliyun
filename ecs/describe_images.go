package main

import (
	"fmt"
	"sort"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ECSImage struct {
	Architecture       string `json:"Architecture"`
	CreationTime       string `json:"CreationTime"`
	Description        string `json:"Description"`
	DiskDeviceMappings struct {
		DiskDeviceMapping []struct {
			Device     string `json:"Device"`
			Size       string `json:"Size"`
			SnapshotId string `json:"SnapshotId"`
		} `json:"DiskDeviceMapping"`
	} `json:"DiskDeviceMappings"`
	ImageId         string `json:"ImageId"`
	ImageName       string `json:"ImageName"`
	ImageOwnerAlias string `json:"ImageOwnerAlias"`
	ImageVersion    string `json:"ImageVersion"`
	IsSubscribed    bool   `json:"IsSubscribed"`
	OSName          string `json:"OSName"`
	ProductCode     string `json:"ProductCode"`
	Size            int64  `json:"Size"`
}

var DESCRIBE_IMAGES cli.Command = cli.Command{
	Name:      "list-images",
	Aliases:   []string{"images", "i"},
	Usage:     "show info of all images",
	ArgsUsage: " ",
	Action: func(c *cli.Context) {
		Print(ECS_INSTANCE.DescribeImages())
	},
}

type DescribeImages struct {
	Images struct {
		Image ECSImages `json:"Image"`
	} `json:"Images"`
	PageNumber int64  `json:"PageNumber"`
	PageSize   int64  `json:"PageSize"`
	RegionId   string `json:"RegionId"`
	RequestId  string `json:"RequestId"`
	TotalCount int64  `json:"TotalCount"`
}

type ECSImages []ECSImage

func (a ECSImages) Len() int           { return len(a) }
func (a ECSImages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ECSImages) Less(i, j int) bool { return a[i].ImageId < a[j].ImageId }

func (ecs *ECS) DescribeImages() (images ECSImages, resp DescribeImages, _ error) {
	defer func() {
		sort.Sort(images)
	}()
	return resp.Images.Image, resp, ecs.Request(map[string]string{
		"Action":   "DescribeImages",
		"RegionId": "",
	}, &resp)
}

func (images ECSImages) Print() {
	for _, image := range images {
		fmt.Println(image.ImageId)
	}
}

func (images ECSImages) PrintTable() {
	fields := []interface{}{"ID", "Owner", "Name"}
	PrintTable(fields, len(images), func(i int) []interface{} {
		image := images[i]
		name := image.ImageName
		if name == image.ImageId {
			name = "-"
		}
		return []interface{}{
			image.ImageId,
			image.ImageOwnerAlias,
			name,
		}
	})
}
