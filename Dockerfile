FROM golang:1.16-alpine AS tempura-builder

USER root
WORKDIR /tempura
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add make

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./ ./
RUN rm tempura; make

FROM scratch

WORKDIR /tempura

COPY --from=tempura-builder /tempura/tempura .
CMD ["./tempura"]
