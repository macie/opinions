# opinions

[![Go Reference](https://pkg.go.dev/badge/github.com/macie/opinions.svg)](https://pkg.go.dev/github.com/macie/opinions)
[![Quality check status](https://github.com/macie/opinions/actions/workflows/check.yml/badge.svg)](https://github.com/macie/opinions/actions/workflows/check.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/macie/opinions)](https://goreportcard.com/report/github.com/macie/opinions)

Find opinions about a given phrase (also URL) on social news websites:

- _[Lemmy](https://en.wikipedia.org/wiki/Lemmy_(social_network))_
- _[Lobsters](https://lobste.rs/about)_
- _[Hacker News](https://en.wikipedia.org/wiki/Hacker_News)_
- _[Reddit](https://en.wikipedia.org/wiki/Reddit)_.

It can be used to [include discussions on static websites/blogs](#static-site-generators).

_opinions_ is a command-line replacement of [discu.eu](https://discu.eu/) service.
It directly calls search engines on underlying websites.

Application is developed with a security-first approach:

- functionality is limited by design
- access to the OS is restricted by [application-level sandboxing](#security-hardening).

## Usage

```sh
$ opinions 'https://grugbrain.dev'
Lemmy	https://lemmy.world/post/5189937	The Grug Brained Developer - A layman's guide to thinking like the self-aware smol brained	https://grugbrain.dev/
Lemmy	https://lemmy.world/post/750451	The Grug Brained Developer	https://grugbrain.dev/
Lemmy	https://lemmy.world/post/685510	The Grug Brained Developer	https://grugbrain.dev/
Hacker News	https://news.ycombinator.com/item?id=31840331	The Grug Brained Developer	https://grugbrain.dev/
Hacker News	https://news.ycombinator.com/item?id=38076886	The Grug Brained Developer (2022)	https://grugbrain.dev/
Lobsters	https://lobste.rs/s/ifaar4/grug_brained_developer	The Grug Brained Developer	https://grugbrain.dev/
Lobsters	https://lobste.rs/s/pmpc9v/grug_brained_developer	The Grug Brained Developer	http://grugbrain.dev
Reddit	https://reddit.com/r/programming/comments/16gt8w4/the_grug_brained_developer/	The grug brained developer	https://grugbrain.dev
Reddit	https://reddit.com/r/brdev/comments/14jhm17/the_grug_brained_developer/	The Grug Brained Developer	https://grugbrain.dev
```

The result is printed to stdout as rows in format: `<service_name><tab><URL><tab><title><tab><source_domain>`.

Websites are queried with User-Agent: `opinions/<version_number> (<os>; +https://github.com/macie/opinions)`.

### Static Site Generators

_opinions_ can be used to extend static websites with comments-like feature. You
can search for active discussions about your article, filter out false positive
results (eg. get only articles from your domain) and generate HTML links with
standard Unix commands:

```sh
opinions 'Grug Brained' | grep 'grugbrain.dev' | awk -F '\t' '
BEGIN { print "<ul>" }
{ print "  <li><a href=\""$2"\" title=\"["$1"] "$3"\">"$1"</a></li>" }
END { print "</ul>" }
'
```

## Installation

Download [latest stable release from GitHub](https://github.com/macie/opinions/releases/latest) .

You can also build it manually with commands: `make && make build` (or without
sandboxing: `make && make unsafe`).

## Development

Use `make` (GNU or BSD):

- `make` - install dependencies
- `make test` - runs test
- `make e2e` - runs e2e tests for CLI
- `make check` - static code analysis
- `make build` - compile binary from the latest commit
- `make unsafe` - compile binary from the latest commit without the security sandbox
- `make dist` - compile binaries from the latest commit for all supported OSes
- `make clean` - removes compilation artifacts
- `make cli-release` - tag the latest commit as a new release of CLI
- `make module-release` - tag the latest commit as a new release of Go module
- `make info` - print system info (useful for debugging).

### Versioning

At the begining, only CLI was released with _[semantic versioning](https://semver.org/)_ scheme (commits marked by tags `v1.0.0`-`v1.5.1`).

Currently repo contains CLI and Go module which can be developed with different
pace. Commits with versions are tagged with:
- `v2.X.X` (_[semantic versioning](https://semver.org/)_) - versions of Go module
- `cli/vYYYY.0M.MICRO` (_[calendar versioning](https://calver.org/)_) - versions of command-line utility.

### Security hardening

On modern Linuxes and OpenBSD, CLI application has restricted access to kernel
calls with [seccomp](https://en.wikipedia.org/wiki/Seccomp) and [pledge](https://man.openbsd.org/pledge.2).

### TODO

- add [sandboxing](https://learnbchs.org/pledge.html) for FreeBSD (with
[Capsicum](https://en.wikipedia.org/wiki/Capsicum_(Unix)) - see:
<https://reviews.freebsd.org/rS308432> and
<https://cgit.freebsd.org/src/tree/lib/libcapsicum/capsicum_helpers.h?id=d66f9c86fa3fd8d8f0a56ea96b03ca11f2fac1fb#n104>))
- verify hardened version of linux arm and arm64.

## Bugs

Results depend on search engines of underlying social news websites. They
may be different than expected.

## License

[MIT](./LICENSE)
