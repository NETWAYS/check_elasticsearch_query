check_elasticsearch_query
=========================

Check the total hits/results of an elasticsearch query over the API of elasticsearch.

The plugin is currently capable to return the total hits of documents based on a provided query string.
For more information to the syntax, please visit: [Elasticserach documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html)

## Usage

```
Arguments:
  -H, --host string             host name, IP Address of the elasticsearch host (default "127.0.0.1")
  -p, --port int                port number of the elasticsearch host (default 9200)
  -U, --user string             Username of the elasticsearch host
  -P, --password string         password of the user
  -q, --query string            Elasticsearch query, e.g.:"event.dataset:sample_web_logs and @timestamp:[now-5m TO now] and message:Parallels"
  -i, --index string            The index which will be used (default "_all")
      --msgchars int            Number of characters to display in latest message as integer. To disable set value to 0 (default 255)
      --msgkey string           For query searches only. index of message to display. eg. message (default "message")
      --excludekey string       Searches for the specified key to be excluded
      --excludevalues strings   The values (comma seperated) to be excluded of a given exclude key
  -c, --critical int            critical threshold for total hits (default 10)
  -w, --warning int             warning threshold for total hits (default 5)
  -t, --timeout int             Abort the check after n seconds (default 10)
  -d, --debug                   Enable debug mode
  -v, --verbose                 Enable verbose mode
  -V, --version                 Print version and exit
```

## Example

```
$ ./check_elasticsearch --host=123.123.123.123 --query=message:207.73.150.150 --index="example-*" --msgkey=message --excludekey=response --excludevalues=200,301
OK - Total hits: 1
message: 207.73.150.150 - - [2018-09-12T13:58:13.033Z] "GET /security-analytics HTTP/1.1" 404 748 "-" "Mozilla/5.0 (X11; Linux x86_64; rv:6.0a1) Gecko/20110421 Firefox/6.0a1"

$ ./check_elasticsearch --host=123.123.123.123 --query=message:"@timestamp:[now-2d TO now] AND message:207.73.150.150" --index="example-*" --msgkey=message
  CRITICAL - Total hits: 16
  message: 207.73.150.150 - - [2018-08-19T14:35:13.301Z] "GET /security-analytics HTTP/1.1" 404 748 "-" "Mozilla/5.0 (X11; Linux x86_64; rv:6.0a1) Gecko/20110421 Firefox/6.0a1"
```

## API Documentation

Full API documentation is available at [https://www.elastic.co/](https://www.elastic.co/guide/en/elasticsearch/reference/7.x/rest-apis.html).

## License

Copyright (c) 2020 [NETWAYS GmbH](mailto:info@netways.de)

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
