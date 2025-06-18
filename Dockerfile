FROM node:24

WORKDIR /usr/src/app

COPY . .

# Install dependencies
RUN npm install 
RUN npm install -g @angular/cli 

# Generate Prisma Client
WORKDIR /usr/src/app/packages/backend
ENV PORT=3000
ENV DATABASE_URL=file:/usr/src/app/wheelgpt.db
RUN npx prisma generate

WORKDIR /usr/src/app
RUN npm run build

EXPOSE $PORT

# Entrypoint: deploy migrations and start backend
RUN chmod +x /usr/src/app/entrypoint.sh
CMD ["/usr/src/app/entrypoint.sh"]