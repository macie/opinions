.POSIX:
.SUFFIXES:

CLIDIR = ./cmd/opinions

# MAIN TARGETS

all: install-dependencies

clean:
	@echo '# Delete binaries: rm -rf ./dist' >&2
	@rm -rf ./dist

info:
	@printf '# OS info: '
	@uname -rsv;
	@echo '# Development dependencies:'
	@go version || true
	@echo '# Go environment variables:'
	@go env || true

check:
	@echo '# Static analysis: go vet' >&2
	@go vet -C $(CLIDIR)
	
test:
	@echo '# Unit tests: go test .' >&2
	@go test .

e2e:
	@echo '# E2E tests of ./dist/opinions' >&2
	@printf 'Hacker News\nLemmy\nLobsters\nReddit\n' >test_case.grugbrain
	@if [ -n "$${GITHUB_ACTIONS}" ]; then sed '/Reddit/d' test_case.grugbrain >filtered; mv filtered test_case.grugbrain; fi
	@printf '' >test_case.unknown
	./dist/opinions --version
	./dist/opinions --timeout 10s 'https://grugbrain.dev' | cut -d'	' -f1 | sort -u | diff test_case.grugbrain -
	./dist/opinions --timeout 8500ms 'zażółćjaźńgęślą' | diff test_case.unknown -

build: *.go
	@echo '# Create release binary: ./dist/opinions' >&2
	@CURRENT_VER_TAG="$$(git tag --points-at HEAD | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		PREV_VER_TAG="$$(git tag | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		CURRENT_COMMIT_TAG="$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')"; \
		PSEUDOVERSION="$${PREV_VER_TAG:-0.0.0}-$$CURRENT_COMMIT_TAG"; \
		VERSION="$${CURRENT_VER_TAG:-$$PSEUDOVERSION}"; \
		go build -C $(CLIDIR) -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/opinions'

unsafe: *.go
	@echo '# Create release binary without sandbox in ./dist/opinions' >&2
	@CURRENT_VER_TAG="$$(git tag --points-at HEAD | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		PREV_VER_TAG="$$(git tag | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		CURRENT_COMMIT_TAG="$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')"; \
		PSEUDOVERSION="$${PREV_VER_TAG:-0.0.0}-$$CURRENT_COMMIT_TAG"; \
		VERSION="$${CURRENT_VER_TAG:-$$PSEUDOVERSION}"; \
		go build -C $(CLIDIR) -tags unsafe -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/opinions'

dist: *.go
	@echo '# Create release binaries in ./dist' >&2
	@CURRENT_VER_TAG="$$(git tag --points-at HEAD | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		PREV_VER_TAG="$$(git tag | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		CURRENT_COMMIT_TAG="$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')"; \
		PSEUDOVERSION="$${PREV_VER_TAG:-0.0.0}-$$CURRENT_COMMIT_TAG"; \
		VERSION="$${CURRENT_VER_TAG:-$$PSEUDOVERSION}"; \
		# hardened \
		GOOS=openbsd GOARCH=amd64 go build -C $(CLIDIR) -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/opinions-openbsd_amd64-hardened'; \
		GOOS=linux GOARCH=amd64 go build -C $(CLIDIR) -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/opinions-linux_amd64-hardened'; \
		# without sandbox \
		GOOS=linux GOARCH=arm GOARM=7 go build -C $(CLIDIR) -tags unsafe -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/opinions-linux_armv7'; \
		GOOS=linux GOARCH=arm64 go build -C $(CLIDIR) -tags unsafe -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/opinions-linux_arm64'; \
		GOOS=freebsd GOARCH=amd64 go build -C $(CLIDIR) -tags unsafe -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/opinions-freebsd_amd64'; \
		GOOS=windows GOARCH=amd64 go build -C $(CLIDIR) -tags unsafe -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/opinions-windows_amd64.exe'; \

	@echo '# Create binaries checksum' >&2
	@cd ./dist; sha256sum * >sha256sum.txt

install-dependencies:
	@echo '# Install CLI dependencies:' >&2
	@go get -C $(CLIDIR) -v -x .

cli-release: check test
	@echo '# Update local branch' >&2
	@git pull --rebase
	@echo '# Create new release tag' >&2
	@PREV_VER_TAG=$$(git tag | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1); \
		printf 'Choose new version number (>%s): ' "$${PREV_VER_TAG:-0.0.0}"
	@read -r VERSION; \
		git tag "v$$VERSION"; \
		git push --tags
