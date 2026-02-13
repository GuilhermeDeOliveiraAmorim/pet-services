# Dockerfile para pet-services-api
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download -C pet-services-api
RUN CGO_ENABLED=0 GOOS=linux go build -C pet-services-api -o /bin/api ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/api /bin/api
EXPOSE 8080
CMD ["/bin/api"]
