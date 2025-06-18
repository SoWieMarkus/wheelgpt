#!/bin/sh
set -e

cd /usr/src/app/packages/backend
npx prisma migrate deploy
npm start