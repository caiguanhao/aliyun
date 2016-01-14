package main

import (
	"errors"
	"fmt"
	"sort"

	"github.com/caiguanhao/aliyun/ecs/opts"
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

type DescribeImages struct {
	Images struct {
		Image []ECSImage `json:"Image"`
	} `json:"Images"`
	PageNumber int64  `json:"PageNumber"`
	PageSize   int64  `json:"PageSize"`
	RegionId   string `json:"RegionId"`
	RequestId  string `json:"RequestId"`
	TotalCount int64  `json:"TotalCount"`
}

type byImagesId []ECSImage

func (a byImagesId) Len() int           { return len(a) }
func (a byImagesId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byImagesId) Less(i, j int) bool { return a[i].ImageId < a[j].ImageId }

func (images *DescribeImages) Do(ecs *ECS) (*DescribeImages, error) {
	if len(opts.Region) < 1 {
		if opts.IsQuiet {
			return images, nil
		}
		return nil, errors.New("Please provide a --region.")
	}
	return images, ecs.Request(map[string]string{
		"Action":   "DescribeImages",
		"RegionId": opts.Region,
	}, images)
}

func (images DescribeImages) Print() {
	sort.Sort(byImagesId(images.Images.Image))
	for _, image := range images.Images.Image {
		fmt.Println(image.ImageId)
	}
}

func (images DescribeImages) PrintTable() {
	idMaxLength := 2
	ownerMaxLength := 3
	for _, image := range images.Images.Image {
		idLength := len(image.ImageId)
		ownerLength := len(image.ImageOwnerAlias)
		if idLength > idMaxLength {
			idMaxLength = idLength
		}
		if ownerLength > ownerMaxLength {
			ownerMaxLength = ownerLength
		}
	}
	format := fmt.Sprintf("%%-%ds  %%-%ds  %%s\n", idMaxLength, ownerMaxLength)
	sort.Sort(byImagesId(images.Images.Image))
	fmt.Printf(format, "ID", "Owner", "Name")
	for _, image := range images.Images.Image {
		name := image.ImageName
		if name == image.ImageId {
			name = ""
		}
		fmt.Printf(
			format,
			image.ImageId,
			image.ImageOwnerAlias,
			name,
		)
	}
}
