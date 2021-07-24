FROM golang:1.16-alpine3.14
COPY . /photo-viewer/auth-service
WORKDIR /photo-viewer/auth-service

ENV POSTGRE_HOST=34.65.207.86
ENV POSTGRE_PORT=5432
ENV POSTGRE_USER=postgres
ENV POSTGRE_PASSWORD=postgres
ENV POSTGRE_DB_NAME=photo-viewer
ENV AUTH_SERVICE_PORT=:10000

RUN go build -o ./out/auth-service .
CMD ["./out/auth-service"]