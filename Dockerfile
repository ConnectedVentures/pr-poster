FROM golang:alpine as builder
WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
WORKDIR /root

COPY --from=builder /go/bin /main

ENTRYPOINT ["/main/pr-poster"]
