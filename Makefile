GOFMT_FILES?=$$(find . -name '*.go')

default: fmt build

fmt:
	@echo "--> Formatting source files"
	gofmt -w $(GOFMT_FILES)

fmt-json:
	@echo "--> Formatting JSON file"
	@jq type bangs.json >/dev/null
	jq -SM . bangs.json | awk 'BEGIN{RS="";getline<"-";print>ARGV[1]}' bangs.json

build: fmt
	@echo "--> Building"
	go build -ldflags="-s -w"

.PHONY: default build fmt fmt-json
