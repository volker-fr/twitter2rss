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
	IgnoreText        []string
	IgnoreSource      []string
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessSecret      string
	Debug             bool
	MaxTweets         int
	CombinedFeed      bool
	CombinedFeedHours int
}

const maxTweetsDefault = 50
const combinedFeedHoursDefault = 6

func LoadConfig() Config {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	combinedFeed := flags.Bool("combined-feed", false, "Combine multiple tweets from the same user into a single RSS entry?")
	// default value 0 is important to identify later where the config came from
	maxTweets := flags.Int("max-tweets", 0, "Maximum tweets per feed")
	combinedFeedHours := flags.Int("combined-feed-hours", 0, "if combined-tweet, how many hours should be combined together?")
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
	// command line argument set
	if *debug == true {
		conf.Debug = *debug
	}
	// command line argument set
	if *combinedFeed == true {
		conf.CombinedFeed = *combinedFeed
	}
	// provided via command line
	if *maxTweets != 0 {
		conf.MaxTweets = *maxTweets
	}
	// not provided via config & command line
	if conf.MaxTweets == 0 {
		conf.MaxTweets = maxTweetsDefault
	}
	// provided via command line
	if *combinedFeedHours != 0 {
		conf.CombinedFeedHours = *combinedFeedHours
	}
	// not provided via config & command line
	if conf.CombinedFeedHours == 0 {
		conf.CombinedFeedHours = combinedFeedHoursDefault
	}

	if conf.CombinedFeedHours > 24 {
		fmt.Println("WARNING: please check your configuration. Combining feed by more as 24hours")
		fmt.Println("         doesn't make a difference and the result will be equal to 24")
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
