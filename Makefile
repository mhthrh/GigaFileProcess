test:
	go test -v -cover ./Test/...

compose_up:
	docker-compose up -d

compose_stop:
	docker-compose stop

compose_down:
	docker-compose down