run:
	cd ./cmd && go run .

up:
	docker compose up -d

dev:
	docker compose -f docker-compose.dev.yml up -d 

upb:
	docker compose up --build -d

down:
	docker compose down