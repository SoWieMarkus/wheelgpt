#!/bin/sh
set -e

npx prisma migrate deploy
npm start