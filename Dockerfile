FROM golang:1.23 AS build

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

ADD . /data/build
WORKDIR /data/build

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o mall.backend main.go

RUN mkdir -p /data/wwwRoot/
RUN pwd && ls -l
RUN mv mall.backend /data/wwwRoot/mall.backend
RUN chmod +x /data/wwwRoot/mall.backend
RUN rm -rf /data/build

# golang mini runtime linux alpine
FROM alpine:3.21

RUN mkdir -p /data/wwwRoot/
COPY --from=build /data/wwwRoot/mall.backend /data/wwwRoot/mall.backend

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repository
RUN apk update && apk add tzdata
RUN echo 'Asia/Shanghai' > /etc/timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /data/wwwRoot