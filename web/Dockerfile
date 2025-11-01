FROM node:22-alpine

RUN corepack enable

WORKDIR /app

COPY . .

RUN if [ -f pnpm-lock.yaml ]; then \
      pnpm install --frozen-lockfile; \
    else \
      pnpm install; \
    fi

EXPOSE 5173

CMD ["pnpm", "dev", "--host", "0.0.0.0"]
