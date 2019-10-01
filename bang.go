package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

// SliceToMap returns a map of Bangs based on the input slice's Bangs' names.
func SliceToMap(slice []Bang) (bangs map[string]Bang) {
	bangs = map[string]Bang{}
	for _, bang := range slice {
		bangs[bang.Name] = bang
	}
	return
}

// MapToNamedBangs returns a slice of Bangs with their names set according to
// keys from an input map.
func MapToNamedBangs(bangs map[string]Bang) (slice []Bang) {
	slice = []Bang{}
	for key, bang := range bangs {
		bang.Name = key
		slice = append(slice, bang)
	}
	return
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
