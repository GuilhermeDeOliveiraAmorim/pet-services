# Dockerfile para pet-services-api
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/api /bin/api
EXPOSE 8080
CMD ["/bin/api"]
