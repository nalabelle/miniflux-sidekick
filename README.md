# Miniflux Sidekick

This is a sidekick container that runs alongside [Miniflux](https://miniflux.app) which is an Open Source RSS feed reader written in Go. You can check out the source code on [GitHub](https://github.com/miniflux/miniflux).

The goal is to support so called [Killfiles](https://en.wikipedia.org/wiki/Kill_file) to filter out items you don't want to see. You can think of it as an ad-blocker for your feed reader. The items are not deleted, they are just marked as *read* so you can still find all items in your feed reader if you have to.

## Features

- Supports a subset of so called UseNet killfiles rules
- Supports local killfiles on disk

## Supported Rules

The general format of the `killfile` is:

```
ignore-article "<feed>" "<filterexpr>"
```

### `<feed>`

This contains the URL of the feed that should be matched. It fuzzy matches the URL so if you only have one feed just use the base URL of the site. Example: `https://example.com` if the feed is on `https://example.com/rss/atom.xml`. A wildcard selector of `*` is also supported instead of the URL.

Alternately, you may specify a comma-separated list of categories whose feeds should be matched by starting the value with `category:`. Example: `category:Photos`.

### `<filterexpr>` Filter Expressions

From the [available rule set](https://newsboat.org/releases/2.15/docs/newsboat.html#_filter_language) and attributes (`Table 5. Available Attributes`) only a small subset are supported right now. These should cover most use cases already though.

**Attributes**

- `title`
- `content`
- `author`
- `tag`

**Comparison Operators**

- `=~`: test whether regular expression matches
- `!~`: logical negation of the `=~` operator
- `#`: contains; this operator matches if a word is contained in a list of space-separated words (useful for matching tags, see below)
- `!#`: contains not; the negation of the `#` operator



### Example

Here's an example of what a `killfile` could look like with these rules. 

This one marks all feed items as read that have a `[Sponsor]` string in the title.
```
ignore-article "https://www.example.com" "title =~ \[Sponsor\]"
```

This one filters out all feed items that have the word `Lunar` OR `moon` in there.
```
ignore-article "https://xkcd.com/atom.xml" "title # Lunar,Moon"
```

This one filters out all feed items that have the word `lunar` OR `moon` in there and is case insensitive.
```
ignore-article "https://xkcd.com/atom.xml" "title =~ (?i)(lunAR|MOON)"
```

This one marks read all feed items without an image in feeds assigned a category of `Photos`.
```
ignore-article "category:Photos" "content !~ (?i)(img src=)"

```

### Testing rules

There are tests in `filter/` that can be used to easily test rules or add new comparison operators.

## Deploy

There are the environment variables that can be set. You can set `MF_KILLFILE_PATH="~/path/to/killfile"`. `MF_API_ENDPOINT` is your Miniflux url and `MF_API_KEY` is your Miniflux credentials. `MF_REFRESH_INTERVAL` is a cron expression that controls how often the rules are evaluated on Miniflux entries.

```
export MF_API_ENDPOINT=https://rss.example.com
export MF_API_KEY=abc123
export MF_KILLFILE_PATH=/path/to/killfile
export MF_REFRESH_INTERVAL="0 30 * * * *"
```

## See Also

Miniflux v2.0.25 added built-in support for [filtering rules](https://miniflux.app/docs/rules.html#filtering-rules).
