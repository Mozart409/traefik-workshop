FROM node:22-alpine AS base
RUN apk add --no-cache libc6-compat
ENV PNPM_HOME=/pnpm
ENV PATH=$PNPM_HOME:$PATH
# Pin the pnpm version package.json declares - corepack cannot read
# packageManager from package.json in the deps stage, which is only given the
# lockfile, so without this every build picks up a different pnpm.
RUN corepack enable && corepack prepare pnpm@10.8.1 --activate
WORKDIR /app

FROM base AS deps
COPY pnpm-lock.yaml ./
RUN pnpm fetch

FROM base AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY . .
# `pnpm fetch` in the deps stage only fills the virtual store (node_modules/.pnpm).
# This links the packages into node_modules proper, offline, from what was fetched.
RUN pnpm install --offline --frozen-lockfile
RUN pnpm run build

FROM base AS runner
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
USER nextjs
EXPOSE 3000
ENV PORT=3000
CMD ["node", "server.js"]
