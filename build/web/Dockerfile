FROM node:16-alpine3.16 as build
WORKDIR /usr/app
COPY /web /usr/app/
RUN npm install && npm run build

FROM nginx:1.21-alpine
COPY --from=build /usr/app/build /usr/app/
COPY --from=build /usr/app/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
