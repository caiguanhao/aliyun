package main

import (
	"regexp"
	"strings"
)

type ECSRegionStatuses []ECSRegionStatus

type ECSRegionStatus struct {
	Disabled bool   `json:"disabled"`
	Title    string `json:"text"`
	Name     string `json:"value"`
}

type GetCommodity struct {
	Code string `json:"code"`
	Data struct {
		Components struct {
			VMRegionNo struct {
				VMRegionNo ECSRegionStatuses `json:"vm_region_no"`
			} `json:"vm_region_no"`
		} `json:"components"`
	} `json:"data"`
	RequestID       string `json:"requestId"`
	SuccessResponse bool   `json:"successResponse"`
}

func getRegionAvailability() func(check string) bool {
	var resp GetCommodity
	Request("https://ecs-buy.aliyun.com/ecs/getCommodity.json?commodityCode=ecs&orderType=BUY", &resp)
	re := regexp.MustCompile("^([^-]+-[^-]+).*$")
	return func(check string) bool {
		check = re.ReplaceAllString(check, "$1")
		for _, region := range resp.Data.Components.VMRegionNo.VMRegionNo {
			if strings.Contains(region.Name, check) {
				if region.Disabled == false {
					return true
				}
			}
		}
		return false
	}
}
