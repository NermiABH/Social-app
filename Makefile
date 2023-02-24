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
test:
	go test -v ./...