package main

import (
	"fmt"
	"strings"
)

// A Bang is a single registered search shortcut
type Bang struct {
	Description string `json:"description" mapstructure:"description"`
	Format      string `json:"format" mapstructure:"format"`
}

const symbol = "{{{s}}}"

// URL returns the direct search result URL for a given Bang
func (b *Bang) URL(q string) string {
	return fmt.Sprint(strings.Replace(b.Format, symbol, q, 1))
}

// Bangs is the registration list of all Bangs. We use a map[string] here for
// faster lookups at runtime than digging through slices.
var Bangs = map[string]Bang{}

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
