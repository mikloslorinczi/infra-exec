FROM alpine

WORKDIR /go/src/app

COPY iserver .

COPY /server/server.env .

RUN ["mkdir", "logs"]

CMD ["./iserver"]