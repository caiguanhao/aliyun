package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caiguanhao/gotogether"
	"github.com/codegangsta/cli"
)

var OSS_DIFF cli.Command = cli.Command{
	Name:      "diff",
	Aliases:   []string{},
	Usage:     "show different files on local and remote OSS",
	ArgsUsage: "LOCAL REMOTE",
	Description: `Status code: 0 - local and remote are identical
                1 - local has different files
                2 - remote has different files
                3 - both local and remote have different files`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "reverse, r",
			Usage: "print LOCAL file paths to stderr, REMOTE to stdout",
		},
		cli.BoolFlag{
			Name:  "md5",
			Usage: "also verify MD5 checksum besides file name and size",
		},
	},
	Action: func(c *cli.Context) {
		totalErrors := 0
		timeStart := time.Now()
		checkMD5 := c.Bool("md5")

		local, remote, localStrPos, remoteStrPos := parseArgsForOSSDiff([]string(c.Args()))

		var remoteFiles, localFiles []OSSFile
		var localTimeUsed, remoteTimeUsed time.Duration

		gotogether.Parallel{
			func() {
				timeStart := time.Now()
				var err error
				remoteFiles, _, err = getOSSFileList(remote, true)
				if err != nil {
					die(err)
				}
				remoteTimeUsed = time.Since(timeStart)
			},
			func() {
				timeStart := time.Now()
				if checkMD5 && verbose {
					fmt.Fprintf(os.Stderr, "MD5 checksum verification using up to %d goroutines\n", concurrency)
				}
				var errs int
				localFiles, errs = getLocalFileList(local)
				totalErrors += errs
				localTimeUsed = time.Since(timeStart)
			},
		}.Run()

		localFilesLength, remoteFilesLength := len(localFiles), len(remoteFiles)

		stdout, stderr := os.Stdout, os.Stderr
		if c.Bool("reverse") {
			stdout, stderr = stderr, stdout
		}

		localOnly, remoteOnly := getDiff(localFiles, remoteFiles, localStrPos, remoteStrPos, checkMD5)
		localOnlyLen, remoteOnlyLen := len(localOnly), len(remoteOnly)

		if verbose {
			fmt.Fprintf(os.Stderr, "Local: %d files, %d different on local\n", localFilesLength, localOnlyLen)
		}
		for i := range localOnly {
			fmt.Fprintln(stdout, localOnly[i].Name)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Remote: %d files, %d different on remote\n", remoteFilesLength, remoteOnlyLen)
		}
		for i := range remoteOnly {
			fmt.Fprintln(stderr, remoteOnly[i].Name)
		}

		retCode := 4
		if totalErrors == 0 {
			if localOnlyLen == 0 && remoteOnlyLen == 0 {
				if verbose {
					fmt.Fprintln(os.Stderr, "Local and remote are identical")
				}
				retCode = 0
			} else if localOnlyLen > 0 && remoteOnlyLen == 0 {
				retCode = 1
			} else if localOnlyLen == 0 && remoteOnlyLen > 0 {
				retCode = 2
			} else if localOnlyLen > 0 && remoteOnlyLen > 0 {
				retCode = 3
			}
			if verbose {
				fmt.Fprintf(
					os.Stderr,
					"Time used:  %s local,  %s remote,  %s total\n",
					localTimeUsed.String(),
					remoteTimeUsed.String(),
					time.Since(timeStart).String(),
				)
			}
		} else {
			fmt.Fprintf(os.Stderr, "%d error(s):\n", totalErrors)
			for _, f := range localFiles {
				if f.err != nil {
					fmt.Fprintln(os.Stderr, f.err)
				}
			}
		}
		os.Exit(retCode)
	},
}

func parseArgsForOSSDiff(args []string) (local, remote string, localStrPos, remoteStrPos int) {
	if len(args) < 2 {
		die("Error: Please specify local and remote.")
	}

	if len(args) > 2 {
		die("Error: Please specify one local and remote.")
	}

	var err error
	local, err = filepath.Abs(filepath.Clean(args[0]))
	if err != nil {
		die(err)
	}
	remote = strings.Trim(filepath.Clean(args[1]), "/")

	info, err := os.Stat(local)
	if err != nil {
		die(err)
	}

	if info.Mode().IsRegular() {
		localStrPos, remoteStrPos = len(filepath.Dir(local)), len(filepath.Dir(remote))
	} else {
		localStrPos, remoteStrPos = len(local), len(remote)
	}
	return
}

func getLocalFileList(local string) (files []OSSFile, errs int) {
	gotogether.Queue{
		Concurrency: concurrency,
		AddJob: func(jobs *chan interface{}) {
			err := filepath.Walk(local, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.Mode().IsRegular() {
					return nil
				}
				*jobs <- OSSFile{
					Name: path,
					Size: info.Size(),
				}
				return nil
			})
			if err != nil {
				debug(err)
				errs++
			}
		},
		DoJob: func(job *interface{}) {
			file := (*job).(OSSFile)
			data, err := ioutil.ReadFile(file.Name)
			if err != nil {
				file.err = err
				errs++
			} else {
				file.ETag = fmt.Sprintf("\"%X\"", md5.Sum(data))
			}
			files = append(files, file)
		},
	}.Run()
	return
}

func getDiff(localFiles, remoteFiles []OSSFile, localStrPos, remoteStrPos int, checkMD5 bool) (localOnly, remoteOnly []OSSFile) {
	gotogether.Parallel{
		func() {
			localOnly = diff(localFiles, remoteFiles, localStrPos, remoteStrPos, checkMD5)
		},
		func() {
			remoteOnly = diff(remoteFiles, localFiles, remoteStrPos, localStrPos, checkMD5)
		},
	}.Run()
	return
}

func diff(left, right []OSSFile, leftStrPos, rightStrPos int, checkMD5 bool) (ret []OSSFile) {
	leftLen, rightLen, retLen := len(left), len(right), 0
	for i := 0; i < leftLen; i++ {
		j, k := 0, 0
		for j < rightLen {
			if right[j].Name[rightStrPos:] == left[i].Name[leftStrPos:] &&
				right[j].Size == left[i].Size &&
				(!checkMD5 || right[j].ETag == left[i].ETag) {
				break
			}
			j++
		}
		for k < retLen {
			if ret[k].Name[leftStrPos:] == left[i].Name[leftStrPos:] &&
				ret[k].Size == left[i].Size {
				break
			}
			k++
		}
		if j == rightLen && k == retLen {
			ret = append(ret, left[i])
			retLen++
		}
	}
	return
}
