up:
	docker compose up -d

down:
	docker compose down

logs:
	docker logs -f rate-limiter -n 20