run:
	go run ./cmd/

docker:
	docker compose -f ./build/docker-compose.yml up --build

lint:
	golangci-lint run -v