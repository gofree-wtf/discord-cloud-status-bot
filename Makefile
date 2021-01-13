init:
	go mod download

build:
	go build -o bin/discord-cloud-status-bot cmd/bot.go

run:
	go run cmd/bot.go

test:
	go test ./pkg/...

docker-build:
	docker build -t discord-cloud-status-bot:latest .

docker-run:
	docker run --rm discord-cloud-status-bot:latest
