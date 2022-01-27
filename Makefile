.PHONY: prepare
prepare:
	@go mod vendor

.PHONY: counter
counter:
	@CGO_ENABLED=0 go build -v -ldflags "-s -w" -o counter cmd/counter/main.go

.PHONY: build
build: counter

.PHONY: clean
clean:
	@rm -rf counter