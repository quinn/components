mod agents

set shell := ["bash", "-lc"]

export PORT := `ports web`

run:
	go run ./cmd/web

dev:
	watchexec --restart \
		--watch cmd \
		--watch internal \
		--watch justfile \
		--watch go.mod \
		--watch go.sum \
		--exts go,css,js -- go run ./cmd/web

test:
	go test ./...

fmt:
	go fmt ./...
