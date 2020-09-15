check_elasticsearch_query
=========================

Check the total hits/results of an elasticsearch query over the API of elasticsearch.

The plugin is currently capable to return the total hits of documents based on a provided query string.
For more information to the syntax, please visit: [Elasticserach documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html)

## Usage

```
Arguments:
  -H, --host string       Host name, IP Address of the elasticsearch host (default "127.0.0.1")
  -p, --port int          Port number of the elasticsearch host (default 9200)
  -U, --user string       Username of the elasticsearch host
  -P, --password string   Password of the user
  -q, --query string      Elasticsearch query, e.g. 'event.dataset=sample_web_logs and @timestamp: [ 2020-09-05T20:44:46.291Z ]'
  -i, --index string      The index which will be used (default "_all")
      --msgchars int      Number of characters to display in latest message as integer. To disable set value to 0 (default 255)
      --msgkey string     For query searches only. Index of message to display. eg. message (default "message")
  -c, --critical int      Critical threshold for total hits (default: 10) (default 10)
  -w, --warning int       Warning threshold for total hits (default: 5) (default 5)
  -t, --timeout int       Abort the check after n seconds (default 10)
  -d, --debug             Enable debug mode
  -v, --verbose           Enable verbose mode
  -V, --version           Print version and exit
```

## API Documentation

Full API documentation is available at [https://www.elastic.co/](https://www.elastic.co/guide/en/elasticsearch/reference/7.x/rest-apis.html).

## License

Copyright (c) 2020 [NETWAYS GmbH](mailto:info@netways.de) \
Copyright (c) 2020 [Philipp Dorschner](mailto:philipp.dorschner@netways.de)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see [gnu.org/licenses](https://www.gnu.org/licenses/).
