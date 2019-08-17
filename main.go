package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func loadConfig() {
	viper.SetConfigName("bangs")
	viper.AddConfigPath("$HOME/.config/bangs/")
	viper.AddConfigPath("$HOME/.bangs")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&Bangs)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "not enough args")
		os.Exit(1)
	}

	loadConfig()

	if bang, ok := Bangs[os.Args[1]]; ok {
		terms := os.Args[2:]
		q := bang.URL(url.QueryEscape(strings.Join(terms[:], " ")))
		fmt.Fprintln(os.Stdout, q)
		// err := browser.OpenURL(q)
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, err)
		// 	os.Exit(1)
		// }
	} else {
		fmt.Fprintln(os.Stderr, os.Args[1], "not found")
		os.Exit(1)
	}
}
