FROM golang:1.16-alpine3.14
COPY . /photo-viewer/image-service
WORKDIR /photo-viewer/image-service

RUN go build -o ./out/image-service .
CMD ["./out/image-service"]