# set the shell to bash in case some environments use sh
SHELL:=/bin/bash

# VERSION is the version of the binary.
VERSION:=$(shell git describe --tags --always)
REPO = $(shell sh -c "git ls-remote --get-url origin | cut -f 2 -d @" | awk -F ".git" '{print $$1}' | sed 's/:/\//')

IMAGE_PREFIX = hub.rms.evolutiongaming.com/prodops

# Determine the arch/os
export XC_OS=linux
export XC_ARCH=amd64

ARCH:=${XC_OS}_${XC_ARCH}
export ARCH

ifeq (${IMAGE_ORG}, )
  IMAGE_ORG = hub.rms.evolutiongaming.com/prod-ops/aws-reporting
  export IMAGE_ORG
endif

# Specify the date of build
DBUILD_DATE=$(shell date -u +'%Y%m%dT%H%M%SZ')

export DBUILD_ARGS=--build-arg DBUILD_DATE=${DBUILD_DATE} --build-arg ARCH=${ARCH}


# -composite: avoid "literal copies lock value from fakePtr"
.PHONY: vet
vet:
	go list ./... | grep -v "./vendor/*" | xargs go vet -composites

.PHONY: fmt
fmt:
	find . -type f -name "*.go" | grep -v "./vendor/*" | xargs gofmt -s -w -l

.PHONY: lint
lint:
	golangci-lint run

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: build.common
build.common: version

.PHONY: clean
clean:
	@echo '--> Cleaning directory...'
	rm -rf bin
	@echo '--> Done cleaning.'

.PHONY: test
test:
	docker build -t ${IMAGE_ORG}:test -f build/test/Dockerfile .

.PHONY: update-deps
update-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif