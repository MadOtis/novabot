FROM golang:1.13-buster as build

WORKDIR /go/src/novabot
ADD . /go/src/novabot

RUN go get -d -v ./...

RUN go build -o /go/bin/novabot

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/novabot /

ENV "BOT_PREFIX"="!"
ENV "BOT_TOKEN"="NjUzMzk1OTY4ODEyNzc3NDcz.XfF8kQ.dAK1QxwA1d1RII3iocr-eqfmsj4"
ENV "SQL_USER"="novabot",
ENV "SQL_PASS"="novabot",
ENV "SQL_HOST"="localhost",
ENV "SQL_PORT"="3306",
ENV "SQL_DATABASE"="novablack"

CMD ["/novabot"]