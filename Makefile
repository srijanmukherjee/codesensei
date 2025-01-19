worker-dev:
	go run worker/main.go

worker: worker/main.go
	go build -o bin/worker worker/main.go

test:
	sudo -E env "PATH=$$PATH" "PROJECT_ROOT=$$(pwd)" go test -v ./...

.PHONY: test