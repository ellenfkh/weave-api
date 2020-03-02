SHA := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	docker build -t weave-api:$(SHA) .

.PHONY: unit
unit: 
	go test

.PHONY: test
test: build
	docker run --entrypoint go weave-api:$(SHA) test

.PHONY: run
run: build
	docker run \
		--mount type=bind,source=$(baseDir),target=/files \
		-p 8080:8080 \
		weave-api:$(SHA)  \
		--baseDir /files --port 8080
