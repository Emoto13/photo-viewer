#!/bin/bash
docker build -t auth-service .

docker tag auth-service gcr.io/photo-viewer-project/auth-service
docker push gcr.io/photo-viewer-project/auth-service:latest

gcloud config set run/region europe-west6
gcloud run deploy auth-service --image gcr.io/photo-viewer-project/auth-service

