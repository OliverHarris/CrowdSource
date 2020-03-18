CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o appbuilt
docker build -t project:5000/gameservice:latest .
docker push project:5000/gameservice:latest