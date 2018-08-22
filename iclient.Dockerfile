FROM golang

WORKDIR /go/src/app

COPY iclient .

COPY /client/client.env .

RUN ["mkdir", "logs"]

CMD ["./iclient"]