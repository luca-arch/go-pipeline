SHELL:=/bin/sh
GCI_LINT:=v1.47.2

ci: lint test test-functional                                  ### Pre-push hook
.PHONY: ci

cover: test                                         ### Run tests, with coverage
	go tool cover -html=cover.out -o cover.html
.PHONY: cover

dist-native:                               ### Compile binaries (native platform)
	go build -o pipeline_native ./cmd
.PHONY: dist

dist-linux: docker-all                              ### Compile binaries (Linux)
	docker run --rm --entrypoint cat gopipeline:linux /bin/pipeline > pipeline_linux_amd64
	chmod +x pipeline_linux_amd64
.PHONY: dist-linux

dist-osx: docker-all                                  ### Compile binaries (OSX)
	docker run --rm --entrypoint cat gopipeline:darwin /bin/pipeline > pipeline_darwin_amd64
	docker run --rm --entrypoint cat gopipeline:m1 /bin/pipeline > pipeline_darwin_arm64
	chmod +x pipeline_darwin_amd64 pipeline_darwin_arm64
.PHONY: dist-osx

docker-alpine:                                     ### Build alpine docker image
	docker build --tag gopipeline:alpine .
.PHONY: docker

docker-all: docker-alpine                            ### Build all docker images
	docker build --build-arg GOARCH=amd64 --build-arg GOOS=darwin --tag gopipeline:darwin .
	docker build --build-arg GOARCH=arm64 --build-arg GOOS=darwin --tag gopipeline:m1 .
	docker build --build-arg GOARCH=amd64 --build-arg GOOS=linux --tag gopipeline:linux .
.PHONY: docker-all

fix:                                                 ### Fix gomft and goimports
	go fmt ./...
	goimports -w ./
	go install mvdan.cc/gofumpt@latest
	gofumpt -e -l -w .
.PHONY: fix

help:                                               ### Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

lint:                                                      ### Run golangci-lint
	docker run --rm -v $(PWD):/mnt -w /mnt golangci/golangci-lint:$(GCI_LINT) golangci-lint run
	echo -e '\033[1mgolangci-lint passed!\033[0m';
	docker run --rm -v $(PWD):/data cytopia/goimports:latest --ci .
	echo -e '\033[1mgoimports passed!\033[0m';
.PHONY: lint

test-functional: docker-alpine                          ### Run functional tests
	testdata/run-tests.sh
	echo -e '\033[1mfunctional tests passed!\033[0m';
.PHONY: test-functional

test:                                                         ### Run unit tests
	go test -v -cover -race -coverprofile cover.out ./...
	echo -e '\033[1mtests passed!\033[0m';
.PHONY: test
