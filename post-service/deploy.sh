#!/bin/bash
docker build -t post-service .

docker tag post-service gcr.io/photo-viewer-project/post-service
docker push gcr.io/photo-viewer-project/post-service:latest

gcloud config set run/region europe-west6
gcloud run deploy post-service --image gcr.io/photo-viewer-project/post-service

