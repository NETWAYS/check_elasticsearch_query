package main

import (
	"github.com/NETWAYS/go-check"
)

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()

	plugin.Name = "check_elasticsearch_query"
	// TODO: Add Readme
	plugin.Readme = ``
	plugin.Version = buildVersion()
	plugin.Timeout = 10


	config := BuildConfigFlags(plugin.FlagSet)
	plugin.ParseArguments()

	err := config.Validate()
	if err != nil {
		check.ExitError(err)
	}

	rc, output, err := config.Run()
	if err != nil {
		check.ExitError(err)
	}

	check.Exit(rc, output)
}
