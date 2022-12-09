## How to Add Module

### What is Module?

Module is a single business vertical focus.

### Location

Modules are located in directory `/modules`.

### How to Add a Module?

To add a new module, create a new directory inside `/modules` directory. It is encouraged that you also separate the module by its version by separating their directory.

For example, these commands will create two modules and each module has two versions.

```
$ mkdir -p modules/example/v1alpha1
$ mkdir -p modules/example/v2
$ mkdir -p modules/toggle/v1beta1
$ mkdir -p modules/toggle/v1
```

After you create module's directory, put your code inside that directory.
Follow [Code Map](CODE_MAP.md) to structurize your code.