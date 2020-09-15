package main

import (
	"github.com/NETWAYS/go-check"
)

const readme = `Check the total hits/results of an elasticsearch query over the API of elasticsearch.

The plugin is currently capable to return the total hits of documents based on a provided query string.
For more information to the syntax, please visit:
https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html

https://github.com/NETWAYS/check_elasticsearch_query

Copyright (c) 2020 NETWAYS GmbH <info@netways.de>
Copyright (c) 2020 Philipp Dorschner <philipp.dorschner@netways.de

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see https://www.gnu.org/licenses/`

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()

	plugin.Name = "check_elasticsearch_query"
	plugin.Readme = readme
	plugin.Version = buildVersion()
	plugin.Timeout = 10

	config := BuildConfigFlags(plugin.FlagSet)
	plugin.ParseArguments()

	err := config.Validate()
	if err != nil {
		check.ExitError(err)
	}

	returnCode, output, err := config.Run()
	if err != nil {
		check.ExitError(err)
	}

	check.Exit(returnCode, output)
}
