## Code Map

Code Map explains how the codes are structured in this project. This document will list down all folders or packages and their purpose.

---

### `bin`

This folder contains any executable binary to support the project.
For example, there is `generate-mock.sh`. It is a shell script file used to generate mocks for all interfaces available in this project.

---

### `cmd`

This folder contains the `main.go`.
The use case may be run and served in multi forms, such as API, cron, or fullstack web.
To cater that case, `cmd` folder can contains subfolders with each folder named accordingly to the form and contain only main package.
e.g: `cmd/api/main.go`, `cmd/cron/main.go`, and `cmd/web/main.go`

For this project, we prefer to use `cmd/server/main.go` as our use cases are only in the form of gRPC server.

---

### `common`

This folder contains common or shared functionalities that are used by (almost) all modules.

---

### `common/config`

This folder contains configuration for the project.

---

### `common/healthcheck`

This folder contains functionality to perform health check. Actually, what's inside this folder is similar to a module since it exposes a gRPC service.
But, this folder doesn't contain any vertical business focus. It is only used by [Kubernetes to check the container's health](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/). Therefore, we put this in `common` folder.

---

### `common/postgres`

This folder contains connection to PostgreSQL.

---

### `common/redis`

This folder contains connection to Redis.

---

### `db/migrations`

This folder contains all database migration files. It has many subdirectories. Each subdirectory represents a single module.

---

### `db/migrations/<module-name>`

This folder contains all module's migration files. Each migration has exactly two files: UP and DOWN. 

---

### `db/schemas`

This folder contains all database schema migration files.

---

### `deployment`

This folder contains Kubernetes manifest related to deployment.

---

### `doc`

This folder contains all documents related to the project.

---

### `features`

This folder contains all [Cucumber](https://cucumber.io/docs/guides/) definitions for the purpose of integration / API test.
We use [Gherkin syntax](https://cucumber.io/docs/gherkin/) to define the features.

---

### `features/modules/<module-name>/<module-version>`

This folder contains all features of a specific module. Inside this folder, there should be only two kind of files.
The first if `*.feature` files which define the feature and scenario. The second is `<module-name>_test.go` which define the Go implementation of [Cucumber](https://cucumber.io/docs/guides/) using [Godog](https://github.com/cucumber/godog).

---

### `modules`

This folder contains all vertical businesses logic.

---

### `modules/<module-name>/<module-version>/entity`

This folder contains the domain of the module.
Mostly, this folder contains only structs, constants, global variables, enumerations, or functions with simple logic related to the core domain of the module (not a business logic).
Since we use Protocol Buffer, entity has a close (tightly coupled) relationship with any struct generated from `.proto` files.

---

### `modules/<module-name>/<module-version>/service`

This folder contains the main business logic of the module. Almost all interfaces and all the business logic flows are defined here.
If someone wants to know the flow of the module, they better start to open this folder.

---

### `modules/<module-name>/<module-version>/pkg`

This folder contains SDK for the module. It basically a gRPC client. This folder exists to be imported and used by another module for the purpose of module-to-module call. The owners of a module MUST provides this SDK. It is their responsibility to implement and provide the SDK.

For example, module `example/v2` wants to call `toggle/v1alpha1`. Then, `example/v2` just needs to import this folder. The actual implementation can be seen in [example/v2](../modules/example/v2/service).

---

### `modules/<module-name>/<module-version>/internal`

All APIs/codes in the internal folder (and all if its subfolders) are designed to [not be able to be imported](https://golang.org/doc/go1.4#internalpackages).
This folder contains all detail implementation specified in the `modules/<module-name>/<module-version>/service` folder.

---

### `modules/<module-name>/<module-version>/internal/builder`

This folder contains the [builder design pattern](https://sourcemaking.com/design_patterns/builder).
It composes all codes needed to build a full usecase.

---

### `modules/<module-name>/<module-version>/internal/grpc/handler`

This folder contains the HTTP/2 gRPC handlers.
Codes in this folder implement gRPC server interface.

---

### `modules/<module-name>/<module-version>/internal/repository`

This folder contains codes that connect to the repository, such as database.
Repository is not limited to databases. Anything that can be a repository can be put here.

---

### `modules/<module-name>/<module-version>/<module-name>.go`

This is special go file inside a module directory. This file exists only as a proxy to connect module to the `cmd/server/main.go`.
You can think this file as a imaginary `main.go` for the module. There are two methods that must be exist in this file. These are the mandatory functions:

```
func InitGrpc() {}
```
This method is used to initialize and register gRPC service. This method can accept any parameters.
The example can be seen in [toggle.go](../doc/modules/toggle/v1alpha1/toggle.go).

---

```
func InitRest() {}
```
This method is used to initialize and register REST service. This method can accept any parameters.
The example can be seen in [toggle.go](../doc/modules/toggle/v1alpha1/toggle.go).

---

Those two functions must be called in `cmd/server/main.go`.

`InitGrpc()` must be called inside `registerGrpcHandlers` function. 
`InitRest()` must be called inside `registerRestHandlers` function. 

The example can be seen in [cmd/server/main.go](../cmd/server/main.go).

---

### `server`

This folder contains all codes needed to define a gRPC and its REST Gateway server.

---

### `test`

This folder contains test related stuffs.
For the case of unit test, the unit test files are put in the same directory as the files that they test. It is one of the Go best practice, so we follow.

---

### `test/fixture`

This folder contains a well defined support for test.

---

### `test/mock`

This folder contains mock for testing.

---