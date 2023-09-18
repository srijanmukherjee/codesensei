.PHONY: test

server:
	go build -o bin/server cmd/server/main.go

worker:
	go build -o bin/worker cmd/worker/main.go

test:
	go test ./test/./... -v