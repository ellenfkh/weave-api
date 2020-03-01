SHA := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	docker build -t weave-api:$(SHA) .

.PHONY: unit
unit: 
	go test

action:
	echo argument is $(argument)

.PHONY: run
run:
	docker run \
		-v $(baseDir):/files \
		-p 8080:8080 \
		weave-api:$(SHA)  \
		--baseDir /files --port 8080