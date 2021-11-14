FROM node:lts-alpine3.12 as build

WORKDIR /app

COPY package.json /app/package.json
RUN npm install

COPY . /app/

RUN npm run build

FROM nginx:1.21-alpine

COPY --from=build /app/dist/frontend /usr/share/nginx/html
COPY --from=build /app/nginx/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]