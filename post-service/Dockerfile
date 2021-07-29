FROM golang:1.16-alpine3.14
COPY . /photo-viewer/post-service
WORKDIR /photo-viewer/post-service

RUN go build -o ./out/post-service .
CMD ["./out/post-service"]