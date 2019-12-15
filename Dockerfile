FROM golang:1.13-buster as build

WORKDIR /go/src/novabot
ADD . /go/src/novabot

RUN go get -d -v ./...

RUN go build -o /go/bin/novabot

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/novabot /

ENV "BOT_PREFIX"="!"
ENV "BOT_TOKEN"=""    
ENV "SQL_USER"="",
ENV "SQL_PASS"="",
ENV "SQL_HOST"="",
ENV "SQL_PORT"="",
ENV "SQL_DATABASE"=""

CMD ["/novabot"]
