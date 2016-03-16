package main

import "encoding/xml"

type OSSResponseError struct {
	XMLName        xml.Name `xml:"Error"`
	Code           string   `xml:"Code"`
	Message        string   `xml:"Message"`
	RequestId      string   `xml:"RequestId"`
	HostId         string   `xml:"HostId"`
	OSSAccessKeyId string   `xml:"OSSAccessKeyId"`
}
