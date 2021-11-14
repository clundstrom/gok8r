FROM golang:1.12.0-alpine3.9

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build main.go

EXPOSE 8000

CMD ["/app/main"]
