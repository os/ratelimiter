FROM golang:1.15-buster AS builder
WORKDIR /go/src/github.com/os/ratelimiter
ADD . .
RUN go install

FROM debian:buster
COPY --from=builder /go/bin/ratelimiter /bin
EXPOSE 8080
ENTRYPOINT ["ratelimiter"]
