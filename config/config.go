package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/coreos/pkg/flagutil"
	"github.com/hashicorp/hcl"
)

type Config struct {
	IgnoreText     []string
	IgnoreSource   []string
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	Debug          bool
	MaxTweets      int
}

const maxTweetsDefault = 50

func LoadConfig() Config {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	// default value 0 is important to identify later where the config came from
	maxTweets := flags.Int("max-tweets", 0, "Maximum tweets per feed")
	debug := flags.Bool("debug", false, "Debug")
	configFile := flags.String("config", "", "Configiguration file")

	flagutil.SetFlagsFromEnv(flags, "TWITTER2RSS")
	flags.Parse(os.Args[1:])

	var conf Config
	if len(*configFile) != 0 {
		fmt.Println("Loading configuration file.")
		conf = LoadConfigFile(*configFile)
	}

	if *consumerKey != "" {
		conf.ConsumerKey = *consumerKey
	}
	if *consumerSecret != "" {
		conf.ConsumerSecret = *consumerSecret
	}
	if *accessToken != "" {
		conf.AccessToken = *accessToken
	}
	if *accessSecret != "" {
		conf.AccessSecret = *accessSecret
	}
	conf.Debug = *debug
	// provided via command line
	if *maxTweets != 0 {
		conf.MaxTweets = *maxTweets
	}
	// not provided via config & command line
	if conf.MaxTweets == 0 {
		conf.MaxTweets = maxTweetsDefault
	}

	if conf.ConsumerKey == "" || conf.ConsumerSecret == "" || conf.AccessToken == "" || conf.AccessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	return conf
}

func LoadConfigFile(configFile string) Config {
	var conf Config

	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error: Couldn't read config file %s: %s", configFile, err)
	}

	err = hcl.Decode(&conf, string(configData))
	if err != nil {
		log.Fatalf("Error parsing config file %s: %s", configFile, err)
	}

	return conf
}
