package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/caiguanhao/gotogether"
	"github.com/codegangsta/cli"
)

var OSS_DOWNLOAD cli.Command = cli.Command{
	Name:      "download",
	Aliases:   []string{"down", "dl", "d", "get"},
	Usage:     "get remote OSS files to local",
	ArgsUsage: "REMOTE ... [LOCAL]",
	Description: `If only one REMOTE but LOCAL is none, it will send the remote file to STDOUT.
   If only one REMOTE and one LOCAL, LOCAL can be a file or directory path.
   If more than one REMOTE, LOCAL must be a directory path.
   It will generated curl command line script if --dry-run provided.`,
	Action: func(c *cli.Context) {
		stat := (&Stat{}).Begin()
		totalErrors := 0

		remotes, locals := parseArgsForOSSDownload([]string(c.Args()))

		toStdout := len(c.Args()) == 1

		if dryrun {
			for index, remote := range remotes {
				segments := []string{"curl"}
				if !toStdout {
					segments = append(segments, fmt.Sprintf("-o %s", strconv.Quote(locals[index])))
				}
				segments = append(segments, strconv.Quote(getDownloadUrl(remote, 3600)))
				fmt.Println(strings.Join(segments, " "))
			}
			return
		}

		defer func() {
			if summary := stat.String(); summary != nil && !toStdout {
				debug(*summary)
			}
			if totalErrors > 0 {
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}()

		if toStdout {
			ret, err := remoteFileToStdOut(remotes[0])
			if err != nil {
				die(err)
			}
			stat.Add(ret)
			return
		}

		gotogether.Queue{
			Concurrency: concurrency,
			AddJob: func(jobs *chan interface{}) {
				for index, remote := range remotes {
					paths := []string{remote, locals[index]}
					*jobs <- paths
				}
			},
			DoJob: func(job *interface{}) {
				paths := (*job).([]string)
				size, err := remoteFileToLocalFile(paths[0], paths[1])
				if err == nil {
					stat.Add(size)
				} else {
					debug(paths[0], err)
					totalErrors++
				}
			},
		}.Run()
	},
}

func parseArgsForOSSDownload(args []string) (remotes, locals []string) {
	l := len(args)

	if l == 0 {
		die("Error: Please specify remote file location.")
	}

	defer func() {
		for i, remote := range remotes {
			remote = filepath.Clean(remote)
			if !strings.HasPrefix(remote, "/") {
				remote = "/" + remote
			}
			remotes[i] = remote
		}
	}()

	if l == 1 {
		remotes = args
		return
	}

	remotes = args[0 : l-1]
	local := filepath.Clean(args[l-1])
	info, err := os.Stat(local)
	if l == 2 {
		if err != nil && !os.IsNotExist(err) {
			die(err)
		}
		if info != nil && info.IsDir() {
			locals = []string{filepath.Join(local, filepath.Base(remotes[0]))}
		} else {
			locals = []string{local}
		}
	} else {
		if (err != nil && os.IsNotExist(err)) || (err == nil && !info.IsDir()) {
			die("Error: LOCAL must be an existing directory.")
		} else if err != nil {
			die(err)
		}
		for _, r := range remotes {
			locals = append(locals, filepath.Join(local, filepath.Base(r)))
		}
	}
	return
}

func getDownloadUrl(remote string, secondsFromNow int64) string {
	date := time.Now().Unix() + secondsFromNow
	signature := (&Signature{Date: fmt.Sprintf("%d", date), URI: remote}).Get()
	url := fmt.Sprintf("%s/%s?OSSAccessKeyId=%s&Expires=%d&Signature=%s",
		api, url.QueryEscape(strings.TrimLeft(remote, "/")), accessKey, date, url.QueryEscape(signature))
	return url
}

func getRemoteFile(remote string) (resp *http.Response, err error) {
	resp, err = sendGetRequest(remote)
	if err != nil {
		return
	}
	err = checkOSSResponse(resp)
	return
}

func remoteFileToStdOut(remote string) (written int64, err error) {
	var resp *http.Response
	resp, err = getRemoteFile(remote)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if sendFileToStdOutConfirm(resp.Header.Get("content-type")) {
		written, err = io.Copy(os.Stdout, resp.Body)
	}
	return
}

func remoteFileToLocalFile(remote, local string) (written int64, err error) {
	var resp *http.Response
	resp, err = getRemoteFile(remote)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var file *os.File
	file, err = os.Create(local)
	defer file.Close()
	if err == nil {
		debugf("Downloading %s to %s ...\n", remote, local)
		written, err = io.Copy(file, resp.Body)
		debugf("Downloaded %s to %s ...\n", remote, local)
	}
	return
}

func sendFileToStdOutConfirm(contentType string) bool {
	if !strings.Contains(contentType, "text") && isTerminal(int(os.Stdout.Fd())) {
		fmt.Fprint(os.Stderr, "REMOTE may be a binary file.  See it anyway? [y] ")
		var answer string
		_, err := fmt.Scanln(&answer)
		if err != nil {
			return false
		}
		if strings.IndexAny(answer, "Yy") == 0 {
			return true
		}
		return false
	}
	return true
}
