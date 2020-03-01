.PHONY: build

DUMMY?=dummy.json

build: clean
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/smtpless
	chmod +x build/smtpless
	build-lambda-zip -o build/smtpless.zip build/smtpless config.json


deploy: build
	sls deploy --verbose

# run function locally
invoke: build
	sls invoke local --function mail --path ${DUMMY}

clean:
	rm -rf ./build
