package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

// TODO: make exported
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
	// fall back to system URL opener
	return browser.OpenURL(url)
}

func main() {
	fs = flag.NewFlagSet("bang", flag.ExitOnError)
	loadConfig()
	loadFlags(fs)

	// check for special subcommands
	switch fs.Arg(0) {
	case "list":
		fmt.Println(listBangs(Bangs))
		os.Exit(0)
	case "help":
		fs.Usage()
		os.Exit(0)
	}

	// check for bang but not enough args
	if len(fs.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "not enough args")
		fs.Usage()
		os.Exit(1)
	}

	// lookup bang
	if bang, ok := Bangs[fs.Arg(0)]; ok {
		var q string
		if fs.Arg(1) == "-" {
			// TODO: read from os.Stdin
			in, _ := os.Stdin.Stat()
			interactive := ((in.Mode() & os.ModeCharDevice) != 0)
			if !interactive {
				bytes, _ := ioutil.ReadAll(os.Stdin)
				q = strings.TrimSpace(string(bytes))
			}
		} else {
			q = strings.Join(fs.Args()[1:], " ")
		}
		url := bang.URL(q)

		if *flagURLOnly {
			fmt.Fprintln(os.Stdout, url)
			os.Exit(0)
		}

		err := launch(url)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, fs.Arg(0), "not a configured bang")
		os.Exit(1)
	}
}
