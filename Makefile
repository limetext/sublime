all: fmt license generate

test:
	@go test -race $(shell go list ./... | grep -v vendor)
fmt:
	@go fmt $(shell go list ./... | grep -v vendor | grep -v testdata)
license:
	@go run $(GOPATH)/src/github.com/limetext/tasks/gen_license.go
generate:
	@go generate $(shell go list ./... | grep -v /vendor/)

check_fmt:
ifneq ($(shell gofmt -l ./ | grep -v vendor | grep -v testdata),)
	$(error code not fmted, run make fmt. $(shell gofmt -l ./ | grep -v vendor | grep -v testdata))
endif

check_license:
	@go run $(GOPATH)/src/github.com/limetext/tasks/gen_license.go -check

check_generate: generate
ifneq ($(shell git status --porcelain | grep "api/"),)
	$(error generated files are not correct, run make generate. $(shell git status --porcelain | grep "api/"))
endif

tasks:
	go get -d -u github.com/limetext/tasks

glide:
	go get -v -u github.com/Masterminds/glide
	glide install
cover_dep:
	go get -v -u github.com/mattn/goveralls
	go get -v -u github.com/axw/gocov/gocov


travis:
ifeq ($(TRAVIS_OS_NAME),osx)
	brew update
	brew install oniguruma python3
endif

travis_test: export PKG_CONFIG_PATH += $(PWD)/vendor/github.com/limetext/rubex:$(GOPATH)/src/github.com/limetext/rubex
travis_test: test check_generate cover report_cover

cover:
	@echo "mode: count" > coverage.cov; \
	for pkg in $$(go list "./..." | grep -v /vendor/); do \
		go test -covermode=count -coverprofile=tmp.cov $$pkg; \
		sed 1d tmp.cov >> coverage.cov; \
		rm tmp.cov; \
	done

report_cover:
ifeq ($(REPORT_COVERAGE),true)
	$$(go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls -coverprofile=coverage.cov -service=travis-ci
endif
	rm coverage.cov
