FROM golang:1.16-alpine3.14
COPY . /photo-viewer/auth-service
WORKDIR /photo-viewer/auth-service

RUN go build -o ./out/auth-service .
CMD ["./out/auth-service"]