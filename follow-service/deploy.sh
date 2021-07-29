#!/bin/bash
docker build -t follow-service .

docker tag follow-service gcr.io/photo-viewer-project/follow-service
docker push gcr.io/photo-viewer-project/follow-service:latest

gcloud config set run/region europe-west6
gcloud run deploy follow-service --image gcr.io/photo-viewer-project/follow-service

