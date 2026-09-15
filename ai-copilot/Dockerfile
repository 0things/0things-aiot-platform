FROM node:25-alpine AS builder
RUN npm install -g pnpm@10
WORKDIR /app
COPY ai-copilot/package.json ai-copilot/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY ai-copilot/tsconfig.json ./
COPY ai-copilot/src ./src
RUN pnpm build && pnpm prune --prod

FROM node:25-alpine
WORKDIR /app
COPY --from=builder --chown=node:node /app/package.json ./
COPY --from=builder --chown=node:node /app/node_modules ./node_modules
COPY --from=builder --chown=node:node /app/dist ./dist
USER node
EXPOSE 8005
CMD ["node", "dist/index.js"]
