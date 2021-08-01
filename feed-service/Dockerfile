FROM golang:1.16-alpine3.14
COPY . /photo-viewer/feed-service
WORKDIR /photo-viewer/feed-service

RUN go build -o ./out/feed-service .
CMD ["./out/feed-service"]