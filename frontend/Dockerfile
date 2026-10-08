ARG REGISTRY=docker.io
FROM ${REGISTRY}/node:25-alpine AS builder
RUN npm install -g pnpm@10
WORKDIR /app
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY index.html vite.config.ts tsconfig*.json .env.production ./
COPY public/ public/
COPY src/ src/
RUN pnpm build

FROM ${REGISTRY}/nginx:stable-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist/ /usr/share/nginx/html/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
