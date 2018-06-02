.PHONY: install
install:
	@go install

.PHONY: test
test:
	go test -cover ./...

.PHONY: test-race
test-race:
	go test -race -cover ./...

.PHONY: e2e
e2e:
	go test -tags e2e ./e2e/...

.PHONY: lint
lint: install
	@gometalinter ./...
