## Module-to-Module Call

Sometimes, it is necessary that a module calls another module. For example, `example/v2` needs to call `toggle/v1alpha1` in its use case.
For this case to happen, there are two options. The first is `toggle/v1alpha1` to provide SDK. The second is `example/v2` implements the client itself.

We prefer the first approach. From previous example, the target is the `toggle/v1alpha1` module. Hence, the `toggle/v1alpha1` module provides SDK in which it only accepts a few configs to work. This way, the idiom API first is built with strong consistency.

By using this approach, the `toggle/v1alpha1` module can define the behavior of its client (SDK) and protect itself from any upcoming harms (because they are the ones who implement the client). It also avoids [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) by preventing clients to re-implement codes that have the same knowledge and goal. For comparison, Google SDK and client are built this way.

This approach enforces developers to provide SDK for other developers.

### SDK

The SDK should be located in `/pkg` directory inside the module's directory. Read more about pkg directory in [https://github.com/golang-standards/project-layout#pkg](https://github.com/golang-standards/project-layout#pkg).

### Example

- The `toggle/v1alpha1` SDK can be seen in [toggle/v1alpha1/pkg](../modules/toggle/v1alpha1/pkg).
- The `example/v2` that uses `toggle/v1alpha1` SDK can be seen in [example/v2/service](../modules/example/v2/service)