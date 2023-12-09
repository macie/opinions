# opinions

[![Quality check status](https://github.com/macie/opinions/actions/workflows/check.yml/badge.svg)](https://github.com/macie/opinions/actions/workflows/check.yml)

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
- access to the OS is restricted by application-level sandboxing (with [pledge](https://man.openbsd.org/pledge.2) and [seccomp](https://en.wikipedia.org/wiki/Seccomp)).

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
- `make build` - compile binary from latest commit
- `make unsafe` - compile binary from latest commit without security sandbox
- `make dist` - compile binaries from latest commit for supported OSes (with [proper version number](https://go.dev/doc/modules/version-numbers))
- `make clean` - removes compilation artifacts
- `make cli-release` - tag latest commit as a new release of CLI
- `make info` - print system info (useful for debugging).

## Bugs

Results depend on search engines of underlying social news websites. They
may be different than expected.

## TODO

- add [sandboxing](https://learnbchs.org/pledge.html) for FreeBSD (with
[Capsicum](https://en.wikipedia.org/wiki/Capsicum_(Unix)) - see:
<https://reviews.freebsd.org/rS308432> and
<https://cgit.freebsd.org/src/tree/lib/libcapsicum/capsicum_helpers.h?id=d66f9c86fa3fd8d8f0a56ea96b03ca11f2fac1fb#n104>))
- verify hardened version of linux arm and arm64.

## License

[MIT](./LICENSE)
