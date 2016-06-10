package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func debug(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func debugf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func die(msg ...interface{}) {
	if len(msg) > 0 {
		fmt.Fprintln(os.Stderr, msg...)
	}
	os.Exit(1)
}

type Signature struct {
	Method, MD5Sum, ContentType, Date, URI string
}

func (signature *Signature) Get() string {
	if signature.Method == "" {
		signature.Method = "GET"
	}
	if signature.Date == "" {
		signature.Date = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT") // don't use time.RFC1123
	}
	msg := strings.Join([]string{
		signature.Method,
		signature.MD5Sum,
		signature.ContentType,
		signature.Date,
		fmt.Sprintf("/%s%s", bucket, signature.URI),
	}, "\n")
	mac := hmac.New(sha1.New, []byte(accessSecret))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (signature *Signature) SetRequest(req *http.Request) {
	sign := signature.Get()
	req.Header.Set("Authorization", fmt.Sprintf("OSS %s:%s", accessKey, sign))
	if signature.MD5Sum != "" {
		req.Header.Set("Content-MD5", signature.MD5Sum)
	}
	if signature.ContentType != "" {
		req.Header.Set("Content-Type", signature.ContentType)
	}
	req.Header.Set("Date", signature.Date)
}

type Stat struct {
	timeStart time.Time
	total     int64
}

func (stat *Stat) Begin() *Stat {
	stat.timeStart = time.Now()
	return stat
}

func (stat *Stat) Add(n int64) {
	stat.total += n
}

func (stat *Stat) String() *string {
	if stat.total == 0 {
		return nil
	}
	since := time.Since(stat.timeStart)
	speed := humanBytes(float64(stat.total) / since.Seconds())
	used := regexp.MustCompile("(\\.[0-9]{3})[0-9]+").ReplaceAllString(since.String(), "$1")
	summary := fmt.Sprintf("transferred: %s (%d bytes)  time used: %s  avg. speed: %s/s",
		humanBytes(float64(stat.total)), stat.total, used, speed)
	return &summary
}

func md5hash(file []byte) []byte {
	md5sum := md5.New()
	md5sum.Write(file)
	return md5sum.Sum(nil)
}

func sendRequest(method, remote string, localFile []byte, localFileMD5 []byte) (resp *http.Response, err error) {
	var req *http.Request
	req, err = http.NewRequest(method, api+remote, bytes.NewReader(localFile))
	if err != nil {
		return
	}
	var md5sum, contentType, remoteNoQS string
	if i := strings.Index(remote, "?"); i > -1 {
		remoteNoQS = remote[:i]
	} else {
		remoteNoQS = remote
	}
	if localFile != nil {
		if localFileMD5 != nil {
			md5sum = base64.StdEncoding.EncodeToString(localFileMD5)
		} else {
			md5sum = base64.StdEncoding.EncodeToString(md5hash(localFile))
		}
		contentType = http.DetectContentType(localFile)
	}
	(&Signature{Method: method, MD5Sum: md5sum, ContentType: contentType, URI: remoteNoQS}).SetRequest(req)
	client := &http.Client{}
	resp, err = client.Do(req)
	return
}

func sendGetRequest(remote string) (resp *http.Response, err error) {
	resp, err = sendRequest("GET", remote, nil, nil)
	return
}

func checkOSSResponse(resp *http.Response) (err error) {
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			errResp := OSSResponseError{}
			err = xml.Unmarshal(body, &errResp)
			if err == nil && len(errResp.Message) > 0 {
				err = errors.New(errResp.Message)
			} else {
				err = errors.New(strings.TrimSpace(string(body)))
			}
		}
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

// https://github.com/golang/crypto/blob/master/ssh/terminal/util.go

var ioctlReadTermios uintptr
var NUM_CPU int

func isTerminalInit() {
	if runtime.GOOS == "darwin" {
		ioctlReadTermios = 0x40487413
	} else {
		ioctlReadTermios = 0x5401
	}
}

func isTerminal(fd int) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

func numOfCPUInit() {
	NUM_CPU = runtime.NumCPU()
}
