# This Makefile intended to be POSIX-compliant (2018 edition with .PHONY target).
#
# .PHONY targets are used by:
#  - task definintions
#  - compilation of Go code (force usage of `go build` to changes detection).
#
# More info:
#  - docs: <https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html>
#  - .PHONY: <https://www.austingroupbugs.net/view.php?id=523>
#
.POSIX:
.SUFFIXES:


#
# PUBLIC MACROS
#

CLI     = opinions
CLIDIR  = ./cmd/opinions
DESTDIR = ./dist
GO      = go
GOFLAGS = 
LDFLAGS = -ldflags "-s -w -X main.AppVersion=$(VERSION)"

LIBSECCOMP_VER  = 2.5.5
LIBSECCOMP_HASH = 248a2c8a4d9b9858aa6baf52712c34afefcf9c9e94b76dce02c1c9aa25fb3375
ZCC     = zig cc
ZAR     = zig ar
ZLD     = zig ld
ZRANLIB = zig ranlib


#
# INTERNAL MACROS
#

LIBSECCOMP_DIR        = /tmp/seccomp
LIBSECCOMP_PREFIX     = $(LIBSECCOMP_DIR)/usr

CLI_CURRENT_VER_TAG   = $$(git tag --points-at HEAD | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)
CLI_LATEST_VERSION    = $$(git tag | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)
CLI_PSEUDOVERSION     = $$(VER="$(CLI_LATEST_VERSION)"; echo "$${VER:-0001.01}")-$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')
CLI_VERSION           = $$(VER="$(CLI_CURRENT_VER_TAG)"; echo "$${VER:-$(CLI_PSEUDOVERSION)}")
MODULE_LATEST_VERSION = $$(git tag | grep "^v" | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)


#
# DEVELOPMENT TASKS
#

.PHONY: all
all: install-dependencies

.PHONY: clean
clean:
	@echo '# Delete bulid directories' >&2
	rm -rf $(DESTDIR) $(LIBSECCOMP_DIR)
	@echo '# Delete libseccomp artifacts' >&2
	# rm libseccomp-$(LIBSECCOMP_VER).tar.gz
	rm libseccomp-*.a

.PHONY: info
info:
	@printf '# OS info: '
	@uname -rsv;
	@echo '# Development dependencies:'
	@$(GO) version || true
	@echo '# Go environment variables:'
	@$(GO) env || true

.PHONY: check
check:
	@echo '# Static analysis' >&2
	$(GO) vet -C $(CLIDIR)

.PHONY: test
test:
	@echo '# Unit tests' >&2
	$(GO) test .

.PHONY: e2e
e2e:
	@echo '# E2E tests of $(DESTDIR)/$(CLI)' >&2
	@printf 'Hacker News\nLemmy\nLobsters\nReddit\n' >test_case.grugbrain
	@if [ -n "$${GITHUB_ACTIONS}" ]; then sed '/Reddit/d' test_case.grugbrain >filtered; mv filtered test_case.grugbrain; fi
	@printf '' >test_case.unknown
	$(DESTDIR)/$(CLI) --version
	$(DESTDIR)/$(CLI) --timeout 10s 'https://grugbrain.dev' | cut -d'	' -f1 | sort -u | diff test_case.grugbrain -
	$(DESTDIR)/$(CLI) --timeout 8500ms 'zażółćjaźńgęślą' | diff test_case.unknown -

.PHONY: build
build:
	@echo '# Build CLI executable: $(DESTDIR)/$(CLI)' >&2
	$(GO) build -C $(CLIDIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/$(CLI)'
	@echo '# Add executable checksum to: $(DESTDIR)/sha256sum.txt' >&2
	cd $(DESTDIR); sha256sum $(CLI) >> sha256sum.txt

.PHONY: unsafe
unsafe:
	@$(MAKE) GOFLAGS='-tags=unsafe' build

.PHONY: dist
dist: opinions-freebsd_amd64 \
 opinions-linux_amd64-hardened \
 opinions-linux_armv7-hardened \
 opinions-linux_arm64-hardened \
 opinions-openbsd_amd64-hardened \
 opinions-windows_amd64.exe

.PHONY: install-dependencies
install-dependencies:
	@echo '# Install CLI dependencies' >&2
	@GOFLAGS='-v -x' $(GO) get -C $(CLIDIR) $(GOFLAGS) .

.PHONY: cli-release
cli-release: check test
	@echo '# Update local branch' >&2
	@git pull --rebase
	@echo '# Create new CLI release tag' >&2
	@VER="$(CLI_LATEST_VERSION)"; printf 'Choose new version number for CLI (calver; >%s): ' "$${VER:-2024.01}"
	@read -r NEW_VERSION; \
		git tag "cli/v$$NEW_VERSION"; \
		git push --tags

.PHONY: module-release
module-release: check test
	@echo '# Update local branch' >&2
	@git pull --rebase
	@echo '# Create new Go module release tag' >&2
	@VER="$(MODULE_LATEST_VERSION)"; printf 'Choose new version number for module (semver; >%s): ' "$${VER:-2.0.0}"
	@read -r NEW_VERSION; \
		git tag "v$$NEW_VERSION"; \
		git push --tags


#
# SUPPORTED EXECUTABLES
#

# this force using `go build` to changes detection in Go project (instead of `make`)
.PHONY: opinions-freebsd_amd64 \
 opinions-linux_amd64-hardened \
 opinions-linux_armv7-hardened \
 opinions-linux_arm64-hardened \
 opinions-openbsd_amd64-hardened \
 opinions-windows_amd64.exe

opinions-freebsd_amd64:
	GOOS=freebsd GOARCH=amd64 $(MAKE) CLI=$@ unsafe

opinions-linux_amd64-hardened:
	GOOS=linux GOARCH=amd64 $(MAKE) CLI=$@ build

opinions-linux_armv7-hardened:
	# FIXME: cross-compilation needs manual compilation of libseccomp:
	#     - libseccomp compilation: <https://github.com/seccomp/libseccomp-golang/blob/main/seccomp_internal.go>
	#     - linking compiled libseccomp: <https://github.com/seccomp/libseccomp-golang/blob/main/.github/workflows/test.yml>
	# GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC="$(ZCC) -target arm-linux-musleabi" $(MAKE) CLI=$@ build
	GOOS=linux GOARCH=arm GOARM=7 $(MAKE) CLI=$@ build

# opinions-linux_arm64-hardened: libseccomp-aarch64.a
# 	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC="$(ZCC) -target aarch64-linux-musl" PKG_CONFIG_PATH="/tmp/seccomp/usr/lib/pkgconfig/" $(MAKE) CLI=$@ build

opinions-linux_arm64-hardened:
	GOOS=linux GOARCH=arm64 $(MAKE) CLI=$@ build

opinions-openbsd_amd64-hardened:
	GOOS=openbsd GOARCH=amd64 $(MAKE) CLI=$@ build

opinions-windows_amd64.exe:
	GOOS=windows GOARCH=amd64 $(MAKE) CLI=$@ unsafe


#
# DEPENDENCIES
#

libseccomp-$(LIBSECCOMP_VER).tar.gz:
	@echo '# Download libseccomp sources' >&2
	wget -q 'https://github.com/seccomp/libseccomp/releases/download/v$(LIBSECCOMP_VER)/libseccomp-$(LIBSECCOMP_VER).tar.gz'
	@if [ "$$(sha256sum libseccomp-$(LIBSECCOMP_VER).tar.gz | cut -d' ' -f1)" != "$(LIBSECCOMP_HASH)" ]; then \
		echo "Checksum mismatch: $$(sha256sum libseccomp-$(LIBSECCOMP_VER).tar.gz | cut -d' ' -f1) != $(LIBSECCOMP_HASH)" >&2; \
		echo "# Delete libseccomp-$(LIBSECCOMP_VER).tar.gz" >&2; \
		rm libseccomp-$(LIBSECCOMP_VER).tar.gz; \
		exit 1; \
	fi

libseccomp-%.a: libseccomp-$(LIBSECCOMP_VER).tar.gz
	@echo '# Build libseccomp' >&2
	mkdir -p $(LIBSECCOMP_DIR)
	cp libseccomp-$(LIBSECCOMP_VER).tar.gz $(LIBSECCOMP_DIR)
	cd $(LIBSECCOMP_DIR); tar -xzf libseccomp-$(LIBSECCOMP_VER).tar.gz
	cd $(LIBSECCOMP_DIR)/libseccomp-$(LIBSECCOMP_VER) ; \
		AR="$(ZAR)" CC="$(ZCC)" CFLAGS="-target $*-linux-musl -static" RANLIB="$(ZRANLIB)" LD="$(ZLD)" ./configure --prefix=/tmp/seccomp/usr --host amd64-pc-linux --build $*-pc-linux && make && make install
	cp $(LIBSECCOMP_PREFIX)/lib/libseccomp.a $@
