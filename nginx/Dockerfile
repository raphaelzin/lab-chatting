FROM arm32v7/nginx:1.23-alpine

RUN apk update && apk add --no-cache git

COPY ./nginx.conf /etc/nginx/nginx.conf

EXPOSE 4000