sense: 
	go build -o bin/app ./
test: 
	go test ./
clean: 
	rm -f bin/app
run:
	go build -o bin/app .
	docker-compose up -d
	./bin/app
deps:
	dep ensure