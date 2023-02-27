build:
	docker-compose build social-app
run:
	docker-compose up social-app
restart:
	docker-compose restart
rebuild:
	docker-compose up -d --no-deps --build social-app
stop:
	docker-compose stop
migrate:
	migrate -path migrations -database 'postgres://postgres:pusinu48@:5436/postgres?sslmode=disable' up
test:
	go test -v ./...
start:
	go build -o social-app ./cmd/main.go && ./social-app