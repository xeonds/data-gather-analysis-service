NAME=data-gather-analysis
BINDIR=build
VERSION=1.0.0
BUILDTIME=$(shell date -u)
GOBUILD=go build
GOFLAGS=-ldflags '-s -w -X "main.version=$(VERSION)" -X "main.buildTime=$(BUILDTIME)"'

all: gen linux-amd64

linux-amd64: 
	GOOS=linux GOARCH=amd64 cd analysis && $(GOBUILD) -o ../$(BINDIR)/$(NAME)-analysis-$@
	GOOS=linux GOARCH=amd64 cd gather && $(GOBUILD) -o ../$(BINDIR)/$(NAME)-gather-$@

web:
	cd display && pnpm i && pnpm run build --outDir=../$(BINDIR)/dist --emptyOutDir