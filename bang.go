package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// A Bang is a single registered search shortcut
type Bang struct {
	Description string `json:"description"`
	Format      string `json:"format"`
}

// URL returns the direct search result URL for a given Bang
func (b *Bang) URL(q string) string {
	return fmt.Sprintf(b.Format, q)
}

// Bangs is the registration list of all Bangs. We use a map[string] here for
// faster lookups at runtime than digging through slices.
var Bangs = map[string]Bang{}

func init() {
	data, err := Asset("bangs.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &Bangs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
