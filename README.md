# GRPC Golang Starter API

## Description

GRPC Golang Starter API is a GRPC microservice that provides a simple and secure way to access GRPC Starter. Build with monorepo, it is easy to use and easy to maintain.

It comes pre-configured with :

1. Google Cloud Error Reporting (<https://cloud.google.com/go/errorreporting>)
2. Google Cloud Profiler (<https://cloud.google.com/go/profiler>)
3. Google Pub Sub(<https://cloud.google.com/go/pubsub>)
4. JWT GO (<https://github.com/dgrijalva/jwt-go>)
5. Go Mock (<https://github.com/golang/mock>)
6. Redis (<https://github.com/gomodule/redigo>)
7. Google UUID (<https://github.com/google/uuid>)
8. GRPC (<https://google.golang.org/grpc>)
9. GRPC Gateway(<https://github.com/grpc-ecosystem/grpc-gateway/v2>)
10. GRPC OpenTracing (<https://github.com/grpc-ecosystem/grpc-opentracing>)
11. GRPC Middleware (<https://github.com/grpc-ecosystem/go-grpc-middleware>)
12. Protocol Buffers (<https://google.golang.org/protobuf>)
13. Protoc Gen Validate (<https://github.com/envoyproxy/protoc-gen-validate>)
14. Prometheus (<https://github.com/grpc-ecosystem/go-grpc-prometheus>)
15. GORM (<https://gorm.io/gorm>)
16. PGX Postgres (<https://github.com/jackc/pgx/v4>)
17. ENV Decode(<https://github.com/joeshaw/envdecode>)
18. Godotenv (<https://github.com/joho/godotenv>)
19. Mailgun (<https://github.com/mailgun/mailgun-go/v4>)
20. Sendgrid (<https://github.com/sendgrid/sendgrid-go>)
21. Testify (<https://github.com/stretchr/testify>)


## Setup

Use this command to install the blueprint

```bash
go get github.com/rifqiakrm/grpc-starter
```

or manually clone the repo.

## How to Run

- Read [Prerequisites](doc/PREREQUISITES.md).
- Then, read [How to Run](doc/HOW_TO_RUN.md).

## Development Guide

- Read [Prerequisites](doc/PREREQUISITES.md).
- Then, read [Development Guide](doc/DEVELOPMENT_GUIDE.md).

## Test

### Unit Test

```sh
$ make tidy
$ make cover
```

### Integration Test / API Test

Before you start creating gherkin scenarios, you need to install [godog](https://github.com/cucumber/godog) and read the step definition documentation [here](doc/GODOG_DOCUMENTATION.md).

To run integration test, we need to start all dependencies needed. We provide all dependencies via [Docker](https://docs.docker.com/engine/install)
Make sure to install [Docker](https://docs.docker.com/engine/install) before running integration test.

Also, we need to build the docker image for grpc-starter first.

```sh
$ make compile-server
$ make docker-build-server
```

After that, run all images needed using `docker-compose` and run the integration test.

```sh
$ docker run -d -p 8080:8080 -p 8081:8081 --env-file .env grpc-starter-server:latest 
$ make test.integration
```

## Deployment

Read [Deployment](doc/DEPLOYMENT.md).

### Staging

TBD

### Production

TBD

### Postman Collection

### Staging

TBD

### Production

TBD

## FAQs

- Read [FAQs](doc/FAQS.md)
