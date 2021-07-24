FROM golang:1.16-alpine3.14
COPY . /photo-viewer/post-service
WORKDIR /photo-viewer/post-service

ENV POSTGRE_HOST=34.65.207.86
ENV POSTGRE_PORT=5432
ENV POSTGRE_USER=postgres
ENV POSTGRE_PASSWORD=postgres
ENV POSTGRE_DB_NAME=photo-viewer
ENV POST_SERVICE_PORT=:10005
ENV REDIS_PORT=:6379
ENV FULL_HOSTNAME=http://localhost

RUN go build -o ./out/post-service .
CMD ["./out/post-service"]