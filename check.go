package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/pflag"
	"log"
	"strconv"
)

var (
	// For deserialization of the response into a map
	Response      map[string]interface{}
	TotalHits     int
	LatestMessage string
	ReturnCode    int
)

type Config struct {
	Host                 string
	Port                 int
	User                 string
	Password             string
	Index                string
	Query                string
	Filter               string
	MessageCharacters    int
	MessageKey           string
	SourceKey            string
	Warning              int
	Critical             int
	Validated         	 bool
	PaginateSearchResult int
}

func BuildConfigFlags(fs *pflag.FlagSet) (config *Config) {
	config = &Config{}

	fs.StringVarP(&config.Host, "host", "H", "127.0.0.1",
		"Host name, IP Address of the elasticsearch host")

	fs.IntVarP(&config.Port, "port", "p", 9200,
		"Port number of the elasticsearch host")

	fs.StringVarP(&config.User, "user", "U", "",
		"Username of the elasticsearch host")

	fs.StringVarP(&config.Password, "password", "P", "",
		"Password of the user")

	fs.StringVarP(&config.Query, "query", "q", "",
		"Elasticsearch query, e.g. 'event.dataset=sample_web_logs and @timestamp: [ 2020-09-05T20:44:46.291Z ]'")

	fs.StringVarP(&config.Index, "index", "i", "_all",
		"The index which will be used")

	// TODO: WIP
	//fs.StringVarP(&config.Filter, "filter", "f", "",
	//	"Name of saved filter in Kibana")

	fs.IntVar(&config.MessageCharacters, "msgchars", 255,
		"Number of characters to display in latest message as integer. To disable set value to 0")

	fs.StringVar(&config.MessageKey, "msgkey", "message",
		"For query searches only. Index of message to display. eg. message")

	fs.IntVar(&config.PaginateSearchResult, "PaginateSearchResult",1,
		"Returns the top x matching documents")

	fs.IntVarP(&config.Critical, "critical", "c", 10,
		"Critical threshold for total hits (default: 10)")

	fs.IntVarP(&config.Warning, "warning", "w", 5,
		"Warning threshold for total hits (default: 5)")

	_ = fs.MarkHidden("PaginateSearchResult")
	
	return
}

func (c *Config) Validate() (err error) {
	c.Validated = false

	if c.User != "" && c.Password == "" {
		err = fmt.Errorf("password must be configured")
		return
	} else if c.Password != "" && c.User == "" {
		err = fmt.Errorf("user must be configured")
		return
	}

	if c.Query == "" {
		err = fmt.Errorf("query has to be configured")
		return
	}

	if c.MessageKey != ""  && c.Query == "" {
			err = fmt.Errorf("query has to be configured to use --msgkey")
			return
	}

	// Validation complete
	c.Validated = true

	return nil
}

func (c *Config) Run() (ReturnCode int, output string, err error) {
	if !c.Validated {
		panic("you need to call Validate() before Run()")
	}

	cfg := elasticsearch.Config{
		Addresses: []string{"http://" + c.Host + ":" + strconv.Itoa(c.Port)},
		Username:  c.User,
		Password:  c.Password,
	}

	ElasticClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		err = fmt.Errorf("error creating the client: %s", err)
	}

	ElasticInfo, err := ElasticClient.Info()
	if err != nil {
		err = fmt.Errorf("error getting response: %s", err)
		return
	}
	defer ElasticInfo.Body.Close()

	if ElasticInfo.IsError() {
		err = fmt.Errorf("error: %s", ElasticInfo.String())
	}

	if err := json.NewDecoder(ElasticInfo.Body).Decode(&Response); err != nil {
		err = fmt.Errorf("error parsing the response body: %s", err)
	}
	if elasticsearch.Version != Response["version"].(map[string]interface{})["number"] {
		err = fmt.Errorf("version of client and server are not equal")
	}

	// TODO: Add more request options, like 'match'.
	// Request body:
	// {
	//  "query": {
	//    "query_string": {
	//      "query": "example"
	//    }
	//  }
	//}
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": c.Query,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		err = fmt.Errorf("error encoding query: %s", err)
	}

	ElasticInfo, err = ElasticClient.Search(
		ElasticClient.Search.WithContext(context.Background()),
		ElasticClient.Search.WithIndex(c.Index),
		ElasticClient.Search.WithBody(&buf),
		ElasticClient.Search.WithTrackTotalHits(true),
		ElasticClient.Search.WithPretty(),
		ElasticClient.Search.WithSize(c.PaginateSearchResult),
		ElasticClient.Search.WithFilterPath(),
	)

	if err != nil {
		err = fmt.Errorf("error getting response: %s", err)
		return
	}
	defer ElasticInfo.Body.Close()

	if ElasticInfo.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(ElasticInfo.Body).Decode(&e); err != nil {
			err = fmt.Errorf("error parsing the response body: %s", err)
		} else {
			log.Fatalf("[%s] %s: %s",
				ElasticInfo.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(ElasticInfo.Body).Decode(&Response); err != nil {
		err = fmt.Errorf("error parsing the response body: %s", err)
	}

	TotalHits = int(Response["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for count, hit := range Response["hits"].(map[string]interface{})["hits"].([]interface{}) {
		if count == (c.PaginateSearchResult - 1) {
			//log.Printf(" * ID=%s", hit.(map[string]interface{})["_id"])
			LatestMessage = hit.(map[string]interface{})["_source"].(map[string]interface{})[c.MessageKey].(string)
		}
	}

	if TotalHits >= c.Critical {
		ReturnCode = check.Critical
	} else if TotalHits >= c.Warning {
		ReturnCode = check.Warning
	} else {
		ReturnCode = check.OK
	}

	output = "Total hits: " + strconv.Itoa(TotalHits)

	if c.MessageCharacters != 0 {
		if len(LatestMessage) <= c.MessageCharacters {
			output += "\n" + c.MessageKey + ": "  + LatestMessage
		} else {
			output += "\n" + c.MessageKey + ": "  + LatestMessage[0:c.MessageCharacters]
		}
	}

	return
}
