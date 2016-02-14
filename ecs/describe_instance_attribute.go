package main

import (
	"fmt"
	"math"
	"time"
)

type ECSInstance DescribeInstanceAttribute

type DescribeInstanceAttributeIPAddress struct {
	IpAddress []string `json:"IpAddress"`
}

func (ipaddr DescribeInstanceAttributeIPAddress) GetIPAddress(n int) string {
	if n > len(ipaddr.IpAddress)-1 {
		return ""
	}
	return ipaddr.IpAddress[n]
}

type DescribeInstanceAttribute struct {
	ClusterId    string `json:"ClusterId"`
	CreationTime string `json:"CreationTime"`
	Description  string `json:"Description"`
	EipAddress   struct {
		AllocationId       string `json:"AllocationId"`
		InternetChargeType string `json:"InternetChargeType"`
		IpAddress          string `json:"IpAddress"`
	} `json:"EipAddress"`
	HostName                string                             `json:"HostName"`
	ImageId                 string                             `json:"ImageId"`
	InnerIpAddress          DescribeInstanceAttributeIPAddress `json:"InnerIpAddress"`
	InstanceId              string                             `json:"InstanceId"`
	InstanceName            string                             `json:"InstanceName"`
	InstanceNetworkType     string                             `json:"InstanceNetworkType"`
	InstanceType            string                             `json:"InstanceType"`
	InternetChargeType      string                             `json:"InternetChargeType"`
	InternetMaxBandwidthIn  int64                              `json:"InternetMaxBandwidthIn"`
	InternetMaxBandwidthOut int64                              `json:"InternetMaxBandwidthOut"`
	OperationLocks          struct {
		LockReason []struct {
			LockReason string `json:"LockReason"`
		} `json:"LockReason"`
	} `json:"OperationLocks"`
	PublicIpAddress  DescribeInstanceAttributeIPAddress `json:"PublicIpAddress"`
	RegionId         string                             `json:"RegionId"`
	SecurityGroupIds struct {
		SecurityGroupId []string `json:"SecurityGroupId"`
	} `json:"SecurityGroupIds"`
	Status        string `json:"Status"`
	VlanId        string `json:"VlanId"`
	VpcAttributes struct {
		NatIpAddress     string                             `json:"NatIpAddress"`
		PrivateIpAddress DescribeInstanceAttributeIPAddress `json:"PrivateIpAddress"`
		VSwitchId        string                             `json:"VSwitchId"`
		VpcId            string                             `json:"VpcId"`
	} `json:"VpcAttributes"`
	ZoneId string `json:"ZoneId"`
}

func (ecs *ECS) DescribeInstanceAttributeById(id string) (instance ECSInstance, err error) {
	return instance, ecs.Request(map[string]string{
		"Action":     "DescribeInstanceAttribute",
		"InstanceId": id,
	}, &instance)
}

func (instance ECSInstance) Print() {
	fmt.Println(instance.InstanceId)
}

func (instance ECSInstance) PrintTable() {
	const format = "%15s: %s\n"
	createdAt, _ := time.Parse(time.RFC3339, instance.CreationTime)
	duration := time.Since(createdAt)
	createdAtStr := fmt.Sprintf("%s (%.0f days ago)",
		createdAt.Local().Format(YMD_HMS_FORMAT),
		math.Floor(duration.Hours()/24))
	fmt.Printf(format, "ID", instance.InstanceId)
	fmt.Printf(format, "Name", instance.InstanceName)
	fmt.Printf(format, "Type", instance.InstanceType)
	fmt.Printf(format, "Image", instance.ImageId)
	fmt.Printf(format, "Status", instance.Status)
	fmt.Printf(format, "Region", instance.RegionId)
	fmt.Printf(format, "Zone", instance.ZoneId)
	fmt.Printf(format, "Public IP", instance.PublicIpAddress.GetIPAddress(0))
	fmt.Printf(format, "Private IP", instance.InnerIpAddress.GetIPAddress(0))
	fmt.Printf(format, "Created At", createdAtStr)
	fmt.Printf(format, "Description", instance.Description)
}
