CGO_ENABLED=0 GOOS=linux GOARCH=amd64
go build -o appbuilt
docker build -t api:latest .
