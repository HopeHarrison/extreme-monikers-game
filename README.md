# Extreme Monikers Game

To run on google cloud:
```
docker build --platform linux/amd64 -t gcr.io/original-storm-432806-p1/monikers-server .

docker push gcr.io/original-storm-432806-p1/monikers-server 

gcloud run deploy monikers-service \
  --source . \
  --platform managed \
  --region us-central1 \
  --set-env-vars PROJECT_ID=original-storm-432806-p1,REGION=us-central1,INSTANCE_NAME=original-storm-432806-p1-dev-db \
  --add-cloudsql-instances original-storm-432806-p1:us-central1:original-storm-432806-p1-dev-db \
  --allow-unauthenticated
```