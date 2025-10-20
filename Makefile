.PHONY: coverage

up:
	docker-compose up --build

down:
	docker-compose down


test:
	go test ./... -v

coverage:
	mkdir -p coverage
	go test ./... -coverprofile=coverage/coverage.out
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

clean:
	docker-compose down --volumes --rmi all

GO_IMAGE ?= golang:1.25-alpine

.PHONY: test-docker coverage-docker

test-docker:
	docker run --rm \
		-v $(shell pwd):/app \
		-v $(HOME)/.cache/go-build:/root/.cache/go-build \
		-v $(HOME)/go/pkg/mod:/go/pkg/mod \
		-w /app \
		$(GO_IMAGE) \
		go test ./... -v

coverage-docker:
	docker run --rm \
		-v $(shell pwd):/app \
		-v $(HOME)/.cache/go-build:/root/.cache/go-build \
		-v $(HOME)/go/pkg/mod:/go/pkg/mod \
		-w /app \
		$(GO_IMAGE) \
		sh -c "mkdir -p coverage && go test ./... -coverprofile=coverage/coverage.out && go tool cover -html=coverage/coverage.out -o coverage/coverage.html"

		