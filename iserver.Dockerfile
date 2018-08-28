FROM golang AS builder

WORKDIR /go/src/github.com/mikloslorinczi/infra-exec/
COPY server ./server
COPY db ./db
COPY common ./common
COPY vendor ./vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iserver server/server.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/mikloslorinczi/infra-exec/iserver .
COPY server/server.yaml .

CMD ["./iserver"]