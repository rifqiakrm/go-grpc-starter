## Development Guide

### Adding a Module

Before developing your changes, you are encouraged to read [How to Add Module](HOW_TO_ADD_MODULE.md).
Please, take your time to read the document before continuing.

### Code Structure

Take your time to read [Code Map](CODE_MAP.md) to familiarize yourself with the code structure.
We provide two examples of the real implementation. You can see it in [Toggle Module](/modules/toggle/v1alpha1) and [Example Module](/modules/example/v2).

### Making Changes

Read [Making Changes](MAKING_CHANGES.md).

### Database Migration

Read [Database Migration](DATABASE_MIGRATION.md).

### Module-to-Module Call

Sometimes, it is necessary that a module calls another module. For example, `example/v2` needs to call `toggle/v1alpha1` in its use case. To know how to implement module-to-module call, read [Module-to-Module Call](MODULE_TO_MODULE_CALL.md).