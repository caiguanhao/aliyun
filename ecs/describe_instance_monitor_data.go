package main

import (
	"fmt"
	"time"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

type ECSInstanceMonitorData []ECSInstanceMonitorDatum

type ECSInstanceMonitorDatum struct {
	BPSRead           int    `json:"BPSRead"`
	BPSWrite          int    `json:"BPSWrite"`
	CPU               int    `json:"CPU"`
	IOPSRead          int    `json:"IOPSRead"`
	IOPSWrite         int    `json:"IOPSWrite"`
	InstanceID        string `json:"InstanceId"`
	InternetBandwidth int    `json:"InternetBandwidth"`
	InternetFlow      int    `json:"InternetFlow"`
	InternetRX        int    `json:"InternetRX"`
	InternetTX        int    `json:"InternetTX"`
	IntranetBandwidth int    `json:"IntranetBandwidth"`
	IntranetFlow      int    `json:"IntranetFlow"`
	IntranetRX        int    `json:"IntranetRX"`
	IntranetTX        int    `json:"IntranetTX"`
	TimeStamp         string `json:"TimeStamp"`
}

type DescribeInstanceMonitorData struct {
	MonitorData struct {
		InstanceMonitorData ECSInstanceMonitorData `json:"InstanceMonitorData"`
	} `json:"MonitorData"`
	RequestID string `json:"RequestId"`
}

var period int

var DESCRIBE_INSTANCE_MONITOR_DATA cli.Command = cli.Command{
	Name:      "monitor-instance",
	Aliases:   []string{"monitor", "m"},
	Usage:     "show CPU and network usage history of an instance",
	ArgsUsage: "[instance IDs...]",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "hours, H",
			Value: 1,
			Usage: "show stats from how many hours ago till now",
		},
		cli.IntFlag{
			Name:  "period, p",
			Value: 0,
			Usage: "period in seconds; must be: 60, 600 or 3600; otherwise will use smaller period as possible",
		},
	},
	Action: func(c *cli.Context) {
		now := time.Now().UTC()
		then := now.Add(time.Duration(-1*c.Int("hours")) * time.Hour)
		period = c.Int("period")
		if period != 60 && period != 600 && period != 3600 {
			if now.Sub(then).Seconds()/60 <= 400 {
				period = 60
			} else if now.Sub(then).Seconds()/600 <= 400 {
				period = 600
			} else {
				period = 3600
			}
		}
		ForAllArgsDo([]string(c.Args()), func(arg string) {
			Print(ECS_INSTANCE.DescribeInstanceMonitorData(arg, then, now, period))
		})
	},
	BashComplete: describeInstancesForBashComplete(nil),
}

func (ecs *ECS) DescribeInstanceMonitorData(id string, startTime, endTime time.Time, period int) (_ ECSInstanceMonitorData, resp DescribeInstanceMonitorData, err error) {
	return resp.MonitorData.InstanceMonitorData, resp, ecs.Request(map[string]string{
		"Action":     "DescribeInstanceMonitorData",
		"InstanceId": id,
		"StartTime":  startTime.Format(TIME_FORMAT),
		"EndTime":    endTime.Format(TIME_FORMAT),
		"Period":     fmt.Sprintf("%d", period),
	}, &resp)
}

func (data ECSInstanceMonitorData) Print() {
	for _, datum := range data {
		fmt.Println(datum.TimeStamp)
	}
}

func (data ECSInstanceMonitorData) PrintTable() {
	PrintTable(
		/* fields     */ []interface{}{"Time", "CPU Usage", "Received", "Sent"},
		/* showFields */ true,
		/* listLength */ len(data),
		/* filter     */ nil,
		/* getInfo    */ func(i int) map[interface{}]interface{} {
			datum := data[i]
			t, _ := time.Parse(TIME_FORMAT, datum.TimeStamp)
			return map[interface{}]interface{}{
				"Time":      t.Local().Format(YMD_HMS_FORMAT),
				"CPU Usage": fmt.Sprintf("%d%%", datum.CPU),
				"Received":  fmt.Sprintf("%d KB/s", datum.InternetRX/8/period),
				"Sent":      fmt.Sprintf("%d KB/s", datum.InternetTX/8/period),
			}
		},
	)
}
