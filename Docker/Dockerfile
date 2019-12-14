FROM golang:1.13-buster as build

WORKDIR /go/src/novabot
ADD . /go/src/novabot

RUN go get -d -v ./...

RUN go build -o /go/bin/novabot

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/novabot /

ENV "BOT_PREFIX"="!"
ENV "BOT_TOKEN"=""    
ENV "SQL_USER"="novabot",
ENV "SQL_PASS"="novabot",
ENV "SQL_HOST"="mariadb",
ENV "SQL_PORT"="3306",
ENV "SQL_DATABASE"="novablack"

CMD ["/novabot"]
