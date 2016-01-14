package opts

import (
	"errors"
	"flag"
	"strings"
)

var (
	IsQuiet        bool
	IsVerbose      bool
	PrintName      bool
	PrintNameAndId bool
	ShowOnlyHidden bool
	ShowAll        bool

	InstanceName  string
	InstanceImage string
	InstanceType  string
	InstanceGroup string
	Region        string
	Description   string
)

func GetInstanceId() (string, error) {
	id := flag.Arg(1)
	if id == "" {
		return "", errors.New("Please provide an instance ID.")
	}
	if strings.Contains(id, ":") {
		id = id[strings.Index(id, ":")+1:]
	}
	return id, nil
}
