FROM node:25-alpine AS builder
ARG VITE_LOGTO_ENDPOINT=http://localhost:3001
ARG VITE_LOGTO_APP_ID=
ARG VITE_LOGTO_RESOURCE=http://localhost:8000
RUN npm install -g pnpm@10
WORKDIR /app
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend ./
ENV VITE_DEVICE_SERVICE_URL=/api \
    VITE_AUTH_SERVICE_URL=/api \
    VITE_NOTIFICATION_SERVICE_URL=/api \
    VITE_AI_GATEWAY_URL="" \
    VITE_LOGTO_ENDPOINT=${VITE_LOGTO_ENDPOINT} \
    VITE_LOGTO_APP_ID=${VITE_LOGTO_APP_ID} \
    VITE_LOGTO_RESOURCE=${VITE_LOGTO_RESOURCE}
RUN pnpm build

FROM nginx:stable-alpine
COPY deploy/compose/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html
RUN chmod -R 755 /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
