# Wholidisuka

Wholidisuka is a command line tool to determine if today is a Japanese holiday.

The exit status is 0 if the date is a holiday and 1 otherwise, but it is 2 if an error occurred.

## Data sources

Wholidisuka refers to [holiday-jp/holiday_jp](https://github.com/holiday-jp/holiday_jp) datasets for Japanese holidays.

You can also use your own calendar that takes priority over holiday_jp.

## Installation

```console
$ go get github.com/quickguard-oss/wholidisuka
```

## Usage

Check if today is a holiday:

```console
$ wholidisuka
```

If Saturday and Sunday are regular holidays:

```console
$ wholidisuka -r 'sat,sun'
```

Use your own calendar:

```console
$ wholidisuka -o ./calendar.yml
```

The format of the calendar is:

```yaml
---
YYYY-MM-DD: '<Description>'  # Holiday
YYYY-MM-DD: ~                # "~ (NULL)" means a business day which overrides the public holiday.
...
```

Shell script example:

```bash
#!/bin/bash

wholidisuka -r 'sat,sun' -o ./calendar.yml

case $? in
  0)  # Holiday
    ...
    ;;
  1)  # Business day
    ...
    ;;
  *)  # Error occurred!
    ...
    ;;
esac
```

## Cache

By default wholidisuka caches [holiday_jp yml](https://github.com/holiday-jp/holiday_jp/blob/master/holidays.yml) for 180 days.

If you want to specify the expiration time, use `-e` option:

```console
$ wholidisuka -e '72h'  # 3 days (in golang "time.ParseDuration()" format)
```

## Caveat

Currently, cache read/write does not support mutual exclusion, so be careful when running multiple processes simultaneously.

## License

MIT
