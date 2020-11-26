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
	lib "github.com/travis-g/bang"
)

var (
	fs            *flag.FlagSet
	flagURLOnly   bool
	flagConfigDir string
)

var helpText = strings.TrimSpace(`
bang - browser launcher, heavily inspired by DuckDuckGo's !bangs

Usage:  bang [OPTIONS] [BANG] [QUERY]...
        bang [OPTIONS] [BANG] -
        bang [OPTIONS] -
        bang list
        bang help

Bangs can be configured in YAML, JSON, TOML, HCL, and several other types, but
must be named "bangs", ex. "bangs.json". The CLI will look for a "bangs" config
file in the following directories, in order:

    ~/.config/bang/
    ~/.config/
    ./ ($PWD)

See the source repository for an example "bangs.json" config file.

Queries can be piped through stdin if a hyphen is passed as the bang or query
argument. Depending on the location of the hyphen within the args list both
the bang and query can be provided via stdin.

The system's URL opener will be used by default, but if set, the BROWSER
environment variable will be executed with the chosen bang's URL passed as the
final argument.

Examples:

    Search for cat pictures on Google Images:
    bang gi cat pictures

    Pipe a full query in via stdin:
    echo "gi cat pictures" | bang -

    Pipe just a query, and output the URL only:
    echo "cat pictures" | bang -url gi -

Options:
`)

var (
	// ConfigPaths is the list of paths where the config file can be
	// stored.
	ConfigPaths = []string{
		"$HOME/.config/bang/",
		"$HOME/.config/",
		".",
	}
)

func loadConfig() error {
	viper.SetConfigName("bangs")
	for _, path := range ConfigPaths {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
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
		var bang lib.Bang
		if jsonpb.Unmarshal(buf, &bang) != nil {
			return fmt.Errorf("fatal error in config file: %s", err)
		}
		bang.Name = key
		lib.Bangs[key] = bang
	}
	return nil
}

func loadFlags(fs *flag.FlagSet) error {
	fs.BoolVar(&flagURLOnly, "url", false, "output URL only")
	fs.StringVar(&flagConfigDir, "config", "", "config file directory")
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

func run(args []string) (err error) {
	// lookup bang
	bang, ok := lib.Bangs[args[0]]
	if !ok {
		return fmt.Errorf("not a configured bang: %s", args[0])
	}
	var q string
	if args[1] == "-" {
		if q, err = stdinToString(); err != nil {
			return
		}
	} else {
		q = strings.Join(args[1:], " ")
	}

	url := bang.URL(q)

	if flagURLOnly {
		fmt.Fprintln(os.Stdout, url)
		return
	}

	err = launch(url)
	return
}

func main() {
	fs = flag.NewFlagSet("bang", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", helpText)
		fs.PrintDefaults()
	}

	if err := loadConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		default:
			panic(fmt.Errorf("error loading config: %s", err))
		}
	}

	if err := loadFlags(fs); err != nil {
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
		fmt.Println(lib.ListBangs(lib.Bangs))
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

	if err := run(args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
