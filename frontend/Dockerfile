ARG REGISTRY=docker.io
FROM ${REGISTRY}/node:25-alpine AS builder
RUN npm install -g pnpm@10
WORKDIR /app
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
ARG VITE_LOGTO_ENDPOINT=https://auth.0thing.com
ARG VITE_LOGTO_RESOURCE=https://gateway.0thing.com/api
ARG VITE_LOGTO_APP_ID
RUN test -n "$VITE_LOGTO_APP_ID" && VITE_DEVICE_SERVICE_URL=/api \
    VITE_AI_GATEWAY_URL=/copilot \
    pnpm build

FROM ${REGISTRY}/nginx:stable-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html
RUN chmod -R 755 /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
