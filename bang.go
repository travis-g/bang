package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

// A Bang is a single registered search shortcut.
type Bang struct {
	ID          string `json:"id" mapstructure:"id"`
	Description string `json:"description" mapstructure:"description"`
	Format      string `json:"format" mapstructure:"format"`

	// If true, queries will be passed without URL encoding/query escaping
	Unescaped bool `json:"unescaped,omitempty" mapstructure:"unescaped"`
}

const symbol = "{{{s}}}"

// URL returns the direct query URL for a Bang.
func (b *Bang) URL(q string) string {
	if !b.Unescaped {
		q = url.QueryEscape(q)
	}
	return fmt.Sprint(strings.Replace(b.Format, symbol, q, 1))
}

// Bangs is the registration list of all Bangs. We use a map[string] here for
// faster lookups at runtime than digging through slices.
var Bangs = map[string]Bang{}

func listBangs(bangs map[string]Bang) string {
	var buf bytes.Buffer
	for key, bang := range bangs {
		fmt.Fprintf(&buf, "%s - %s\n", key, bang.Description)
	}
	return strings.TrimSpace(buf.String())
}

// Filter filters a map of Bangs based on a predicate.
//
//    fmt.Println(Filter(Bangs, func(v string) bool {
//      return strings.Contains(strings.ToLower(v), strings.ToLower("Bing"))
//    }))
func Filter(vs map[string]Bang, f func(string) bool) map[string]Bang {
	vsf := make(map[string]Bang, 0)
	for i, v := range vs {
		if f(i) || f(v.Description) {
			vsf[i] = v
		}
	}
	return vsf
}
