## How to Run

There are two ways to run the application. The first is to setup all dependencies manually.
The second is to use docker images.

### Manual
- Install Protoc Compiler

    Linux 
    ```bash
    $ apt-get update && apt-get install -y protobuf-compiler
    ```
    Mac OSX
    ```bash
    $ brew install protobuf
    ```

- Install golang-migrate
  ```bash
  $ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

- Install all protocol buffers dependencies
  ```bash
  $ export GO111MODULE=auto && go install github.com/envoyproxy/protoc-gen-validate@latest && \
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
  go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest
  ```

- Generate protocol buffers
  ```bash 
  $ make generate-pb
  ```

- Create `.env` file

  You can copy the `env.sample` and change its values.

    ```bash
    $ cp env.sample .env
    ```

- Run the migration

    ```bash
    $ make all-db-migrate url="postgres://YOUR_POSTGRES_USER:YOUR_POSTGRES_PASSWORD@YOUR_POSTGRES_HOST:YOUR_POSTGRES_PORT/YOUR_DATABASE_NAME"
    ```

- Run the initial user seeder

    ```bash
    $ make seed
    ```

- Run the application

    ```bash
    $ go run cmd/server/main.go
    ```

### Docker

- Install [Docker](https://docs.docker.com/engine/install).

- Create `.env` file

  You can copy the `env.sample` and change its values.

    ```bash
    $ cp env.sample .env
    ```

- Download the dependencies

    ```bash
    $ make tidy
    ```

- Compile the backend binary

    ```bash
    $ make compile-server
    ```

- Build backend-server image

    ```bash
    $ make docker-build-server
    ```

- Run docker

    ```bash
    docker run -d -p 8080:8080 -p 8081:8081 --env-file .env grpc-starter-server:latest
    ```