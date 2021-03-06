package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
)

type ECS struct {
	KEY    string
	SECRET string
}

var ECS_INSTANCE ECS = ECS{KEY: KEY, SECRET: SECRET}

var IsQuiet bool
var IsVerbose bool

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Version = "1.0.0"
	app.Usage = "control Aliyun ECS instances"
	app.HideHelp = true
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		DESCRIBE_INSTANCES,
		DESCRIBE_IMAGES,
		DESCRIBE_REGIONS,
		DESCRIBE_INSTANCE_TYPES,
		DESCRIBE_SECURITY_GROUPS,
		CREATE_INSTANCE,
		ALLOCATE_PUBLIC_IP_ADDRESS,
		START_INSTANCE,
		STOP_INSTANCE,
		RESTART_INSTANCE,
		REMOVE_INSTANCE,
		UPDATE_INSTANCE,
		HIDE_INSTANCE,
		UNHIDE_INSTANCE,
		DESCRIBE_INSTANCE_MONITOR_DATA,
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "quiet, q",
			Usage:       "show only name or ID",
			Destination: &IsQuiet,
		},
		cli.BoolFlag{
			Name:        "verbose, V",
			Usage:       "show more info",
			Destination: &IsVerbose,
		},
	}
	app.BashComplete = func(c *cli.Context) {
		for _, command := range c.App.Commands {
			for _, name := range command.Names() {
				if len(name) <= 2 {
					continue
				}
				fmt.Fprintln(c.App.Writer, name)
			}
		}
	}
	app.Run(os.Args)
}
