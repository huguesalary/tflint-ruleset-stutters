default: build

test:
	go test ./...

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-stutters ~/.tflint.d/plugins
