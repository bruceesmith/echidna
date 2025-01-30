
File: Makefile

# Change these variables as necessary.
main_package_path = ./
binary_name = ""

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	go test -race -buildvcs -cover ./...

## test/verbose: run all tests with verbose output
.PHONY: test
test/verbose:
	go test -v -race -buildvcs -cover ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go generate ./...
	go build -a -tags osusergo,netgo -ldflags "-s -X 'github.com/bruceesmith/echidna.BuildDate=$(shell date)' -w -extldflags '-static'" -o=/tmp/${binary_name} ${main_package_path}

## run: run the  application
.PHONY: run
run: build
	/tmp/${binary_name}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/${binary_name}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push

## prod: deploy the application to production
.PHONY: prod 
prod: audit no-dirty
	GOOS=linux GOARCH=amd64 go build -a -tags osusergo,netgo -ldflags "-s -X 'github.com/bruceesmith/echidna.BuildDate=$(shell date)' -w -extldflags '-static'" -o ${HOME}/go/bin/${binary_name} ${main_package_path}
	upx -5 ~/go/bin/${binary_name}
	# Include additional deployment steps here...
library: audit
	GOOS=linux GOARCH=amd64 go build ./...
