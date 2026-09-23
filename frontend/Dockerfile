FROM node:25-alpine AS builder
RUN npm install -g pnpm@10
WORKDIR /app
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend ./
ENV VITE_DEVICE_SERVICE_URL=/api \
    VITE_AUTH_SERVICE_URL=/api \
    VITE_NOTIFICATION_SERVICE_URL=/api \
    VITE_AI_GATEWAY_URL=""
RUN pnpm build

FROM nginx:stable-alpine
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html
RUN chmod -R 755 /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
