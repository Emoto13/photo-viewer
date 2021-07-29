FROM golang:1.16-alpine3.14
COPY . /photo-viewer/follow-service
WORKDIR /photo-viewer/follow-service

RUN go build -o ./out/follow-service .
CMD ["./out/follow-service"]