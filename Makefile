tempura: *.go */*.go go.mod
	gofmt -w *.go */*.go
	go build -o tempura main.go

.PHONY: init
init:
	go mod init tempura

.PHONY: test
test:
	gofmt -w *.go
	go test
