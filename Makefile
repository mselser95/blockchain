pre-commit:
	pre-commit install

test:
	go test -cover -v ./...

mockgen:
	go generate ./...