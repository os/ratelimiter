FROM golang:1.15-buster
ENTRYPOINT go test -v ./...