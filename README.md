# Extreme Monikers Game

To run on google cloud:
```
docker build --platform linux/amd64 -t gcr.io/original-storm-432806-p1/monikers-server .

docker push gcr.io/original-storm-432806-p1/monikers-server 

gcloud run deploy monikers-server  --image gcr.io/original-storm-432806-p1/monikers-server  --platform managed  --region us-west1  --port 8080  --allow-unauthenticated
```