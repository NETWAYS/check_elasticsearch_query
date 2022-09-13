package main

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	c := &Config{}
	assert.Error(t, c.Validate())

	c.host = "127.0.0.1"
	c.query = "machine.ram: 16106127360"
	c.index = "_all"
	c.messageCharacters = 10
	c.messageKey = "message"

	assert.NoError(t, c.Validate())
}

func TestBuildConfigFlags(t *testing.T) {
	fs := &pflag.FlagSet{}
	BuildConfigFlags(fs)

	assert.True(t, fs.HasFlags())
}

func TestConfig_Run(t *testing.T) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
		Username: "User",
		Password: "Password",
	}

	elasticClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		t.Fatalf("error client: %s", err)
	}

	response, err := elasticClient.Info()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response == nil {
		t.Fatalf("Unexpected response: %v", response)
	}
}
