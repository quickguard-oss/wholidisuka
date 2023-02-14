# Wholidisuka

Wholidisuka is a command line tool that checks whether today is a holiday in Japan.

When the date is a holiday, the exit status will be 0, and when it is not a holiday, the exit status will be 1.    
However, if an error occurs, the exit status will be 2.

## Data sources

Wholidisuka uses the [holiday-jp/holiday_jp](https://github.com/holiday-jp/holiday_jp) dataset as a source of information for Japanese holidays.

Alternatively, you can opt to provide their own calendar, which takes precedence over the holiday_jp dataset.

## Installation

```console
$ go install github.com/quickguard-oss/wholidisuka/cmd/wholidisuka@latest
```

## Usage

Check if today is a holiday:

```console
$ wholidisuka
```

If Saturday and Sunday are regular holidays:

```console
$ wholidisuka -r 'sat' -r 'sun'
```

Use your own calendar:

```console
$ wholidisuka -o ./calendar.yml
```

The format of the calendar is:

```yaml
---
YYYY-MM-DD: '<Description>'  # Holiday
YYYY-MM-DD: ~                # "~ (NULL)" refers to a business day that takes precedence over public holidays.
...
```

Shell script example:

```bash
#!/bin/bash

wholidisuka -r 'sat' -r 'sun' -o ./calendar.yml

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

By default, wholidisuka caches the [holiday-jp YAML](https://github.com/holiday-jp/holiday_jp/blob/master/holidays.yml) data for 180 days.

If you wish to specify a different expiration date for the cache, you can use the `-e` option:

```console
$ wholidisuka -e '72h'  # 3 days (in golang "time.ParseDuration()" format)
```

## Caveat

Note that wholidisuka does not support mutual exclusion for cache read/write operations.

Running multiple processes simultaneously may lead to corrupted cache data.

## License

MIT
