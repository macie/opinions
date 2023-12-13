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


#
# INTERNAL MACROS
#

CURRENT_VER_TAG = $$(git tag --points-at HEAD | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)
LATEST_VERSION  = $$(git tag | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)
PSEUDOVERSION   = $$(VER="$(LATEST_VERSION)"; echo "$${VER:-0.0.0}")-$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')
VERSION         = $$(VER="$(CURRENT_VER_TAG)"; echo "$${VER:-$(PSEUDOVERSION)}")


#
# DEVELOPMENT TASKS
#

.PHONY: all
all: install-dependencies

.PHONY: clean
clean:
	@echo '# Delete bulid directory' >&2
	rm -rf $(DESTDIR)

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

.PHONY: unsafe
unsafe:
	@$(MAKE) GOFLAGS='-tags=unsafe' build

.PHONY: dist
dist:
	@echo '# Create CLI executables in $(DESTDIR)' >&2
	# hardened
	GOOS=openbsd GOARCH=amd64 $(GO) build -C $(CLIDIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/opinions-openbsd_amd64-hardened'
	GOOS=linux GOARCH=amd64 $(GO) build -C $(CLIDIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/opinions-linux_amd64-hardened'
	# without sandbox
	GOFLAGS='-tags=unsafe' GOOS=linux GOARCH=arm GOARM=7 $(GO) build -C $(CLIDIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/opinions-linux_armv7'
	GOFLAGS='-tags=unsafe' GOOS=linux GOARCH=arm64 $(GO) build -C $(CLIDIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/opinions-linux_arm64'
	GOFLAGS='-tags=unsafe' GOOS=freebsd GOARCH=amd64 $(GO) build -C $(CLIDIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/opinions-freebsd_amd64'
	GOFLAGS='-tags=unsafe' GOOS=windows GOARCH=amd64 $(GO) build -C $(CLIDIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/opinions-windows_amd64.exe'

	@echo '# Create checksums' >&2
	@cd $(DESTDIR); sha256sum * >sha256sum.txt

.PHONY: install-dependencies
install-dependencies:
	@echo '# Install CLI dependencies' >&2
	@GOFLAGS='-v -x' $(GO) get -C $(CLIDIR) $(GOFLAGS) .

.PHONY: cli-release
cli-release: check test
	@echo '# Update local branch' >&2
	@git pull --rebase
	@echo '# Create new release tag' >&2
	@VER="$(LATEST_VERSION)"; printf 'Choose new version number (>%s): ' "$${VER:-0.0.0}"
	@read -r NEW_VERSION; \
		echo "New tag: $${NEW_VERSION}"
		git tag "v$$NEW_VERSION"; \
		git push --tags
