package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dogtools/dog/types"
)

type userArgs struct {
	help     bool
	version  bool
	info     bool
	taskName string
	taskArgs []string
}

var knownFlags = [...]string{
	"-i", "--info",
	"-h", "--help",
	"-v", "--version",
}

func printVersion() {
	fmt.Println("Dog version: " + version)
}

func printHelp() {
	fmt.Println(`Usage: dog
       dog [OPTIONS] TASK [ARGS]
       dog [--help] [--version]

Dog is a command line application that executes tasks.

Options:

  -i, --info     Print execution info (duration, statuscode) after task execution
  -h, --help     Print usage information and help
  -v, --version  Print version information`)
}

func printNoValidDogfile() {
	fmt.Println(`Error: No valid Dogfile in current directory
Need help? --> dog --help
More info  --> https://github.com/dogtools/dog`)
}

func printTasks(tm types.TaskMap) {

	maxCharSize := 0
	for taskName, _ := range tm {
		if len(taskName) > maxCharSize {
			maxCharSize = len(taskName)
		}
	}

	var tasks []string
	for k := range tm {
		tasks = append(tasks, k)
	}
	sort.Strings(tasks)

	for _, taskName := range tasks {
		spaces := strings.Repeat(" ", maxCharSize-len(taskName)+2)
		fmt.Printf("%s%s%s\n", taskName, spaces, tm[taskName].Description)
	}
}

func parseArgs(args []string) (a userArgs, err error) {

	// default values
	a = userArgs{
		help:     false,
		version:  false,
		info:     false,
		taskName: "",
		taskArgs: []string{},
	}

	// iterate over all provided arguments
	for i, arg := range args {

		if arg == "--help" || arg == "-h" {
			if i == 0 && len(args) == 1 && a.taskName == "" {
				a.help = true
				return a, nil
			} else {
				return a, fmt.Errorf("Error: %s doesn't accept additional parameters", arg)
			}
		}

		if arg == "--version" || arg == "-v" {
			if i == 0 && len(args) == 1 && a.taskName == "" {
				a.version = true
				return a, nil
			} else {
				return a, fmt.Errorf("Error: %s doesn't accept additional parameters", arg)
			}
		}

		if arg == "--info" || arg == "-i" {
			if a.taskName == "" {
				a.info = true
			} else {
				return a, fmt.Errorf("Error: %s is not a valid task argument", arg)
			}
		}

		if string(arg[0]) != "-" {
			if a.taskName == "" {
				a.taskName = arg
			} else {
				return a, fmt.Errorf("Error: only one task can be executed at a time")
			}
		} else {
			validArg := false
			for _, f := range knownFlags {
				if arg == f {
					validArg = true
				}
			}
			if !validArg {
				return a, fmt.Errorf("Error: %s is not a valid argument", arg)
			}
		}
	}

	return a, nil
}
