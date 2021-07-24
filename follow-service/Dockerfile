FROM golang:1.16-alpine3.14
COPY . /photo-viewer/follow-service
WORKDIR /photo-viewer/follow-service

ENV POSTGRE_HOST=34.65.207.86
ENV POSTGRE_PORT=5432
ENV POSTGRE_USER=postgres
ENV POSTGRE_PASSWORD=postgres
ENV POSTGRE_DB_NAME=photo-viewer
ENV FOLLOW_SERVICE_PORT=:10001
ENV FULL_HOSTNAME=http://localhost

RUN go build -o ./out/follow-service .
CMD ["./out/follow-service"]