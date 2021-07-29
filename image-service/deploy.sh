#!/bin/bash
docker build -t image-service .

docker tag image-service gcr.io/photo-viewer-project/image-service
docker push gcr.io/photo-viewer-project/image-service:latest

gcloud config set run/region europe-west6
gcloud run deploy image-service --image gcr.io/photo-viewer-project/image-service

