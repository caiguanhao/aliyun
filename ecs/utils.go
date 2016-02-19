package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/caiguanhao/aliyun/ecs/errors"
	"github.com/caiguanhao/aliyun/vendor/cli"
)

const TIME_FORMAT = "2006-01-02T15:04:05Z"
const YMD_HMS_FORMAT = "2006-01-02 15:04:05"

func exit(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}

func sign(secret string, query string) string {
	mac := hmac.New(sha1.New, []byte(secret+"&"))
	mac.Write([]byte("GET&%2F&" + query))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func randomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func urlEncode(input string) string {
	return strings.Replace(url.QueryEscape(input), "+", "%20", -1)
}

func buildQueryString(input map[string]string) string {
	keys := make([]string, 0, len(input))
	for val := range input {
		keys = append(keys, val)
	}
	sort.Strings(keys)
	queries := make([]string, 0, len(input))
	for _, key := range keys {
		query := fmt.Sprintf("%s=%s", urlEncode(key), urlEncode(input[key]))
		queries = append(queries, query)
	}
	queryString := strings.Join(queries, "&")
	return queryString
}

func (ecs *ECS) Request(queries map[string]string, target interface{}) error {
	params := map[string]string{
		"Format":           "JSON",
		"Version":          "2014-05-26",
		"AccessKeyId":      ecs.KEY,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   randomString(64),
		"Timestamp":        time.Now().UTC().Format(TIME_FORMAT),
		"PageSize":         "50",
		"PageNumber":       "1",
	}
	for k, v := range queries {
		params[k] = v
	}
	query := buildQueryString(params)
	signature := sign(ecs.SECRET, urlEncode(query))
	url := fmt.Sprintf("http://ecs.aliyuncs.com/?%s&Signature=%s", query, urlEncode(signature))

	if IsVerbose {
		fmt.Println(url)
	}

	client := http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if IsVerbose {
		fmt.Printf("%s %s\n", res.Proto, res.Status)
		for key, values := range res.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
		body, ioerr := ioutil.ReadAll(res.Body)
		if ioerr != nil {
			return ioerr
		}
		if res.StatusCode == 200 {
			err = json.Unmarshal(body, target)
		} else {
			errResp := ECSResponseError{}
			err = json.Unmarshal(body, &errResp)
			if err != nil {
				return err
			}
			pretty, jsonerr := json.MarshalIndent(&errResp, "", "  ")
			if jsonerr != nil {
				fmt.Printf("%s\n", body)
			} else {
				fmt.Printf("%s\n", pretty)
			}
			return &errResp
		}
		if err != nil {
			fmt.Printf("%s\n", body)
		} else {
			pretty, jsonerr := json.MarshalIndent(target, "", "  ")
			if jsonerr != nil {
				fmt.Printf("%s\n", body)
			} else {
				fmt.Printf("%s\n", pretty)
			}
		}
	} else {
		if res.StatusCode == 200 {
			err = json.NewDecoder(res.Body).Decode(target)
		} else {
			errResp := ECSResponseError{}
			err = json.NewDecoder(res.Body).Decode(&errResp)
			if err != nil {
				return err
			}
			return &errResp
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func ForAllRegionsDo(do func(region string) (err error)) (err error) {
	var regions ECSRegions
	var wg sync.WaitGroup
	var errs errors.Errors
	regions, _, err = ECS_INSTANCE.DescribeRegions()
	if err != nil {
		return
	}
	for _, region := range regions {
		wg.Add(1)
		go func(region string) {
			err := do(region)
			if err != nil {
				errs.Add(err.Error())
			}
			wg.Done()
		}(region.RegionID)
	}
	wg.Wait()
	if errs.HaveError() {
		err = errs.Errorify()
		return
	}
	return
}

type ECSInterface interface {
	PrintTable()
	Print()
}

func getFirstPart(input string) string {
	if strings.Contains(input, "@") {
		input = input[:strings.Index(input, "@")]
	}
	return input
}

func ForAllArgsDo(args []string, call func(arg string)) {
	for _, arg := range args {
		call(getFirstPart(arg))
	}
}

func Print(printable ECSInterface, others ...interface{}) {
	err := others[len(others)-1]
	if err == nil {
		if IsQuiet {
			printable.Print()
		} else {
			printable.PrintTable()
		}
	} else {
		exit(err)
	}
}

func PrintTable(
	fields []interface{},
	showFields bool,
	listLength int,
	filter func(i int) bool,
	getInfo func(i int) map[interface{}]interface{},
) {
	fieldsLen := len(fields)
	maxlengths := make([]interface{}, fieldsLen)
	for i, field := range fields {
		maxlengths[i] = len(field.(string))
	}
	var lines [][]interface{}
	for i := 0; i < listLength; i++ {
		if filter != nil && filter(i) != true {
			continue
		}
		info := getInfo(i)
		line := make([]interface{}, fieldsLen)
		for j, field := range fields {
			line[j] = info[field]
			l := len(line[j].(string))
			if l > maxlengths[j].(int) {
				maxlengths[j] = l
			}
		}
		lines = append(lines, line)
	}
	format := fmt.Sprintf(strings.TrimSpace(strings.Repeat("%%-%ds  ", fieldsLen))+"\n", maxlengths...)
	if showFields {
		fmt.Printf(format, fields...)
	}
	for _, line := range lines {
		fmt.Printf(format, line...)
	}
}

func printFlagsForCommand(c *cli.Context, name string) {
	var flags []cli.Flag
	for _, command := range c.App.Commands {
		if command.Name == name {
			flags = command.Flags
			break
		}
	}
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
}

func hintsForBashComplete(c *cli.Context, flagName *string) {
	if flagName == nil {
		return
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
			fmt.Printf("%s@%dCPU,%.6gGMem\n", _type.InstanceTypeId, _type.CpuCoreCount, _type.MemorySize)
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

func checkValuesForBashComplete(c *cli.Context) bool {
	bashCompletionFlag := "--" + cli.BashCompletionFlag.Name
	for _, flag := range c.Command.Flags {
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

func describeInstancesForBashComplete(filter func(instance ECSInstance) bool) func(c *cli.Context) {
	return func(c *cli.Context) {
		instances, _ := ECS_INSTANCE.DescribeInstances()
		for _, instance := range instances {
			if filter != nil && !filter(instance) {
				continue
			}
			fmt.Printf("%s@%s\n", instance.InstanceId, instance.InstanceName)
		}
	}
}

func ensureInstanceOfTheSameNameDoesNotExist(name string) {
	instances, _ := ECS_INSTANCE.DescribeInstances()
	for _, instance := range instances {
		if instance.InstanceName == name {
			exit("Instance of the same name already exists. Please choose a different name.")
		}
	}
}
