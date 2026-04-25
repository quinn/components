set shell := ["bash", "-lc"]

export PORT := `ports web`

run:
	go run ./cmd/dev

dev:
	watchexec --restart \
		--watch cmd \
		--watch internal \
		--watch ssg \
		--watch components \
		--watch css \
		--watch justfile \
		--watch go.mod \
		--watch go.sum \
		--exts go,css,js -- go run ./cmd/dev

build out="./dist":
	go run ./cmd/build -o {{out}}

test:
	go test ./...

fmt:
	go fmt ./...
