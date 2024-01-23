proto_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
proto_src := $(shell find . -type f -name '*.proto')
proto_go := $(proto_src:%.proto=%.pb.go)
gorums_go := $(proto_src:%.proto=%_gorums.pb.go)



build: protos
	go build -o bin/oncetreenode cmd/oncetreenode/main.go
	go build -o bin/oncetreeclient cmd/oncetreeclient/main.go

test:
	go test ./...

format:
	find . -type f -name "*.go" | xargs gofumpt -w

clean:
	rm -rf bin/

deps:
	go install mvdan.cc/gofumpt@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/relab/gorums/cmd/protoc-gen-gorums@latest

.PHONY: protos

protos: $(proto_go) $(gorums_go)

%.pb.go %_gorums.pb.go : %.proto
	@protoc -I=$(proto_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<
