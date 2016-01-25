package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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

func download(to string) (written int64, err error) {
	var resp *http.Response
	resp, err = get(remote)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			err = errors.New(string(body))
		}
		return
	}

	if useSTDIN {
		written, err = io.Copy(os.Stdout, resp.Body)
		return
	}

	var file *os.File
	file, err = os.Create(to)
	defer file.Close()
	if err == nil {
		fmt.Printf("Downloading %s to %s ...\n", remote, to)
		written, err = io.Copy(file, resp.Body)
	}
	return
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
var useSTDIN bool

func init() {
	flag.BoolVar(&curlScript, "curl", false, "")
	flag.Usage = func() {
		fmt.Println("oss-get [--curl] REMOTE-FILE [LOCAL-FILE]")
		fmt.Println()
		fmt.Println("    --curl    generate curl script")
	}
	flag.Parse()
}

func main() {
	timeStart := time.Now()

	l := flag.NArg()
	if l < 1 {
		fmt.Fprintln(os.Stderr, "Error: Please specify remote file location.")
		return
	}
	if l > 2 {
		fmt.Fprintln(os.Stderr, "Error: Please specify one remote file and one local file location.")
		return
	}
	remote = flag.Arg(0)
	if l == 1 {
		useSTDIN = true
	} else if l == 2 {
		local = flag.Arg(1)
	}
	apiPrefix = DEFAULT_API_PREFIX
	bucket = DEFAULT_BUCKET
	api = makeAPI()
	if !useSTDIN {
		info, err := os.Stat(local)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if info.IsDir() {
			local = path.Join(local, path.Base(remote))
		}
	}
	if curlScript {
		uri := strings.TrimLeft(remote, "/")
		date := time.Now().Unix() + 3600
		signature := sign(fmt.Sprintf("%d", date), remote)
		var output string
		if !useSTDIN {
			output = fmt.Sprintf("-o '%s' ", local)
		}
		fmt.Printf("curl %s'%s/%s?OSSAccessKeyId=%s&Expires=%d&Signature=%s'\n",
			output, api, url.QueryEscape(uri), KEY, date, url.QueryEscape(signature))
	} else {
		downloaded, err := download(local)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if downloaded > 0 {
			since := time.Since(timeStart)
			speed := humanBytes(float64(downloaded) / since.Seconds())
			used := regexp.MustCompile("(\\.[0-9]{3})[0-9]+").ReplaceAllString(since.String(), "$1")
			fmt.Fprintf(os.Stderr, "transferred: %s (%d bytes)  time used: %s  avg. speed: %s/s\n",
				humanBytes(float64(downloaded)), downloaded, used, speed)
		}
	}
}
