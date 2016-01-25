package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var api string
var apiPrefix string
var bucket string
var remote string
var local string

func sign(date, uri string) string {
	msg := strings.Join([]string{"GET", "", "", date, fmt.Sprintf("/%s%s", bucket, uri)}, "\n")
	mac := hmac.New(sha1.New, []byte(SECRET))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func get(remotePath string) (*http.Response, error) {
	req, reqErr := http.NewRequest("GET", api+remotePath, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT") // don't use time.RFC1123
	signature := sign(date, remotePath)
	auth := fmt.Sprintf("OSS %s:%s", KEY, signature)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Date", date)
	client := &http.Client{}
	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, respErr
	}
	return resp, nil
}

func makeAPI() string {
	if strings.Count(apiPrefix, "%s") == 1 {
		api = fmt.Sprintf(apiPrefix, bucket)
	} else {
		api = apiPrefix
	}
	return api
}

func download() (int64, error) {
	to := local
	info, err := os.Stat(to)
	if err != nil {
		return 0, err
	}

	if info.IsDir() {
		to = path.Join(to, path.Base(remote))
	}

	resp, err := get(remote)
	defer resp.Body.Close()
	if err != nil {
		return 0, err
	}

	file, err := os.Create(to)
	defer file.Close()
	if err != nil {
		return 0, err
	}

	fmt.Printf("Downloading %s to %s ...\n", remote, to)

	n, err := io.Copy(file, resp.Body)

	if err != nil {
		return 0, err
	}

	return n, nil
}

func fmtFloat(float float64, suffix string) string {
	return strings.TrimSuffix(strings.TrimRight(fmt.Sprintf("%.3f", float), "0"), ".") + suffix
}

func humanBytes(bytes float64) string {
	const TB = 1 << 40
	const GB = 1 << 30
	const MB = 1 << 20
	const KB = 1 << 10
	abs := bytes
	if bytes < 0 {
		abs = bytes * -1
	}
	if abs >= TB {
		return fmtFloat(bytes/TB, " TB")
	}
	if abs >= GB {
		return fmtFloat(bytes/GB, " GB")
	}
	if abs >= MB {
		return fmtFloat(bytes/MB, " MB")
	}
	if abs >= KB {
		return fmtFloat(bytes/KB, " KB")
	}
	return fmt.Sprintf("%.0f bytes", bytes)
}

var curlScript bool

func init() {
	flag.BoolVar(&curlScript, "curl", false, "")
	flag.Usage = func() {
		fmt.Println("oss-get [--curl] REMOTE-FILE LOCAL-FILE")
		fmt.Println()
		fmt.Println("  --curl   generate curl script")
	}
	flag.Parse()
}

func main() {
	timeStart := time.Now()

	l := len(flag.Args())
	if l != 2 {
		fmt.Println("Error: Please specify remote file and local file location.")
		return
	}
	remote = flag.Args()[0]
	local = flag.Args()[1]
	apiPrefix = DEFAULT_API_PREFIX
	bucket = DEFAULT_BUCKET
	api = makeAPI()
	if curlScript {
		uri := strings.TrimLeft(remote, "/")
		date := time.Now().Unix() + 3600
		signature := sign(fmt.Sprintf("%d", date), remote)
		fmt.Printf("curl -#o '%s' '%s/%s?OSSAccessKeyId=%s&Expires=%d&Signature=%s'\n",
			path.Base(remote), api, url.QueryEscape(uri), KEY, date, url.QueryEscape(signature))
	} else {
		downloaded, err := download()
		if err != nil {
			fmt.Println(err)
			return
		}
		if downloaded > 0 {
			since := time.Since(timeStart)
			speed := humanBytes(float64(downloaded) / since.Seconds())
			used := regexp.MustCompile("(\\.[0-9]{3})[0-9]+").ReplaceAllString(since.String(), "$1")
			fmt.Printf(
				"transferred: %s (%d bytes)  time used: %s  avg. speed: %s\n",
				humanBytes(float64(downloaded)), downloaded, used, speed+"/s",
			)
		}
	}
}
