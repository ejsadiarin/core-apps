FROM node:22-alpine AS base

FROM base AS pruner
RUN apk add --no-cache libc6-compat
WORKDIR /app
RUN npm install -g turbo
COPY . .
RUN turbo prune --scope=core --docker

FROM base AS installer
RUN apk add --no-cache libc6-compat
WORKDIR /app
RUN npm install -g pnpm
COPY --from=pruner /app/out/json/ .
COPY --from=pruner /app/out/pnpm-lock.yaml ./pnpm-lock.yaml
RUN pnpm install --frozen-lockfile

FROM base AS builder
WORKDIR /app
RUN npm install -g pnpm
COPY --from=pruner /app/out/full/ .
COPY --from=installer /app/node_modules ./node_modules
COPY --from=installer /app/packages ./packages
RUN pnpm turbo run build --filter=core...

FROM base AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/apps/core/public ./apps/core/public
COPY --from=builder --chown=nextjs:nodejs /app/apps/core/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/apps/core/.next/static ./apps/core/.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000
CMD ["node", "apps/core/server.js"]
