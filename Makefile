# Makefile raiz - Pet Services

.PHONY: up down logs build api infra

up:
	cd pet-services-infra && docker compose up --build

down:
	cd pet-services-infra && docker compose down

logs:
	cd pet-services-infra && docker compose logs -f

build:
	cd pet-services-api && docker build -t pet-services-api:latest .

api:
	cd pet-services-api && make

infra:
	cd pet-services-infra
