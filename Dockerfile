# server
FROM golang:1.19 AS build

WORKDIR /app
ADD server .

RUN go build -o room .

# web
FROM node:18.0 AS node
WORKDIR /app
# 拷贝配置文件
COPY web/package*.json ./
# 如果node里没有预装yarn
RUN command -v yarn >/dev/null 2>&1 || { npm install -g yarn; }
RUN yarn
COPY web .
RUN yarn build


FROM ubuntu:latest

WORKDIR /app

COPY --from=build /app/room .
COPY --from=node /app/dist web/dist
COPY --from=build /app/configs/openai_config.ini configs/openai_config.ini
COPY start.sh start.sh

EXPOSE 8322 8080

RUN chmod +x start.sh
CMD ["./start.sh"]


## 以下是单独部署后端
#FROM golang:1.19 AS build
#
#WORKDIR /app
#ADD server .
#
#RUN go build -o room .


#FROM ubuntu:latest
#
#WORKDIR /app
#
#COPY --from=build /app/room .
#COPY --from=build /app/configs/openai_config.ini configs/openai_config.ini
#
#EXPOSE 8322
#
#CMD ["./room"]