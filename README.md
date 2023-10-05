# opinions

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

## Deploying

> TODO: this should be declared inside `Makefile`.

```sh
CLI_VERSION="$(date '+%y.%m')"
git tag "v${CLI_VERSION}"
git push --tags
git push origin

GOOS=openbsd GOARCH=amd64 go build -C cmd/ -ldflags="-s -w -X main.AppVersion=$CLI_VERSION" -o "../dist/opinions-$GOARCH"_"$GOOS"
```

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