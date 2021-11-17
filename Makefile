install:
	go install .

build:
	go build -o golang-test-task

test:
	go test ./... -v

lint-fix:
	go fmt ./...
