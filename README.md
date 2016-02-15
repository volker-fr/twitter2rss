# Twitter2RSS
The application will create a RSS feed from a Twitter home timeline.

It was created out of a test project to play around with go.

## Build
```
% make build
```

If you build on another platform as Linux (e.g. MacOS) and want to copy it
to your Linux or FreeBSD server:
```
% make linux
% make freebsd
```

## Usage
```
% ./twitter2rss --help
Usage of user-auth:
  -access-secret string
        Twitter Access Secret
  -access-token string
        Twitter Access Token
  -config string
        Configiguration file
  -consumer-key string
        Twitter Consumer Key
  -consumer-secret string
        Twitter Consumer Secret
  -debug
        Debug
  -max-tweets int
        Maximum tweets per feed
```

## Configuration
You can provide all required configuration from the command line (see Usage).

The required API keys can be received by creating your own Twitter app on https://apps.twitter.com/


## Environment variables
```
% export TWITTER2RSS_CONSUMER_KEY="YOUR_TWITTER_CONSUMER_KEY"
% export TWITTER2RSS_CONSUMER_SECRET="YOUR_TWITTER_CONSUMER_SECRET"
% export TWITTER2RSS_ACCESS_TOKEN="YOUR_TWITTER_ACCESS_TOKEN"
% export TWITTER2RSS_ACCESS_SECRET="YOUR_TWITTER_ACCESS_SECRET"
```

### Configuration File

See the configuration file twitter.hcl in this repository.

## Filters
You can ignore tweets by configuring regex filters.

### Source filter
The source filter ignores all tweets based on the client used for a tweet. That is helpful to filter all automatically created tweets from various services such as Foursquare out.

Currently I do not know how to find out the Source string of a tweet without querying the API.

### Text filter
The text filter ignores all tweets that match a specific substring in the tweet text.

Note: The tweet text used is not the raw tweet text provided by the Twitter API but the final text after building the entire tweet object together. If there is an URL in a tweet, this filter can also match the summary of the URL.

## Thanks
Thanks to Jon Bodner for helping with go related questions.
