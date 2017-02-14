# go_mirrors make file

SHELL:=/bin/bash
BASEPATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))


all:print fmt test build
#all:print fmt lint test build

print:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making print<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo SHELL:$(SHELL)
	@echo BASEPATH:$(BASEPATH)
	@echo -e "\n"

fmt:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making fmt<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	go fmt ./...
	@echo -e "\n"

lint:
	-golint ./...

test:

build:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	go build -o ${GOPATH}/bin/hero ${BASEPATH}/hero/main.go
	@echo -e "\n"