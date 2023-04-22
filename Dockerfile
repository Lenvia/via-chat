FROM golang:1.19 AS build

WORKDIR /app
ADD . .

RUN go build -o room .

# 安装 Node.js 和 yarn
FROM node:latest AS node
WORKDIR /app
# 拷贝配置文件
COPY web/chat/package*.json ./
RUN npm install -g yarn
RUN yarn
COPY web/chat .
RUN yarn build


FROM ubuntu:latest

WORKDIR /app

COPY --from=build /app/room .
COPY --from=node /app/dist web/chat/dist
COPY --from=build /app/configs/openai_config.ini configs/openai_config.ini

EXPOSE 8322 8080

CMD ["./start.sh"]


## 以下是单独部署后端
#FROM golang:1.19 AS build
#
#WORKDIR /app
#ADD . .
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