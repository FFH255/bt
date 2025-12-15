docker:
	open --background -a Docker

debug:
	docker compose -f docker-compose.debug.yaml down --volumes
	docker-compose -f docker-compose.debug.yaml up

local:
	docker compose -f docker-compose.local.yaml down --volumes
	docker-compose -f docker-compose.local.yaml up --build

prod:
	docker compose -f docker-compose.prod.yaml down --volumes
	docker-compose -f docker-compose.prod.yaml up --build

swag:
	swag init -g ./src/internal/shared/proto/swagger/handler.go -o ./src/docs
