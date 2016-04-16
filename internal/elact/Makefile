# sampleapp Makefile
SHELL    = /bin/bash
PRG     ?= $(shell basename $$PWD)

ARCH     ?= amd64
ALLARCH  ?= "linux/amd64 linux/386 windows/amd64 darwin/386"

KEY      ?= main-1511

all: build-all

build: bin bin/$(PRG)_linux_$(ARCH) sum

bin/$(PRG)_linux_$(ARCH): *.go
	@echo "*** $@ ***"
	@pushd bin ; \
	gox -osarch="linux/$(ARCH)" -output=$(PRG)_linux_$(ARCH) ../ && popd || { popd ; exit 1 ; }

build-all: bin
	@echo "*** $@ ***"
	@pushd bin ; \
	gox -osarch=$(ALLARCH) ../ && popd || { popd ; exit 1 ; }

bin:
	@[ -d bin ] || mkdir bin

sum:
	@echo "*** $@ ***"
	@pushd bin ; sha256sum $(PRG)* > SHA256SUMS ; popd

run: bin bin/$(PRG)_linux_$(ARCH)
	@echo "*** $@ ***"
	bin/$(PRG)_linux_$(ARCH)

.PHONY: buildall sum bin
