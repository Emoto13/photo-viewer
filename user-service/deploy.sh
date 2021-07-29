#!/bin/bash
docker build -t user-service .

docker tag user-service gcr.io/photo-viewer-project/user-service
docker push gcr.io/photo-viewer-project/user-service:latest

gcloud config set run/region europe-west6
gcloud run deploy user-service --image gcr.io/photo-viewer-project/user-service

