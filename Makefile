# Makefile raiz - Pet Services

.PHONY: help up down logs build api infra clean

help:
	@echo "Comandos disponíveis:"
	@echo "  make up         - Sobe toda a stack (infra: API, banco, MinIO)"
	@echo "  make down       - Para e remove os containers da stack"
	@echo "  make logs       - Mostra os logs dos serviços da stack"
	@echo "  make build      - Builda a imagem da API"
	@echo "  make api        - Executa make na pasta da API"
	@echo "  make infra      - Entra na pasta de infraestrutura"
	@echo "  make clean      - Limpa artefatos de build da API"

up:
	cd pet-services-infra && docker compose up -d

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

clean:
	cd pet-services-api && make clean
