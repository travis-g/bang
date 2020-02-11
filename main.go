package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/kballard/go-shellquote"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

var (
	fs          *flag.FlagSet
	flagURLOnly *bool
)

var helpText = strings.TrimSpace(`
bang - browser launcher, heavily inspired by DuckDuckGo's !bangs

Usage:  bang [options] <bang> <query ...>
        bang [options] <bang> -
        bang [options] -
        bang list
        bang help

Bangs can be configured with any filetype supported by viper
(github.com/spf13/viper) but must be named "bangs", ex. "bangs.json". The CLI
will look for a "bangs" config file in the following directories, in order:

    ~/.config/bang/
    ~/.config/
    ./ ($PWD)

See the source for an example "bangs.json" config file to get started.

Queries can be piped through Stdin if a hyphen is passed as the bang or query
argument. Depending on the location of the hyphen within the args list the full
bang and query can be provided via Stdin.

The system's URL opener will be used by default, but if set, the BROWSER
environment variable will be executed with the chosen bang's URL passed as the
final argument.

Examples:

  Search for cat pictures on Google Images:
  bang gi cat pictures

  Pipe a full query in via Stdin:
  echo "gi cat pictures" | bang -

  Pipe just a query, and output the URL only:
  echo "cat pictures" | bang -url gi -

Options:
`)

func loadConfig() error {
	viper.SetConfigName("bangs")
	viper.AddConfigPath("$HOME/.config/bang/")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	// Unmarshal each key of config file as a proto.Message then stick it back
	// into the global map of named Bangs
	for key, i := range viper.AllSettings() {
		jsonBytes, err := json.Marshal(i)
		if err != nil {
			return fmt.Errorf("fatal error in config file: %s", err)
		}
		buf := bytes.NewReader(jsonBytes)
		var bang Bang
		if jsonpb.Unmarshal(buf, &bang) != nil {
			return fmt.Errorf("fatal error in config file: %s", err)
		}
		bang.Name = key
		Bangs[key] = bang
	}
	return nil
}

func loadFlags(fs *flag.FlagSet) error {
	flagURLOnly = fs.Bool("url", false, "output URL only")
	return fs.Parse(os.Args[1:])
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

func stdinToString() (string, error) {
	in, _ := os.Stdin.Stat()
	interactive := ((in.Mode() & os.ModeCharDevice) != 0)
	if !interactive {
		bytes, err := ioutil.ReadAll(os.Stdin)
		q := strings.TrimSpace(string(bytes))
		return q, err
	}
	return "", fmt.Errorf("%s", "unable to read os.Stdin")
}

func main() {
	fs = flag.NewFlagSet("bang", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", helpText)
		fs.PrintDefaults()
	}

	err := loadConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		default:
			panic(fmt.Errorf("error loading config: %s", err))
		}
	}

	err = loadFlags(fs)
	if err != nil {
		panic(fmt.Errorf("error loading flags: %s", err))
	}

	var args = fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(0)
	}

	// check for special subcommands
	switch args[0] {
	case "-":
		// pull os.Stdin
		q, err := stdinToString()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		args = strings.Split(q, " ")
	case "list":
		fmt.Println(listBangs(Bangs))
		os.Exit(0)
	case "help":
		fs.Usage()
		os.Exit(0)
	}

	// check for bang but not enough args
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "not enough args")
		fs.Usage()
		os.Exit(1)
	}

	// lookup bang
	if bang, ok := Bangs[args[0]]; ok {
		var q string
		if args[1] == "-" {
			q, err = stdinToString()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		} else {
			q = strings.Join(args[1:], " ")
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
		fmt.Fprintln(os.Stderr, args[0], "not a configured bang")
		os.Exit(1)
	}
}
