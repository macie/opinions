# opinions

[![Quality check status](https://github.com/macie/opinions/actions/workflows/check.yml/badge.svg)](https://github.com/macie/opinions/actions/workflows/check.yml)

Find opinions about given phrase (also URL) on social news websites:
_Hacker News_ and _Lobsters_. It can be used for including discussions on static
websites/blogs.

_opinions_ is command-line replacement of [discu.eu](https://discu.eu/) service.
It directly calls search engines on underlying websites.

Application is developed with security-first approach:

- functionality is limited by design
- access to the OS is restricted by application-level sandboxing (currently OpenBSD only).

## Usage

```sh
$ opinions 'https://grugbrain.dev'
Lobsters	https://lobste.rs/s/ifaar4/grug_brained_developer	The Grug Brained Developer	https://grugbrain.dev/
Lobsters	https://lobste.rs/s/pmpc9v/grug_brained_developer	The Grug Brained Developer	http://grugbrain.dev
Hacker News	https://news.ycombinator.com/item?id=31840331	The Grug Brained Developer	https://grugbrain.dev/
```

The result is printed to stdout as rows in format: `<service_name><tab><URL><tab><title><tab><source_domain>`.

Websites are queried with User-Agent: `opinions/<version_number> (<os>; +https://github.com/macie/opinions)`.

## Installation

Download [latest stable release from GitHub](https://github.com/macie/opinions/releases/latest) .

You can also build it manually with commands: `make && make build` (or without
sandboxing: `make && make unsafe`).

## Development

Use `make` (GNU or BSD):

- `make` - install dependencies
- `make test` - runs test
- `make check` - static code analysis
- `make build` - compile binaries from latest commit for supported OSes (with [proper version number](https://go.dev/doc/modules/version-numbers))
- `make unsafe` - compile binaries from latest commit without security sandbox
- `make release` - mark latest commit with choosen version tag and compile binaries for supported OSes
- `make clean` - removes compilation artifacts
- `make info` - print system info (useful for debugging).

## Bugs

Results depends on search engines of underlying social news websites. They
may be different than expected.

## TODO

Add [sandboxing](https://learnbchs.org/pledge.html) for other OSes:

- Linux: [seccomp](https://en.wikipedia.org/wiki/Seccomp) (see:
<https://github.com/stephane-martin/skewer/blob/master/sys/scomp/seccomp.go> and
<https://blog.heroku.com/applying-seccomp-filters-on-go-binaries>)
- FreeBSD: [Capsicum](https://en.wikipedia.org/wiki/Capsicum_(Unix)) (see:
<https://reviews.freebsd.org/rS308432> and
<https://cgit.freebsd.org/src/tree/lib/libcapsicum/capsicum_helpers.h?id=d66f9c86fa3fd8d8f0a56ea96b03ca11f2fac1fb#n104>)

## License

[MIT](./LICENSE)
