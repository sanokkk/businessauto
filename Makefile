swag-init:
	swag init -g ./cmd/server/main.go  -o ./docs
swag-fmt:
	swag fmt

minio:
	docker run \
       -p 9000:9000 \
       -p 9001:9001 \
       --name minio \
       -v ./content:/data \
       -e "MINIO_ROOT_USER=minio-admin" \
       -e "MINIO_ROOT_PASSWORD=minio-admax" \
       -d \
       quay.io/minio/minio server /data --console-address ":9001"

run:
	go run cmd/server/main.go

start: swag-fmt swag-init minio run

