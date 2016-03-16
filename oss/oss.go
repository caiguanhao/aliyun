package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/caiguanhao/aliyun/vendor/cli"
)

var api string
var bucket string
var prefix string
var concurrency int
var dryrun bool
var verbose bool

func init() {
	isTerminalInit()
	numOfCPUInit()
}

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Version = "1.0.0"
	app.Usage = "control Aliyun OSS"
	app.HideHelp = true
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		OSS_UPLOAD,
		OSS_DOWNLOAD,
		OSS_DIFF,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "bucket, b",
			Usage:       "bucket name",
			Destination: &bucket,
		},
		cli.StringFlag{
			Name:        "prefix, p",
			Usage:       "API prefix",
			Destination: &prefix,
		},
		cli.IntFlag{
			Name:        "concurrency, c",
			Value:       NUM_CPU,
			Usage:       fmt.Sprintf("job concurrency, defaults to number of CPU (%d)", NUM_CPU),
			Destination: &concurrency,
		},
		cli.BoolFlag{
			Name:        "dry-run, D",
			Usage:       "do not actually run",
			Destination: &dryrun,
		},
		cli.BoolFlag{
			Name:        "verbose, V",
			Usage:       "show more info",
			Destination: &verbose,
		},
	}
	app.BashComplete = func(c *cli.Context) {
		for _, command := range c.App.Commands {
			for _, name := range command.Names() {
				if len(name) < 2 {
					continue
				}
				fmt.Fprintln(c.App.Writer, name)
			}
		}
	}
	app.Before = func(context *cli.Context) error {
		prefix = DEFAULT_API_PREFIX
		bucket = DEFAULT_BUCKET

		if concurrency < 1 || concurrency > 16 {
			concurrency = NUM_CPU
		}

		if strings.Count(prefix, "%s") == 1 {
			api = fmt.Sprintf(prefix, bucket)
		} else {
			api = prefix
		}

		return nil
	}
	app.Run(os.Args)
}
