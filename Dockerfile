FROM golang:1.17.2 as builder

WORKDIR /go/src/github.com/sidharthjs/todo
COPY . /go/src/github.com/sidharthjs/todo

RUN CGO_ENABLED=0 go build -mod vendor

FROM alpine:3.14.2 AS final

WORKDIR /app
COPY --from=builder /go/src/github.com/sidharthjs/todo/todo /app
COPY --from=builder /go/src/github.com/sidharthjs/todo/db /app/db

CMD ["./todo"]