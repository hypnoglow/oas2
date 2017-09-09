.PHONY: install
install:
	@go install

.PHONY: test
test:
	go test -cover ./...

.PHONY: race
race:
	go test -race -cover ./...

.PHONY: lint
lint: install
	@ gometalinter --concurrency=1 --deadline=60s --vendor ./...
