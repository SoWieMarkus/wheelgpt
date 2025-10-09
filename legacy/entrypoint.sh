#!/bin/sh
set -e

cd /usr/src/app
npx prisma migrate deploy
npx prisma generate --sql
npm start --workspace=@wheelgpt/backend