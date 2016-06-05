package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
)

var api string
var bucket string
var prefix string
var accessKey string
var accessSecret string
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
		OSS_LIST,
		OSS_DIFF,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "bucket, b",
			Usage:       "bucket name",
			Value:       DEFAULT_BUCKET,
			Destination: &bucket,
		},
		cli.StringFlag{
			Name:        "prefix, p",
			Usage:       "API prefix",
			Value:       DEFAULT_API_PREFIX,
			Destination: &prefix,
		},
		cli.StringFlag{
			Name:        "key",
			Usage:       "access key",
			Value:       KEY,
			EnvVar:      "ACCESS_KEY",
			Destination: &accessKey,
		},
		cli.StringFlag{
			Name:        "secret",
			Usage:       "access key secret",
			EnvVar:      "ACCESS_SECRET",
			Destination: &accessSecret,
		},
		cli.IntFlag{
			Name:        "concurrency, c",
			Value:       NUM_CPU,
			Usage:       fmt.Sprintf("job concurrency, defaults to number of CPU (%d), max is 16", NUM_CPU),
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
	app.Before = func(c *cli.Context) error {
		if !c.GlobalIsSet("secret") {
			accessSecret = SECRET
		}

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
