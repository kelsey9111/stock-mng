APP_NAME = server

dev:
	go run ./cmd/${APP_NAME}/ 

swager:
	swag init -g ./cmd/server/main.go -o ./cmd/swager/docs

