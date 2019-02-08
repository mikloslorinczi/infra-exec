FROM golang AS builder

WORKDIR /go/src/github.com/mikloslorinczi/infra-exec/
COPY client ./client
COPY common ./common
COPY executor ./executor
COPY vendor ./vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iclient client/client.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/mikloslorinczi/infra-exec/iclient .
COPY client/client.yaml .

CMD ["./iclient"]