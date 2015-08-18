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
	mac := hmac.New(sha1.New, SECRET)
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

func removeFirstPartOfPath(input string) string {
	index := strings.Index(input, "/")
	if index > -1 {
		return input[index:]
	}
	return input
}

func pathsForFile(src *string) [2]string {
	if strings.HasSuffix(target, "/") {
		return [2]string{*src, target + getLastPartOfPath(*src)}
	} else {
		return [2]string{*src, target}
	}
}

func pathsForDirectory(src *string) [2]string {
	if strings.HasSuffix(target, "/") {
		return [2]string{*src, target + *src}
	} else {
		return [2]string{*src, target + removeFirstPartOfPath(*src)}
	}
}

func walkFiles(done <-chan struct{}, root string) (<-chan [2]string, <-chan error) {
	paths := make(chan [2]string)
	errc := make(chan error, 1)
	go func() {
		defer close(paths)
		fi, err := os.Stat(root)
		if err != nil {
			errc <- err
		} else if fi.Mode().IsRegular() {
			defer close(errc)
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
				case paths <- pathsForDirectory(&path):
				case <-done:
					return errors.New("walk canceled")
				}
				return nil
			})
		}
	}()
	return paths, errc
}

var source, target string

var api string
var bucket string
var apiPrefix string
var concurrency int

var dryrun bool
var verbose bool

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
	flag.StringVar(&bucket, "b", string(DEFAULT_BUCKET), "")
	flag.StringVar(&apiPrefix, "z", string(DEFAULT_API_PREFIX), "")
	flag.IntVar(&concurrency, "c", 2, "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.Usage = func() {
		fmt.Println("oss [OPTION] SOURCE TARGET")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("    -c <num>   Specify how many files to process concurrently, default is 2, max is 10")
		fmt.Println()
		fmt.Println("    -b <name>  Specify bucket other than:", string(DEFAULT_BUCKET))
		fmt.Println("    -z <url>   Specify API URL prefix other than:", string(DEFAULT_API_PREFIX))
		fmt.Println("       You can use custom domain or official URL like this:")
		fmt.Println("       {http, https}://%s.oss-cn-{beijing, hangzhou, hongkong, qingdao, shenzhen}{, -internal}.aliyuncs.com")
		fmt.Println("       Note: %s will be replaced with the bucket name if specified")
		fmt.Println()
		fmt.Println("    -v  Be verbosive")
		fmt.Println("    -d  Dry-run. See list of files that will be transferred,")
		fmt.Println("        show full URL if -v is also set")
		fmt.Println()
		fmt.Println("Built with key ID:", string(KEY))
		fmt.Println("API:", makeAPI())
		fmt.Println("Source: https://github.com/caiguanhao/oss")
	}
	flag.Parse()

	source = flag.Arg(0)
	target = flag.Arg(1)

	makeAPI()

	if concurrency < 1 || concurrency > 10 {
		fmt.Println("Warning: bad concurrency value:", concurrency, ". Fall back to 2.")
		concurrency = 2
	}
}

func main() {
	if source == "" {
		fmt.Println("Error: Please specify source.")
		return
	}
	if target == "" {
		fmt.Println("Error: Please specify target.")
		return
	}
	target = regexp.MustCompile("/{2,}").ReplaceAllLiteralString(target, "/")
	if !strings.HasPrefix(target, "/") {
		target = "/" + target
	}

	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, source)

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
}
