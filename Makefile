swag:
	@swag init -d cmd --parseDependency --parseDepth 3 -o docs/

run:
	@go run cmd/main.go

generate:
	@go generate -v ./...