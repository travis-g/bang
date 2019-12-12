# bang

CLI for quick browser launching, heavily inspired by DuckDuckGo's [!bangs][ddg-bangs].

Queries can also be piped through Stdin.

```console
$ echo "cat pictures" | bang -url gi -
https://www.google.com/search?tbm=isch&q=cat+pictures
```

Supports launching custom programs as subshells using the `BROWSER` environment variable:

```console
$ export BROWSER="chromium-browser --incognito"
$ bang gi cat pictures
# open the search for cat pictures in an incognito Chromium session
```

[ddg-bangs]: https://duckduckgo.com/bang
