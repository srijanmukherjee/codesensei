.PHONY: test isolate

server:
	go build -o bin/server cmd/server/main.go

worker:
	go build -o bin/worker cmd/worker/main.go

dev: server
	docker compose -f docker-compose.dev.yml up db --no-recreate -d
	./bin/server

install: install-isolate
	echo "NOT IMPLEMENTED"

test:
	go test ./test/./... -v

install-isolate:
	cp isolate.cf isolate/default.cf
	cd isolate && $(MAKE) install