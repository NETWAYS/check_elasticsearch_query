package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/pflag"
	"strconv"
	"strings"
)

var (
	response             map[string]interface{}
	totalHits            int
	latestMessage        string
	availableMessageKeys string
)

type Config struct {
	host                 string
	port                 int
	user                 string
	password             string
	index                string
	query                string
	// filter               string
	messageCharacters    int
	messageKey           string
	exclude              bool
	excludeKey           string
	excludeValues        []string
	warning              int
	critical             int
	paginateSearchResult int
}

func BuildConfigFlags(fs *pflag.FlagSet) (config *Config) {
	config = &Config{}

	fs.StringVarP(&config.host, "host", "H", "127.0.0.1",
		"host name, IP Address of the elasticsearch host")

	fs.IntVarP(&config.port, "port", "p", 9200,
		"port number of the elasticsearch host")

	fs.StringVarP(&config.user, "user", "U", "",
		"Username of the elasticsearch host")

	fs.StringVarP(&config.password, "password", "P", "",
		"password of the user")

	fs.StringVarP(&config.query, "query", "q", "",
		"Elasticsearch query, e.g.:" +
		      "\"event.dataset:sample_web_logs AND @timestamp:[now-5m TO now] AND message:Parallels\"")

	fs.StringVarP(&config.index, "index", "i", "_all",
		"The index which will be used")

	/*
	TODO: filter objects should be implemented WIP
	fs.StringVarP(&config.filter, "filter", "f", "",
		"Name of saved filter in Kibana")
	 */

	fs.IntVar(&config.messageCharacters, "msgchars", 255,
		"Number of characters to display in latest message as integer. To disable set value to 0")

	fs.StringVar(&config.messageKey, "msgkey", "message",
		"For query searches only. index of message to display. eg. message")

	fs.BoolVar(&config.exclude, "exclude", false,
		"Exludes specified values of a key")

	fs.StringVar(&config.excludeKey, "excludekey", "",
		"Searches for the specified key to be excluded")

	fs.StringSliceVar(&config.excludeValues, "excludevalues", config.excludeValues,
		"The values (comma seperated) to be excluded of a given exclude key")

	fs.IntVar(&config.paginateSearchResult, "paginateSearchResult",1,
		"Returns the top x matching documents")

	fs.IntVarP(&config.critical, "critical", "c", 10,
		"critical threshold for total hits")

	fs.IntVarP(&config.warning, "warning", "w", 5,
		"warning threshold for total hits")

	_ = fs.MarkHidden("paginateSearchResult")
	_ = fs.MarkHidden("exclude")

	return
}

func (c *Config) Validate() (err error) {

	if c.user != "" && c.password == "" {
		err = fmt.Errorf("password must be configured")
		return
	} else if c.password != "" && c.user == "" {
		err = fmt.Errorf("user must be configured")
		return
	}

	if c.query == "" {
		err = fmt.Errorf("query has to be configured")
		return
	}

	if strings.Contains(c.query, "="){
		err = fmt.Errorf("the character \"=\" in the query is not allowed")
		return
	}

	if c.messageKey != ""  && c.query == "" {
		err = fmt.Errorf("query has to be configured to use --msgkey")
		return
	}

	if c.excludeKey != "" && len(c.excludeValues) == 0 {
		err = fmt.Errorf("exlude option needs an exclude key and exclue values")
		return
	} else if c.excludeKey == "" && len(c.excludeValues) > 0 {
		err = fmt.Errorf("exlude option needs an exclude key and exclue values")
		return
	} else if c.excludeKey != "" && len(c.excludeValues) > 0 {
		c.exclude = true
	}

	// Validation complete

	return nil
}

func (c *Config) Run() (returnCode int, output string, err error) {

	cfg := elasticsearch.Config{
		Addresses: []string{"http://" + c.host + ":" + strconv.Itoa(c.port)},
		Username:  c.user,
		Password:  c.password,
	}

	elasticClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		err = fmt.Errorf("error creating the client: %s", err)
		return
	}

	elasticInfo, err := elasticClient.Info()
	if err != nil {
		err = fmt.Errorf("error getting response: %s", err)
		return
	}
	defer elasticInfo.Body.Close()

	if elasticInfo.IsError() {
		err = fmt.Errorf("error: %s", elasticInfo.String())
		return
	}

	err = json.NewDecoder(elasticInfo.Body).Decode(&response)
	if err != nil {
		err = fmt.Errorf("error parsing the response body: %s", err)
		return
	}

	serverVersion := strings.Split(
		response["version"].(
			map[string]interface{})["number"].(string), ".")

	clientVersion := strings.Split(elasticsearch.Version, ".")

	if serverVersion[0] != clientVersion[0] {
		err = fmt.Errorf("major version of client and server does not match")
		return
	}

	/*
	TODO: Add more request options, like 'match'.
	Request body:
	{
	 "query": {
	   "query_string": {
	     "query": "example"
	   }
	 }
	}
	 */
	if c.exclude {
		c.query += " NOT ("
		for _, excludeValue := range c.excludeValues {
			c.query += c.excludeKey + ":" + excludeValue + " "
		}

		c.query += ")"
	}

	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": c.query,
			},
		},
	}

	err = json.NewEncoder(&buf).Encode(query)
	if err != nil {
		err = fmt.Errorf("error encoding query: %s", err)
		return
	}

	elasticInfo, err = elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(c.index),
		elasticClient.Search.WithBody(&buf),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
		elasticClient.Search.WithSize(c.paginateSearchResult),
		elasticClient.Search.WithFilterPath(),
	)
	if err != nil {
		err = fmt.Errorf("error getting response: %s", err)
		return
	}
	defer elasticInfo.Body.Close()

	if elasticInfo.IsError() {
		var e map[string]interface{}
		err = json.NewDecoder(elasticInfo.Body).Decode(&e)

		if err != nil {
			err = fmt.Errorf("error parsing the response body: %s", err)
			return
		} else {
			err = fmt.Errorf("[%s] %s: %s",
				elasticInfo.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			return
		}
	}

	if err = json.NewDecoder(elasticInfo.Body).Decode(&response); err != nil {
		err = fmt.Errorf("error parsing the response body: %s", err)
		return
	}

	totalHits = int(
		response["hits"].(
			map[string]interface{})["total"].(
				map[string]interface{})["value"].(float64))

	// TODO: Refactor
	hitObjects := response["hits"].(
		map[string]interface{})["hits"].(
			[]interface{})

	for count, hit := range hitObjects {
		if count == (c.paginateSearchResult - 1) {
			if hit.(map[string]interface{})["_source"].(map[string]interface{})[c.messageKey] == nil {
				for key, _ := range hit.(map[string]interface{})["_source"].(map[string]interface{}){
					availableMessageKeys += key + "\n"
				}

				err = fmt.Errorf("no message key for \"" +
					c.messageKey + "\" was found. Available keys:\n%v", availableMessageKeys)

				return
			} else {
				if c.exclude {
					if hit.(map[string]interface{})["_source"].(map[string]interface{})[c.excludeKey] == nil{
						for key, _ := range hit.(map[string]interface{})["_source"].(map[string]interface{}){
							availableMessageKeys += key + "\n"
						}
						err = fmt.Errorf("exclude key: \"" +
							c.excludeKey + "\" was not found. Available keys:\n%v", availableMessageKeys)

						return
					}
				}
				latestMessage = hit.(
					map[string]interface{})["_source"].(
						map[string]interface{})[c.messageKey].(string)
			}
		}
	}

	if totalHits >= c.critical {
		returnCode = check.Critical
	} else if totalHits >= c.warning {
		returnCode = check.Warning
	} else if totalHits == 0 {
		returnCode = check.OK
		output = "The query return 0 hits"
		return
	} else {
		returnCode = check.OK
	}

	output = "Total hits: " + strconv.Itoa(totalHits)

	if c.messageCharacters != 0 {
		if len(latestMessage) <= c.messageCharacters {
			output += "\n" + c.messageKey + ": "  + latestMessage
		} else {
			output += "\n" + c.messageKey + ": "  + latestMessage[0:c.messageCharacters]
		}
	}

	return
}
