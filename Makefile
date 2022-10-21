run:
	go run ./cmd/

compile:
	go build -o webshop ./cmd/

test:
	go test -v ./...

short:
	go test -v -short ./...

docker:
	docker compose -f ./build/docker-compose.yml up --build

lint:
	golangci-lint run -v