FROM golang:1.22-alpine AS build_base

COPY ./ /tmp/build_dir

WORKDIR /tmp/build_dir

COPY go.mod go.sum /
RUN go mod download

COPY . .

RUN go build -o businessauto_backend ./cmd/server/main.go

FROM alpine:3.9

COPY --from=build_base /tmp/build_dir/businessauto_backend /app/businessauto_backend
COPY --from=build_base /tmp/build_dir/config.yml /app/config.yml
COPY --from=build_base /tmp/build_dir/migrations /app/migrations

WORKDIR /app

EXPOSE 8080

CMD ["./businessauto_backend", "main"]