GOFMT_FILES?=$$(find . -name '*.go' | grep -v pb.go)

default: protoc fmt build

protoc:
	@echo "--> Compiling protobufs"
	protoc *.proto --go_out=plugins=grpc,paths=source_relative:.

fmt:
	@echo "--> Formatting source files"
	gofmt -w $(GOFMT_FILES)

fmt-json:
	@echo "--> Formatting JSON file"
	@jq type bangs.json >/dev/null
	jq -SM '.' bangs.json | awk 'BEGIN{RS="";getline<"-";print>ARGV[1]}' bangs.json

build: protoc fmt
	@echo "--> Building"
	go build -ldflags="-s -w"

.PHONY: default build protoc fmt fmt-json
