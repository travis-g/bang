GOFMT_FILES?=$$(find . -name '*.go')

default: build

fmt:
	@echo "--> Formatting source files"
	gofmt -w $(GOFMT_FILES)
	jq -SM . bangs.json | awk 'BEGIN{RS="";getline<"-";print>ARGV[1]}' bangs.json

bindata:
	@echo "--> Generating static files"
	go-bindata bangs.json

build: bindata fmt
	@echo "--> Building"
	go build -ldflags="-s -w"

.PHONY: default build fmt bindata
