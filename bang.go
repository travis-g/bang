package main

import "fmt"

// A Bang is a single registered search shortcut
type Bang struct {
	Description string
	Format      string
}

// URL returns the direct search result URL for a given Bang
func (b *Bang) URL(q string) string {
	return fmt.Sprintf(b.Format, q)
}

// Bangs is the registration list of all Bangs. We use a map[string] here for
// faster lookups at runtime than digging through slices.
var Bangs = map[string]Bang{
	"w": Bang{
		"WikiPedia",
		"https://wikipedia.org/wiki/Special:Search/%s",
	},
	"bi": Bang{
		"Bing Images",
		"https://www.bing.com/images/search?q=%s",
	},
	"gh": Bang{
		"GitHub",
		"https://github.com/search?q=%s",
	},
	"gi": Bang{
		"Google Images",
		"https://www.google.com/search?tbm=isch&q=%s",
	},
	"r": Bang{
		"Reddit",
		"https://www.reddit.com/search?q=%s",
	},
	"reddits": Bang{
		"Reddit subreddit",
		"https://www.reddit.com/r/%s",
	},
}
