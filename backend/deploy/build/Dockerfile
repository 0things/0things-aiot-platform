ARG REGISTRY=docker.io
FROM ${REGISTRY}/golang:1.25-alpine AS builder

ARG APP_RELATIVE_PATH=./cmd/server

WORKDIR /data/app
COPY . .

WORKDIR /data/app/backend
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/server ${APP_RELATIVE_PATH}

FROM ${REGISTRY}/alpine:3.20

RUN apk --no-cache add tzdata ca-certificates && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

ARG APP_ENV
ENV APP_ENV=${APP_ENV}

WORKDIR /data/app
COPY --from=builder /bin/server /data/app/server
COPY backend/config /data/app/config

EXPOSE 8000
ENTRYPOINT [ "./server" ]
