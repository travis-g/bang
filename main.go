package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

var (
	fs          *flag.FlagSet
	flagURLOnly *bool
)

func loadConfig() {
	viper.SetConfigName("bangs")
	viper.AddConfigPath("$HOME/.config/bangs/")
	viper.AddConfigPath("$HOME/.bangs")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	err = viper.Unmarshal(&Bangs)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func loadFlags(fs *flag.FlagSet) {
	flagURLOnly = fs.Bool("url", false, "output URL only")
	fs.Parse(os.Args[1:])
}

func launch(url string) error {
	if browser, ok := os.LookupEnv("BROWSER"); ok {
		args, err := shellquote.Split(browser)
		if err != nil {
			return err
		}
		args = append(args, url)
		var launcher *exec.Cmd
		if len(args) > 1 {
			launcher = exec.Command(args[0], args[1:]...)
		} else {
			launcher = exec.Command(args[0])
		}
		return launcher.Start()
	}
	// fall back to system opener
	return browser.OpenURL(url)
}

func main() {
	fs = flag.NewFlagSet("default", flag.ExitOnError)
	loadConfig()
	loadFlags(fs)

	if len(fs.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "not enough args")
		os.Exit(1)
	}

	if bang, ok := Bangs[fs.Arg(0)]; ok {
		terms := fs.Args()
		q := bang.URL(url.QueryEscape(strings.Join(terms[1:], " ")))

		if *flagURLOnly {
			fmt.Fprintln(os.Stdout, q)
			os.Exit(0)
		}

		err := launch(q)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, fs.Arg(0), "not a configured bang")
		os.Exit(1)
	}
}
