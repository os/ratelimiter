FROM golang:1.18.1-buster
ENTRYPOINT go test -v ./...
