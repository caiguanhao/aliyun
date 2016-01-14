package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type httpResponse http.Response

func (resp *httpResponse) GetResponse() (*string, error) {
	var response bytes.Buffer

	response.WriteString(fmt.Sprintf("%s %s\n", resp.Proto, resp.Status))
	for key, values := range resp.Header {
		for _, value := range values {
			response.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}

	body, respErr := ioutil.ReadAll(resp.Body)
	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()
	if len(body) > 0 {
		response.WriteByte('\n')
		response.Write(body)
	}

	ret := response.String()
	return &ret, nil
}

func md5hash(file []byte) []byte {
	md5sum := md5.New()
	md5sum.Write(file)
	return md5sum.Sum(nil)
}

func request(method, remotePath string, localFile []byte, localFileMD5 []byte) (*httpResponse, error) {
	req, reqErr := http.NewRequest(method, api+remotePath, bytes.NewReader(localFile))
	if reqErr != nil {
		return nil, reqErr
	}

	var md5sum string
	if localFileMD5 != nil {
		md5sum = base64.StdEncoding.EncodeToString(localFileMD5)
	} else {
		md5sum = base64.StdEncoding.EncodeToString(md5hash(localFile))
	}
	contentType := http.DetectContentType(localFile)
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	path := fmt.Sprintf("/%s%s", bucket, remotePath)
	msg := strings.Join([]string{method, md5sum, contentType, date, path}, "\n")
	mac := hmac.New(sha1.New, []byte(SECRET))
	mac.Write([]byte(msg))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	auth := fmt.Sprintf("OSS %s:%s", KEY, signature)

	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-MD5", md5sum)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Date", date)

	client := &http.Client{}
	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, respErr
	}
	response := httpResponse(*resp)
	return &response, nil
}

func getHeader(remotePath, headerName string) (*string, error) {
	resp, reqErr := request("HEAD", remotePath, nil, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	value := resp.Header.Get(headerName)
	return &value, nil
}

func upload(remotePath, localPath string, checkETag bool) (*string, error) {
	if dryrun {
		if verbose {
			fmt.Println(localPath, "->", api+remotePath)
		} else {
			fmt.Println(localPath, "->", remotePath)
		}
		return nil, nil
	}

	if verbose {
		fmt.Println(remotePath, "- added")
	}

	localFile, readErr := ioutil.ReadFile(localPath)
	if readErr != nil {
		return nil, readErr
	}

	var localFileMD5 []byte
	if checkETag {
		var localMD5, remoteMD5 string
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			localFileMD5 = md5hash(localFile)
			localMD5 = fmt.Sprintf("%x", localFileMD5)
			wg.Done()
		}()
		go func() {
			etag, headErr := getHeader(remotePath, "Etag")
			if headErr == nil {
				remoteMD5 = strings.ToLower(strings.Replace(*etag, "\"", "", -1))
			}
			wg.Done()
		}()
		wg.Wait()
		if localMD5 != "" && localMD5 == remoteMD5 {
			fmt.Println(remotePath, "- no changes, ignored")
			return nil, nil
		}
	}

	fmt.Println(remotePath, "- uploading")
	resp, reqErr := request("PUT", remotePath, localFile, localFileMD5)
	if reqErr != nil {
		return nil, reqErr
	}

	response, respErr := resp.GetResponse()
	if respErr != nil {
		return nil, respErr
	}
	if resp.StatusCode != 200 {
		fmt.Println(remotePath, "- fail -", resp.Status)
		return response, nil
	}
	totalBytes += len(localFile)
	fmt.Println(remotePath, "- done")

	return response, nil
}

type result struct {
	path     *string
	response *string
	err      error
}

func process(done <-chan struct{}, paths <-chan [2]string, c chan<- result) {
	for path := range paths {
		ret, err := upload(path[1], path[0], true)
		p := path[0]
		select {
		case c <- result{&p, ret, err}:
		case <-done:
			return
		}
	}
}

func getLastPartOfPath(input string) string {
	index := strings.LastIndex(input, "/")
	if index > -1 {
		return input[index+1:]
	}
	return input
}

func pathsForFile(src *string) [2]string {
	path := target
	if parentsPath {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		path += filepath.Dir(*src) + "/"
	}
	if len(source) > 1 {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		return [2]string{*src, path + getLastPartOfPath(*src)}
	}
	if strings.HasSuffix(path, "/") {
		return [2]string{*src, path + getLastPartOfPath(*src)}
	} else {
		return [2]string{*src, path}
	}
}

func pathsForDirectory(root, src *string) [2]string {
	if parentsPath {
		if strings.HasSuffix(target, "/") {
			return [2]string{*src, target + *src}
		}
		return [2]string{*src, target + "/" + *src}
	}
	if strings.HasSuffix(target, "/") {
		r := *root
		if strings.HasSuffix(r, "/") {
			r = strings.TrimSuffix(r, "/")
		}
		path, _ := filepath.Rel(filepath.Dir(r), *src)
		return [2]string{*src, target + path}
	} else {
		path, _ := filepath.Rel(*root, *src)
		return [2]string{*src, target + "/" + path}
	}
}

func walkFiles(done <-chan struct{}) (<-chan [2]string, <-chan error) {
	paths := make(chan [2]string)
	errc := make(chan error, len(source))
	go func() {
		defer close(paths)
		defer close(errc)
		for _, root := range source {
			fi, err := os.Stat(root)
			if err != nil {
				errc <- err
			} else if fi.Mode().IsRegular() {
				paths <- pathsForFile(&root)
			} else {
				errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.Mode().IsRegular() {
						return nil
					}
					select {
					case paths <- pathsForDirectory(&root, &path):
					case <-done:
						return errors.New("walk canceled")
					}
					return nil
				})
			}
		}
	}()
	return paths, errc
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

var source []string
var target string

var api string
var bucket string
var apiPrefix string
var concurrency int

var parentsPath bool

var dryrun bool
var verbose bool

var totalBytes int

func makeAPI() string {
	if strings.Count(apiPrefix, "%s") == 1 {
		api = fmt.Sprintf(apiPrefix, bucket)
	} else {
		api = apiPrefix
	}
	return api
}

func init() {
	flag.BoolVar(&dryrun, "d", false, "")
	flag.StringVar(&bucket, "b", DEFAULT_BUCKET, "")
	flag.StringVar(&apiPrefix, "z", DEFAULT_API_PREFIX, "")
	flag.IntVar(&concurrency, "c", 2, "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&parentsPath, "parents", false, "")
	flag.Usage = func() {
		fmt.Println("oss [OPTION] SOURCE ... TARGET")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("    -c <num>   Specify how many files to process concurrently, default is 2, max is 10")
		fmt.Println()
		fmt.Println("    -b <name>  Specify bucket other than:", DEFAULT_BUCKET)
		fmt.Println("    -z <url>   Specify API URL prefix other than:", DEFAULT_API_PREFIX)
		fmt.Println("       You can use custom domain or official URL like this:")
		fmt.Println("       {http, https}://%s.oss-cn-{beijing, hangzhou, hongkong, qingdao, shenzhen}{, -internal}.aliyuncs.com")
		fmt.Println("       Note: %s will be replaced with the bucket name if specified")
		fmt.Println()
		fmt.Println("    --parents  Use full source file name under TARGET")
		fmt.Println()
		fmt.Println("    -v  Be verbosive")
		fmt.Println("    -d  Dry-run. See list of files that will be transferred,")
		fmt.Println("        show full URL if -v is also set")
		fmt.Println()
		fmt.Println("Built with key ID", KEY, MADE)
		fmt.Println("API:", makeAPI())
		fmt.Println("Source: https://github.com/caiguanhao/aliyun")
	}
	flag.Parse()

	makeAPI()

	if concurrency < 1 || concurrency > 10 {
		fmt.Println("Warning: bad concurrency value:", concurrency, ". Fall back to 2.")
		concurrency = 2
	}
}

func main() {
	l := len(flag.Args())
	if l < 2 {
		fmt.Println("Error: Please specify source and target.")
		return
	}
	source = flag.Args()[0 : l-1]
	target = flag.Args()[l-1]

	target = regexp.MustCompile("/{2,}").ReplaceAllLiteralString(target, "/")
	if !strings.HasPrefix(target, "/") {
		target = "/" + target
	}

	timeStart := time.Now()

	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done)

	c := make(chan result)
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			process(done, paths, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	for r := range c {
		if verbose && r.response != nil {
			path := r.path
			ret := r.response
			fmt.Println(*path, "returned:")
			fmt.Println(*ret)
		}
	}

	if err := <-errc; err != nil {
		fmt.Println(err)
	}

	if totalBytes > 0 {
		since := time.Since(timeStart)
		speed := humanBytes(float64(totalBytes) / since.Seconds())
		used := regexp.MustCompile("(\\.[0-9]{3})[0-9]+").ReplaceAllString(since.String(), "$1")
		fmt.Printf(
			"transferred: %s (%d bytes)  time used: %s  avg. speed: %s\n",
			humanBytes(float64(totalBytes)), totalBytes, used, speed+"/s",
		)
	}
}
