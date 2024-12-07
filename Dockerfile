FROM node:20

WORKDIR /usr/src/app

ENV DATABASE_URL=file:/usr/src/app/wheelgpt.db
COPY package*.json ./

RUN npm install

COPY . .

RUN npx prisma generate
RUN npm run build

ENV PORT=3000

EXPOSE $PORT

CMD ["sh", "-c", "npx prisma migrate deploy && npm start"]