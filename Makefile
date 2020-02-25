.PHONY: build

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/smtpless main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose