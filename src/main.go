package main

import (
	"flag"
	"fmt"
	"os"
)

var SITES = []string{
	"offspring",
	"topps",
	"funko",
	"end",
}

func main() {
	var (
		tasks    string
		proxies  string
		monitor  bool
		newTasks bool
		site     string
	)

	flag.StringVar(&tasks, "t", "", "tasks")
	flag.StringVar(&proxies, "p", "", "proxies")
	flag.BoolVar(&monitor, "m", false, "monitor")
	flag.BoolVar(&newTasks, "n", false, "create new task sheet")
	flag.StringVar(&site, "s", "", "site")

	flag.Parse()

	if err := run(tasks, proxies, monitor, site, newTasks); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(tasks string, proxies string, monitor bool, site string, newTasks bool) error {
	if err := deviceCheck(); err != nil {
		return err
	}
	err := initialCheck()
	if err != nil {
		return err
	}

	if newTasks {
		if site == "" {
			return ErrSiteRequired
		}
		if err = createTaskTemplate(site, newTasks); err != nil {
			return err
		}
	}

	return nil
}
