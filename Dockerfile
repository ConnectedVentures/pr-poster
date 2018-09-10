FROM golang:alpine
WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

ENTRYPOINT ["/go/src/app/pr-poster"]
