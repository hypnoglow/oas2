.PHONY: install
install:
	go install

.PHONY: test
test:
	go test -cover ./...