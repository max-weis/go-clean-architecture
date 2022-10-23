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

db:
	docker rm -f postgres
	docker volume rm build_pg_data
	docker compose -f ./build/docker-compose.yml up --build postgres

apitest:
	docker compose -f ./build/docker-compose.yml up -d --build
	sleep 5
	newman run apitests/webshop.postman_collection.json
	docker compose -f ./build/docker-compose.yml kill
	docker compose -f ./build/docker-compose.yml rm -fv

lint:
	golangci-lint run -v