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

.PHONY: docker-run
docker-run:
	docker rm -f tempura
	docker build \
		-f Dockerfile \
		--target tempura-builder \
		./ -t tempura
	docker run -itd \
		--restart=always \
		--network proxy \
		--name tempura tempura \
		/tempura/tempura --serve 8080
