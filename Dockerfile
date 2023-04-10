FROM golang:1.19 AS build

WORKDIR /app
ADD . .

RUN go build -o room .

FROM ubuntu:20.04

WORKDIR /app
COPY --from=build /app/room .

EXPOSE 8322

CMD ["./room"]