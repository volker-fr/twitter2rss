package config

import (
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	conf := LoadConfigFile("config_test.hcl")

	if conf.IgnoreText[0] != "I just backed.*kickstarter.com" {
		t.Error("Expected IgnoreText[0] to be 'I just backed.*kickstarter.com', got:", conf.IgnoreText[0])
	}
	if conf.IgnoreSource[0] != "untappd" {
		t.Error("Expected IgnoreSource[0] to be 'untappd', got:", conf.IgnoreSource[0])
	}
	if conf.ConsumerKey != "myConsumerKey" {
		t.Error("Expected ConsumerKey myConsumerKey, got:", conf.ConsumerKey)
		}
	if conf.ConsumerSecret != "myConsumerSecret" {
		t.Error("Expected ConsumerSecret myConsumerSecret, got:", conf.ConsumerSecret)
	}
	if conf.AccessToken != "myAccessToken" {
		t.Error("Expected AccessToken myAccessToken, got:", conf.AccessToken)
	}
	if conf.AccessSecret != "myAccessSecret" {
		t.Error("Expected AccessSecret myAccessSecret, got:", conf.AccessSecret)
	}
	if conf.Debug != true {
		t.Error("Expected Debug true, got:", conf.Debug)
	}
}
