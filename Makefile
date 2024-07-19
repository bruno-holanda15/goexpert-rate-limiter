up:
	docker compose up -d

down:
	docker compose down

logs:
	docker logs -f rate-limiter -n 20

tests:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html