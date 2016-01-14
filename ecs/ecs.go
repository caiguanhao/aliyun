package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

type ECS struct {
	KEY    string
	SECRET string
}

type ECSInterface interface {
	Print()
	PrintTable()
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
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"PageSize":         "50",
		"PageNumber":       "1",
	}
	for k, v := range queries {
		params[k] = v
	}
	query := buildQueryString(params)
	signature := sign(ecs.SECRET, urlEncode(query))
	url := fmt.Sprintf("http://ecs.aliyuncs.com/?%s&Signature=%s", query, urlEncode(signature))

	if opts.IsVerbose {
		fmt.Println(url)
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if opts.IsVerbose {
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

func (ecs *ECS) Do(task string, target *ECSInterface) bool {
	var err error
	switch task {
	case "list", "list-instances":
		if flag.Arg(1) != "" {
			*target, err = (&DescribeInstanceAttribute{}).Do(ecs)
		} else {
			*target, err = (&DescribeInstances{}).Do(ecs)
		}
	case "images", "list-images":
		*target, err = (&DescribeImages{}).Do(ecs)
	case "regions", "list-regions":
		*target, err = (&DescribeRegions{}).Do(ecs)
	case "types", "list-instance-types":
		*target, err = (&DescribeInstanceTypes{}).Do(ecs)
	case "groups", "list-security-groups":
		*target, err = (&DescribeSecurityGroups{}).Do(ecs)
	case "create", "create-instance":
		*target, err = (&CreateInstance{}).Do(ecs)
	case "allocate", "allocate-public-ip":
		*target, err = (&AllocatePublicIP{}).Do(ecs)
	case "start", "start-instance":
		*target, err = (&StartInstance{}).Do(ecs)
	case "stop", "stop-instance":
		*target, err = (&StopInstance{}).Do(ecs)
	case "restart", "restart-instance":
		*target, err = (&RestartInstance{}).Do(ecs)
	case "remove", "remove-instance":
		*target, err = (&RemoveInstance{}).Do(ecs)
	case "update", "update-instance":
		*target, err = (&ModifyInstanceAttribute{}).Do(ecs)
	case "hide", "hide-instance":
		*target, err = (&ModifyInstanceAttribute{}).HideInstance(ecs, true)
	case "unhide", "unhide-instance":
		*target, err = (&ModifyInstanceAttribute{}).HideInstance(ecs, false)
	default:
		return false
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return true
}
