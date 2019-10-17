GOFMT_FILES?=$$(find . -name '*.go' | grep -v pb.go)

.PHONY: default
default: fmt build

.PHONY: protoc
protoc:
	@echo "--> Compiling protobufs"
	protoc *.proto --go_out=plugins=grpc,paths=source_relative:.

.PHONY: fmt
fmt:
	@echo "--> Formatting source files"
	gofmt -w $(GOFMT_FILES)

.PHONY: fmt-json
fmt-json:
	@echo "--> Formatting JSON file"
	@jq type bangs.json >/dev/null
	jq -SM '.' bangs.json | awk 'BEGIN{RS="";getline<"-";print>ARGV[1]}' bangs.json

.PHONY: build
build: fmt
	@echo "--> Building"
	go build -ldflags="-s -w"
