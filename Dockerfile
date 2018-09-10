FROM golang:alpine
WORKDIR /go/src/pr-poster
COPY . .
RUN go install
ENTRYPOINT ["/go/bin/pr-poster"]
