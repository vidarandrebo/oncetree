proto_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
proto_src := $(shell find . -type f -name '*.proto')
proto_go := $(proto_src:%.proto=%.pb.go)
gorums_go := $(proto_src:%.proto=%_gorums.pb.go)

build:
	go build -o bin/oncetreenode cmd/main.go

test:
	go test -v ./pkg/...

format:
	find . -type f -name "*.go" | xargs gofmt -w

.PHONY: protos

protos: $(proto_go) $(gorums_go)

%.pb.go %_gorums.pb.go : %.proto
	@protoc -I=$(proto_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<
