FROM node:24.10.0-alpine AS build
WORKDIR /wheelgpt
RUN npm cache clean --force
COPY . .
RUN npm install
RUN npm run build


FROM nginx:alpine AS nginx
COPY --from=build /wheelgpt/dist/frontend/browser /usr/share/nginx/html
COPY /nginx.conf  /etc/nginx/conf.d/default.conf
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]