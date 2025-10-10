TEST_FLAGS := -v -cover -timeout 30s

# test
.PHONY: test
test:
	go test -v -cover -timeout 30s ./... -coverprofile=cover.out.tmp
	cat cover.out.tmp | grep -v "config/*" | grep -v "cmd/*" | grep -v "mock/*" | grep -v "pkg/*" | grep -v "infrastructure/*" | grep -v "app/*" | grep -v "domain/repository/*" > cover.out
	rm cover.out.tmp
	go tool cover -html=cover.out -o cover.html
	open cover.html

# lint
.PHONY: lint
lint:
	golangci-lint run ./...

# generate
.PHONY: generate
generate:
	go generate ./...

# fmt
.PHONY: fmt
fmt:
	go fmt ./...

# docker compose down all
.PHONY: down-all
down-all:
	docker compose down --rmi local --volumes

# docker compose down rmi
.PHONY: down-rmi
down-rmi:
	docker compose down --rmi local

# docker compose down
.PHONY: down
down:
	docker compose down

# docker compose up
.PHONY: up
up:
	docker compose up -d

# modernize
.PHONY: modernize
modernize:
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test ./...
