CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o appbuilt
docker build -t 10.8.0.1:5000/api:latest .
docker push 10.8.0.1:5000/api:latest