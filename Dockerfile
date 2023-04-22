FROM golang:1.19 AS build

WORKDIR /app
ADD . .

RUN go build -o room .

FROM ubuntu:latest

WORKDIR /app

COPY --from=build /app/room .
COPY --from=build /app/configs/openai_config.ini configs/openai_config.ini

EXPOSE 8322

CMD ["./room"]