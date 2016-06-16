package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/caiguanhao/gotogether"
	"github.com/codegangsta/cli"
)

var OSS_UPLOAD cli.Command = cli.Command{
	Name:        "upload",
	Aliases:     []string{"up", "u", "put"},
	Usage:       "upload local files to remote OSS",
	ArgsUsage:   "[LOCAL ...] REMOTE",
	Description: `If only one REMOTE but LOCAL is none, it will read contents of STDIN and upload them to REMOTE.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "parents, p",
			Usage: "use full SOURCE file name under TARGET",
		},
	},
	Action: func(c *cli.Context) {
		stat := (&Stat{}).Begin()
		totalErrors := 0

		locals, remote := parseArgsForOSSUpload([]string(c.Args()))

		fromStdin := len(c.Args()) == 1

		defer func() {
			if summary := stat.String(); summary != nil {
				debug(*summary)
			}
			if totalErrors > 0 {
				debugf("%d error(s) occurred during uploading.\n", totalErrors)
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}()

		if fromStdin {
			if dryrun {
				die()
			}
			if strings.HasSuffix(remote, "/") {
				die("Error: REMOTE must be a file path.")
			}
			readFileFromStdInPrompt()
			ret, err := localReaderToRemoteFile(os.Stdin, remote)
			if err != nil {
				die(remote+":", err)
			}
			stat.Add(ret)
			return
		}

		parentsPath := c.Bool("parents")

		gotogether.Queue{
			Concurrency: concurrency,
			AddJob: func(jobs *chan interface{}) {
				localsMoreThanOne := len(locals) > 1
				for _, local := range locals {
					fi, err := os.Stat(local)
					if err != nil {
						debug(err)
						totalErrors++
					} else if fi.Mode().IsRegular() {
						*jobs <- []string{local, localPathToRemotePath(local, remote, localsMoreThanOne, parentsPath)}
					} else {
						err := filepath.Walk(local, func(path string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if !info.Mode().IsRegular() {
								return nil
							}
							*jobs <- []string{path, localDirectoryToRemotePath(local, path, remote, parentsPath)}
							return nil
						})
						if err != nil {
							debug(err)
							totalErrors++
						}
					}
				}
			},
			DoJob: func(job *interface{}) {
				paths := (*job).([]string)
				if dryrun {
					fmt.Println(paths[0], " -> ", paths[1])
					return
				}
				size, err := localFileToRemoteFile(paths[0], paths[1])
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

func parseArgsForOSSUpload(args []string) (locals []string, remote string) {
	l := len(args)

	if l == 0 {
		die("Error: Please specify remote file location.")
	}

	defer func() {
		suffix := false
		if strings.HasSuffix(remote, "/") {
			suffix = true
		}
		remote = filepath.Clean(remote)
		if !strings.HasPrefix(remote, "/") {
			remote = "/" + remote
		}
		if !strings.HasSuffix(remote, "/") && suffix {
			remote += "/"
		}
	}()

	if l == 1 {
		remote = args[0]
		return
	}

	for _, arg := range args[0 : l-1] {
		locals = append(locals, filepath.Clean(arg))
	}
	remote = args[l-1]
	return
}

func getHeader(remote, key string) (value string, err error) {
	var resp *http.Response
	resp, err = sendRequest("HEAD", remote, nil, nil)
	if err != nil {
		return
	}
	value = resp.Header.Get(key)
	return
}

func localFileToRemoteFile(local, remote string) (size int64, err error) {
	var localFile []byte
	localFile, err = ioutil.ReadFile(local)
	if err != nil {
		return
	}
	size, err = localBytesToRemoteFile(localFile, remote)
	return
}

func localReaderToRemoteFile(r io.Reader, remote string) (size int64, err error) {
	var localFile []byte
	localFile, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	size, err = localBytesToRemoteFile(localFile, remote)
	return
}

func localBytesToRemoteFile(localFile []byte, remote string) (size int64, err error) {
	var localFileMD5 []byte
	var localMD5, remoteMD5 string
	gotogether.Parallel{
		func() {
			localFileMD5 = md5hash(localFile)
			localMD5 = fmt.Sprintf("%x", localFileMD5)
		},
		func() {
			etag, err := getHeader(remote, "Etag")
			if err == nil {
				remoteMD5 = strings.ToLower(strings.Replace(etag, "\"", "", -1))
			}
		},
	}.Run()
	if localMD5 != "" && localMD5 == remoteMD5 {
		fmt.Println(remote+":", "no changes, ignored")
		return
	}
	fmt.Println(remote+":", "uploading")
	var resp *http.Response
	resp, err = sendRequest("PUT", remote, localFile, localFileMD5)
	if err != nil {
		return
	}

	err = checkOSSResponse(resp)
	if err != nil {
		return
	}

	size = int64(len(localFile))
	fmt.Println(remote+":", "done")
	return
}

func getLastPartOfPath(input string) string {
	index := strings.LastIndex(input, "/")
	if index > -1 {
		return input[index+1:]
	}
	return input
}

func localPathToRemotePath(local, remote string, localsMoreThanOne, parentsPath bool) (ret string) {
	defer func() {
		ret = filepath.Clean(ret)
	}()
	path := remote
	if parentsPath {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		path += filepath.Dir(local) + "/"
	}
	if localsMoreThanOne {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		ret = path + getLastPartOfPath(local)
		return
	}
	if strings.HasSuffix(path, "/") {
		ret = path + getLastPartOfPath(local)
		return
	} else {
		ret = path
		return
	}
}

func localDirectoryToRemotePath(directory, local, remote string, parentsPath bool) (ret string) {
	defer func() {
		ret = filepath.Clean(ret)
	}()
	if parentsPath {
		if strings.HasSuffix(remote, "/") {
			ret = remote + local
			return
		}
		ret = remote + "/" + local
		return
	}
	if strings.HasSuffix(remote, "/") {
		if strings.HasSuffix(directory, "/") {
			directory = strings.TrimSuffix(directory, "/")
		}
		path, _ := filepath.Rel(filepath.Dir(directory), local)
		ret = remote + path
		return
	} else {
		path, _ := filepath.Rel(directory, local)
		ret = remote + "/" + path
		return
	}
}

func readFileFromStdInPrompt() {
	if isTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintln(os.Stderr, "Press CTRL-D to send, CTRL-C to abort.")
	}
}
