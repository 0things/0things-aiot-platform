ARG REGISTRY=docker.io
FROM --platform=$BUILDPLATFORM ${REGISTRY}/golang:1.27.1-alpine AS builder
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

ARG APP_RELATIVE_PATH=./cmd/server
ARG TARGETOS
ARG TARGETARCH

WORKDIR /data/app
COPY . .

RUN rm -rf /data/app/bin/
RUN export GOPROXY=https://goproxy.cn,direct && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o ./bin/server ${APP_RELATIVE_PATH}
RUN mv config /data/app/bin/

FROM ${REGISTRY}/alpine:3.24
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

ARG APP_ENV
ENV APP_ENV=${APP_ENV}

WORKDIR /data/app
COPY --from=builder /data/app/bin /data/app

ENTRYPOINT [ "./server" ]
