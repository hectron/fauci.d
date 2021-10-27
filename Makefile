.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/slackbot_backend functions/slack/backend/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/slackbot_handler functions/slack/handler/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

test:
	go test ./...
