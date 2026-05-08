FROM golang:1.26-alpine AS build

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

ADD . /data/build
WORKDIR /data/build

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o mall.backend main.go

RUN mkdir -p /data/wwwRoot/web/vendor/go-captcha-jslib
RUN cp mall.backend /data/wwwRoot/mall.backend
RUN cp -r web/*.html /data/wwwRoot/web/
RUN cp -r web/vendor /data/wwwRoot/web/

# golang mini runtime linux alpine
FROM alpine:3.21

RUN mkdir -p /data/wwwRoot/
COPY --from=build /data/wwwRoot/ /data/wwwRoot/

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add tzdata
RUN echo 'Asia/Shanghai' > /etc/timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /data/wwwRoot