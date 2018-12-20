package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/browser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "not enough args")
		os.Exit(1)
	}
	if bang, ok := Bangs[os.Args[1]]; ok {
		terms := os.Args[2:]
		q := url.QueryEscape(strings.Join(terms[:], " "))
		browser.OpenURL(bang.URL(q))
	} else {
		fmt.Fprintln(os.Stderr, os.Args[1], "not found")
		os.Exit(1)
	}
}
