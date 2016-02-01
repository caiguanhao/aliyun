package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/caiguanhao/aliyun/ecs/opts"
)

var ACTIONS = [][]string{
	{"list-instances", "list", "l", "List all instances, show one if ID is specified"},
	{"list-images", "images", "i", "List all images"},
	{"list-regions", "regions", "n", "List all regions"},
	{"list-instance-types", "types", "t", "List all instance types"},
	{"list-security-groups", "groups", "g", "List all security groups"},
	{"create-instance", "create", "c", "Create an instance"},
	{"allocate-public-ip", "allocate", "a", "Allocate an IP address for an instance"},
	{"start-instance", "start", "s", "Start an instance"},
	{"stop-instance", "stop", "k", "Stop an instance"},
	{"restart-instance", "restart", "r", "Restart an instance"},
	{"remove-instance", "remove", "d", "Remove an instance"},
	{"update-instance", "update", "e", "Update attributes of an instance"},
	{"hide-instance", "hide", "h", "Hide instance from instance list"},
	{"unhide-instance", "unhide", "u", "Un-hide instance from instance list"},
}

func init() {
	flag.BoolVar(&opts.IsQuiet, "q", false, "")
	flag.BoolVar(&opts.IsQuiet, "quiet", false, "")
	flag.BoolVar(&opts.IsVerbose, "v", false, "")
	flag.BoolVar(&opts.IsVerbose, "verbose", false, "")
	flag.BoolVar(&opts.PrintName, "print-name", false, "")
	flag.BoolVar(&opts.PrintNameAndId, "print-name-id", false, "")
	flag.BoolVar(&opts.ShowOnlyHidden, "hidden-only", false, "")
	flag.BoolVar(&opts.ShowAll, "all", false, "")
	flag.StringVar(&opts.InstanceName, "name", "", "")
	flag.StringVar(&opts.InstanceImage, "image", "", "")
	flag.StringVar(&opts.InstanceType, "type", "", "")
	flag.StringVar(&opts.InstanceGroup, "group", "", "")
	flag.StringVar(&opts.Region, "region", "", "")
	flag.StringVar(&opts.Description, "description", "\x00", "")
	flag.Usage = func() {
		if opts.IsQuiet {
			for _, action := range ACTIONS {
				fmt.Println(action[1])
			}
			for _, action := range ACTIONS {
				fmt.Println(action[0])
			}
			if opts.IsVerbose {
				flag.VisitAll(func(flag *flag.Flag) {
					fmt.Print("-")
					if len(flag.Name) > 1 {
						fmt.Print("-")
					}
					fmt.Println(flag.Name)
				})
				fmt.Println("-h")
				fmt.Println("--help")
			}
			return
		}
		fmt.Printf("Usage: %s [OPTION] [ACTION] [TARGET]\n", path.Base(os.Args[0]))
		fmt.Println()
		fmt.Printf("Using Access Key %s\n", KEY)
		fmt.Println()
		fmt.Printf("%-20s  %-11s  %s\n", "Action", "Alias", "Description")
		for _, action := range ACTIONS {
			fmt.Printf("%-20s  %-1s, %-8s  %s\n", action[0], action[2], action[1], action[3])
		}
	}
	flag.Parse()
}

func main() {
	client := ECS{KEY: KEY, SECRET: SECRET}
	var target ECSInterface
	if client.Do(flag.Arg(0), &target) {
		if opts.IsQuiet {
			target.Print()
		} else {
			target.PrintTable()
		}
	} else {
		flag.Usage()
	}
}
