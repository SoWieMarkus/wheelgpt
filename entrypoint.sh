#!/bin/sh
set -e

cd /usr/src/app
npx prisma migrate deploy
npm start --workspace=@wheelgpt/backend