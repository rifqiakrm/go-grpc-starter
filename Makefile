GO_FILES = $(shell go list ./... | grep -v /test/integration/ | grep -v /features/)

.PHONY: format
format:
	bin/format.sh

.PHONY: check.import
check.import:
	bin/check-import.sh

.PHONY: cleanlintcache
cleanlintcache:
	golangci-lint cache clean

.PHONY: lint
lint: cleanlintcache
	golangci-lint run ./...

.PHONY: pretty	
pretty: tidy format lint

.PHONY: cleantestcache
cleantestcache:
	go clean -testcache

.PHONY: test.unit
test.unit: cleantestcache
	go test -v -race $(GO_FILES)

.PHONY: mockgen
mockgen:
	bin/generate-mock.sh

.PHONY: dep-download
dep-download:
	GO111MODULE=on go mod download

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: cover
cover: cleantestcache
	go test -v -race $(GO_FILES) -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out 

.PHONY: coverhtml
coverhtml: cleantestcache
	go test -v -race $(GO_FILES) -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: compile-server
compile-server:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o deploy/grpc-starter cmd/server/main.go

.PHONY: docker-build-server
docker-build-server:
	docker build --no-cache -t grpc-starter-server:latest -f Dockerfile .

.PHONY: tests
tests:
	bin/godog.sh

.PHONY: test
test:
	bin/godog-feature.sh $(feature)

.PHONY: generate.report
generate.report:
	bin/generate-report.sh

.PHONY: all-db-migrate
all-db-migrate:
	bin/migrate.sh $(url)

.PHONY: migration
migration:
	migrate create -ext sql -dir db/migrations/$(module) $(name)

.PHONY: migrate
migrate:
	migrate -path db/migrations/$(module) -database "$(url)?sslmode=disable&search_path=$(module)" -verbose up

.PHONY: rollback
rollback:
	migrate -path db/migrations/$(module) -database "$(url)?sslmode=disable&search_path=$(module)" -verbose down 1

.PHONY: rollback-all
rollback-all:
	migrate -path db/migrations/$(module) -database "$(url)?sslmode=disable&search_path=$(module)" -verbose down -all

.PHONY: force-migrate
force-migrate:
	migrate -path db/migrations/$(module) -database "$(url)?sslmode=disable&search_path=$(module)" -verbose force $(version)

.PHONY: schema
schema:
	migrate create -ext sql -dir db/schemas $(name)

.PHONY: migrate-schema
migrate-schema:
	migrate -path db/schemas -database "$(url)?sslmode=disable" -verbose up

.PHONY: rollback-schema
rollback-schema:
	migrate -path db/schemas -database "$(url)?sslmode=disable" -verbose down 1

.PHONY: force-schema
force-schema:
	migrate -path db/schemas -database "$(url)?sslmode=disable" -verbose force $(version)

.PHONY: rollback-schema-all
rollback-schema-all:
	migrate -path db/schemas -database "$(url)?sslmode=disable" -verbose down -all

.PHONY: validate-migration
validate-migration:
	bin/validate-migration.sh

.PHONY: generate-pb
generate-pb:
	bin/generate-pb.sh
