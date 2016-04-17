default: test

test:
	@go test -race $(shell go list ./... | grep -v vendor)
fmt:
	@go fmt $(shell go list ./... | grep -v vendor | grep -v testdata)
license:
	@go run gen_license.go ./
generate:
	@go generate $(shell go list ./... | grep -v /vendor/)

check_fmt:
ifneq ($(shell gofmt -l ./ | grep -v vendor | grep -v testdata),)
	$(error code not fmted, run make fmt. $(shell gofmt -l ./ | grep -v vendor | grep -v testdata))
endif

check_license:
ifneq ($(shell go run gen_license.go ./),)
	$(error license is not added to all files, run make license)
endif

check_generate: config_oniguruma generate
ifneq ($(shell git status --porcelain),)
	$(error generated files are not correct, run make generate)
endif

glide:
	go get -v -u github.com/Masterminds/glide
	glide install

travis:
ifeq ($(TRAVIS_OS_NAME),osx)
	brew update
	brew install oniguruma python3
endif

config_oniguruma: export PKG_CONFIG_PATH += $(PWD)/vendor/github.com/limetext/rubex:$(GOPATH)/src/github.com/limetext/rubex
travis_test: config_oniguruma test
