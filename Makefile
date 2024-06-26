proto_include := $(shell go list -m -f {{.Dir}} github.com/relab/gorums)
proto_src := $(shell find . -type f -name '*.proto')
proto_go := $(proto_src:%.proto=%.pb.go)
gorums_go := $(proto_src:%.proto=%_gorums.pb.go)



.PHONY: build
build: protos logs
	go build -o bin/benchmarkclient cmd/benchmarkclient/main.go
	go build -o bin/benchmarkreplica cmd/benchmarkreplica/main.go

.PHONY: test
test: protos
	go test ./... -race
	go test ./... -v

logs:
	mkdir -p logs

.PHONY: bench
bench:
	go test -run=None ./... -bench=. -benchmem -benchtime=100000x

.PHONY: format
format:
	find . -type f -name "*.go" | xargs gofumpt -w

.PHONY: clean
clean:
	rm -rf bin/
	rm -rf logs
	find protos/ -name "*pb.go" -type f | xargs rm
	go clean -cache -testcache

.PHONY: deps
deps:
	go install mvdan.cc/gofumpt@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/relab/gorums/cmd/protoc-gen-gorums@master

.PHONY: protos
protos: $(proto_go) $(gorums_go) format

%.pb.go %_gorums.pb.go : %.proto
	protoc -I=$(proto_include):. \
		--go_out=paths=source_relative:. \
		--gorums_out=paths=source_relative:. \
		$<


.PHONY: publish
publish: build
	ssh ssh4.ux.uis.no mkdir -p oncetree
	rsync -aP bin sh ssh4.ux.uis.no:oncetree/
	ssh ssh4.ux.uis.no ./oncetree/sh/publish_bbchain.sh

.PHONY: clean_remote
clean_remote:
	-ssh ssh4.ux.uis.no ./oncetree/sh/clean.sh

.PHONY: bench_remote
bench_remote: stop_remote clean_remote publish
	ssh ssh4.ux.uis.no ./oncetree/sh/bench.sh
	ssh ssh4.ux.uis.no ./oncetree/sh/clients.sh

.PHONY: stop_remote
stop_remote:
	-ssh ssh4.ux.uis.no ./oncetree/sh/stop.sh

.PHONY: get_logs
get_logs:
	-rm -r writes reads replicas
	ssh ssh4.ux.uis.no ./oncetree/sh/get_logs.sh
	rsync -aP ssh4.ux.uis.no:oncetree/writes ./
	rsync -aP ssh4.ux.uis.no:oncetree/reads ./
	rsync -aP ssh4.ux.uis.no:oncetree/replicas ./
