# bang

Browser launcher, heavily inspired by DuckDuckGo's [!bangs][ddg-bangs].

Search for cat pictures on Google Images:

```console
$ bang gi cat pictures
# opens a browser to 'cat pictures' on Google Images
```

Queries can also be piped through stdin if a hyphen is passed in as the argument, and the `-url` flag can be used to print the bang's URL to stdout rather than launching a browser.

```console
$ echo "cat pictures" | bang -url gi -
https://www.google.com/search?tbm=isch&q=cat+pictures
$ echo "reddit cat pictures" | bang -url -
https://www.reddit.com/search?q=cat+pictures
```

The system's URL opener will be used by default, but if set, the `BROWSER` environment variable will be executed with the chosen bang's URL passed as the final argument.

## Config

The CLI looks for a config file named `bangs.(json|yml|yaml|toml|hcl)` in the following locations, in order:

```plain
~/.config/bang/
~/.config/
./ (current directory)
```

Each key of the config file should be the unique `name` of a Bang, with the following properties:

- `description` `(string: <req>)` is a friendly description for the Bang.
- `escape_method` `(int: 0)` defines how the query is escaped prior to it being substituted within the Bang's `format`:
  - `0` - Escapes with `url.QueryEscape`: `cat pictures` &rArr; `cat+pictures`. This is the default method.
  - `1` - Pass-through without escaping: `cat pictures` &rArr; `cat pictures`
  - `2` - Escapes with `url.PathEscape`: `cat pictures` &rArr; `cat%20pictures`
- `format` `(string: <req>)` defines the template used to create the Bang's resulting query string. Use `{{{s}}}` to denote where the query should be substituted.

Example YAML entry for GoDoc:

```json
godoc:
  # try it: bang godoc github.com/travis-g/bang
  description: GoDoc
  escape_method: 1
  format: "https://godoc.org/{{{s}}}"
```

See the [`bang.proto`](bang.proto) file for the Bang object format.

[ddg-bangs]: https://duckduckgo.com/bang
