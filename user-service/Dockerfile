FROM golang:1.16-alpine3.14
COPY . /photo-viewer/user-service
WORKDIR /photo-viewer/user-service

RUN go build -o ./out/user-service .
CMD ["./out/user-service"]