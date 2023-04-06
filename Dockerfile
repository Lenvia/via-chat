FROM golang:1.19 AS build

WORKDIR /app
ADD . .

RUN go get github.com/go-sql-driver/mysql
RUN go build -o app .

FROM ubuntu:20.04

WORKDIR /app
COPY --from=build /app/app .

EXPOSE 8322

CMD ["./app"]