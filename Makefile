.PHONY: run-containers
run-containers:
	docker run --rm -d -p 9001:80 --name server1 kennethreitz/httpbin
	docker run --rm -d -p 9002:80 --name server2 kennethreitz/httpbin
	docker run --rm -d -p 9003:80 --name server3 kennethreitz/httpbin
	docker run --rm -d -p 9004:80 --name server4 kennethreitz/httpbin

.PHONY: stop
stop:
	docker stop server1
	docker stop server2
	docker stop server3
	docker stop server4

.PHONY: help
help:
	@echo "Available commands:"
	@grep -E '^\.[A-Z]+:' $(MAKEFILE_LIST) | sed 's/^\.//'

.PHONY: run-proxy-server
run-proxy-server:
	go run cmd/main.go

.PHONY: health_check
health_check:
	go run cmd/main.go & sleep 2
	curl -I http://localhost:8080/ping