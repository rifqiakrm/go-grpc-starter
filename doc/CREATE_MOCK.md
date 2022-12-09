## How to Create Mock

If there are some interfaces that you need some mocks for them, you can use `mockgen`.
We provide a simple command to generate necessary mocks.

```
$ make mockgen
```

That command will create all mocks for all interfaces available in this repository.
The mocks will be available in directory `/test/mock`.
The mocks will structured as its original source directory.

For example, mocks for interfaces in directory `/modules/example/v2/service` will be created in directory `/test/mock/modules/example/v2/service`.