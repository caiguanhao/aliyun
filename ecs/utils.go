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
)

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

	if IsVerbose {
		fmt.Println(url)
	}

	res, err := http.Get(url)
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

func Print(printable ECSInterface, others ...interface{}) {
	err := others[len(others)-1]
	if err == nil {
		if IsQuiet {
			printable.Print()
		} else {
			printable.PrintTable()
		}
	} else {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func PrintTable(fields []interface{}, length int, process func(i int) []interface{}) {
	maxlengths := make([]interface{}, len(fields))
	for i, field := range fields {
		maxlengths[i] = len(field.(string))
	}
	var lines [][]interface{}
	for i := 0; i < length; i++ {
		line := process(i)
		if line == nil {
			continue
		}
		for j := 0; j < len(fields); j++ {
			l := len(line[j].(string))
			if l > maxlengths[j].(int) {
				maxlengths[j] = l
			}
		}
		lines = append(lines, line)
	}
	format := fmt.Sprintf(strings.TrimSpace(strings.Repeat("%%-%ds  ", len(fields)))+"\n", maxlengths...)
	fmt.Printf(format, fields...)
	for _, line := range lines {
		fmt.Printf(format, line...)
	}
}