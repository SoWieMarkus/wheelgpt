# --- Builder stage ---
FROM node:24.7.0-alpine AS builder

WORKDIR /usr/src/app

COPY . .

# Install dependencies and build all packages
RUN npm install
RUN npm install -g @angular/cli
ENV DATABASE_URL=file:./wheelgpt.db
RUN npx prisma migrate deploy
RUN npx prisma generate --sql
RUN npm run build

FROM node:24.7.0-alpine

WORKDIR /usr/src/app

# Copy only necessary files from each package
COPY --from=builder /usr/src/app/package*.json ./
COPY --from=builder /usr/src/app/entrypoint.sh ./entrypoint.sh
COPY --from=builder /usr/src/app/prisma ./prisma

# Copy dist folders from each package
COPY --from=builder /usr/src/app/packages/backend/dist ./packages/backend/dist
COPY --from=builder /usr/src/app/packages/backend/package.json ./packages/backend/package.json
COPY --from=builder /usr/src/app/packages/frontend/dist ./packages/frontend/dist
COPY --from=builder /usr/src/app/packages/frontend/package.json ./packages/frontend/package.json
COPY --from=builder /usr/src/app/packages/shared/dist ./packages/shared/dist
COPY --from=builder /usr/src/app/packages/shared/package.json ./packages/shared/package.json

# Install only production dependencies
RUN npm install --omit=dev

# Entrypoint: deploy migrations and start backend
RUN chmod +x /usr/src/app/entrypoint.sh
CMD ["/usr/src/app/entrypoint.sh"]