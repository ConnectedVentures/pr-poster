FROM circleci/golang:1.10 as builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main

FROM scratch
WORKDIR /root

COPY --from=builder /main /main

ENTRYPOINT ["/main/pr-poster"]
