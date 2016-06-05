package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
)

type OSSFileList struct {
	Name        string
	Prefix      string
	Marker      string
	MaxKeys     int
	Delimiter   string
	IsTruncated bool
	NextMarker  string
	Files       []OSSFile      `xml:"Contents"`
	Directories []OSSDirectory `xml:"CommonPrefixes"`
}

type OSSDirectory struct {
	Name string `xml:"Prefix"`
}

type OSSFile struct {
	Name         string `xml:"Key"`
	LastModified string
	ETag         string
	Size         int64
	err          error
}

var OSS_LIST cli.Command = cli.Command{
	Name:      "list",
	Aliases:   []string{"ls", "l"},
	Usage:     "show list of files on remote OSS",
	ArgsUsage: "REMOTE ...",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "recursive, r",
			Usage: "recursively list subdirectories",
		},
	},
	Action: func(c *cli.Context) {
		remotes := parseArgsForOSSList([]string(c.Args()))
		multiple := len(remotes) > 1
		for i, remote := range remotes {
			if multiple {
				if i > 0 {
					fmt.Println()
				}
				fmt.Printf("%s:\n", remote)
			}
			remoteFiles, remoteDirs, err := getOSSFileList(remote, c.Bool("recursive"))
			if err != nil {
				die(err)
			}
			hasSuffix := strings.HasSuffix(remote, "/")
			for _, dir := range remoteDirs {
				if hasSuffix {
					fmt.Println(dir.Name[len(remote):])
				} else {
					fmt.Println(dir.Name)
				}
			}
			for _, file := range remoteFiles {
				if hasSuffix {
					fmt.Println(file.Name[len(remote):])
				} else {
					fmt.Println(file.Name)
				}
			}
		}
	},
}

func parseArgsForOSSList(args []string) (remotes []string) {
	if len(args) < 1 {
		remotes = append(remotes, "")
		return
	}
	for _, arg := range args {
		hasSuffix := strings.HasSuffix(arg, "/")
		arg = strings.TrimLeft(filepath.Clean(arg), "/")
		if hasSuffix {
			arg += "/"
		}
		remotes = append(remotes, arg)
	}
	return
}

func getOSSFileListWithMarker(prefix string, marker *string, files *[]OSSFile, dirs *[]OSSDirectory, recursive bool) (err error) {
	queryString := "?max-keys=1000"
	if !recursive {
		queryString += "&delimiter=/"
	}
	queryString += "&prefix=" + url.QueryEscape(prefix)
	if marker != nil {
		queryString += "&marker=" + url.QueryEscape(*marker)
	}
	var resp *http.Response
	resp, err = sendGetRequest("/" + queryString)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = checkOSSResponse(resp)
	if err != nil {
		return
	}
	var list OSSFileList
	if err = xml.NewDecoder(resp.Body).Decode(&list); err != nil {
		return
	}
	*files = append(*files, list.Files...)
	*dirs = append(*dirs, list.Directories...)
	if verbose {
		fmt.Fprintf(os.Stderr, "Remote: received %d file names (out of %d) ...\n", len(list.Files), len(*files))
	}
	if list.IsTruncated {
		err = getOSSFileListWithMarker(prefix, &list.NextMarker, files, dirs, recursive)
	}
	return
}

func getOSSFileList(prefix string, recursive bool) (files []OSSFile, dirs []OSSDirectory, err error) {
	err = getOSSFileListWithMarker(prefix, nil, &files, &dirs, recursive)
	return
}
