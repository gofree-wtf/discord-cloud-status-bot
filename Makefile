init:
	go mod download

build:
	go build -o bin/discord-cloud-status-bot cmd/bot.go

run:
	go run cmd/bot.go -config.file ./configs/app.yaml

test:
	go test -v ./pkg/...

docker-build:
	docker build -t discord-cloud-status-bot:latest .

docker-run:
	docker run --rm \
	-e BOT_TOKEN=INSERT_HERE \
	discord-cloud-status-bot:latest
