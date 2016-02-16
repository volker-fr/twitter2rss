# Filter based on the source of a tweet. Most helpful for services such as
# Foursquare
ignoreSource = [
    "untappd",
    "foursquare"
]

# Filter based on the text of a tweet.
ignoreText = [
    "I just backed.*kickstarter.com",
    "I just backed.*kck.st"
]

# Authentication data from twitter
#consumerKey = "myConsumerKey"
#consumerSecret = "myConsumerSecret"
#accessToken = "myAccessToken"
#accessSecret = "myAccessSecret"

# Should the program output debug information
#   default: false
#debug = true

# Maxium tweets to get from the home timeline
#   default: 50
#maxTweets = 100

# Combine feed entries of the same author into one rss entry?
#   default: false
combinedFeed = true

# If tweets are combined, define the amount of hours you want to use
# E.g. 6 means every 6 hours there is a new feed entry
# This only works in combination with combedFeed. Else it is ignored
#   default: 6
#combinedFeedHours = 12
