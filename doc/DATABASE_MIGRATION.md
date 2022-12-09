## Database Migration

Database migration is used to migrate database definition to the actual database.
The migration files are located in directory [db/migrations](../db/migrations).
Each module has its own directory inside [db/migrations](../db/migrations).
Each module must define its database migration inside its own directory. See [db/migrations/assessment](../db/migrations/assessment) for example. 

We use different schema for different module. The schemas are defined in [db/schemas](../db/schemas).

### General Rule

#### Tool

To migrate database, you will use [golang-migrate](https://github.com/golang-migrate/migrate). Follow [golang-migrate installation](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md).

#### SQL File

There will always be two files for each migration. `<seq>_<name>.up.sql` and `<seq>_<name>.down.sql`. The former is used to ENFORCE the changes, the later is used to ROLLBACK the changes.

Both files must be executed successfully. So, you **MUST RUN** `<seq>_<name>.up.sql` and you **MUST ALSO RUN** `<seq>_<name>.down.sql`.
After both files are successfully executed, you then run the `<seq>_<name>.up.sql` once again to make your changes real.

    Q: Why we have to run it three times?
    A: To make sure that our UP and DOWN sql files are executed successfully.

### Schema

Before creating table for your module's use case, you **MUST** define its schema. The schema's name must be the same as module's name. See [db/schemas/assessment](../db/schemas/000001_assessment.up.sql) for example.

Always remember to create your schema first before creating table migrations.

To create schema, run this command:

```
$ make schema name=<module-name>
```

e.g:

```
$ make schema name=user
```

Then, three commands: UP, DOWN, and UP migration. These three commands must all success.

```
$ make migrate-schema url=<database-url>
$ make rollback-schema url=<database-url>
$ make migrate-schema url=<database-url>
```

e.g:

```
$ make migrate-schema url="postgres://user:password@host:port/dbname"
$ make rollback-schema url="postgres://user:password@host:port/dbname"
$ make migrate-schema url="postgres://user:password@host:port/dbname"
```

If one command fails, the database will be dirty and can be cleaned by running this command:

```
make force-schema url=<database-url> version=<latest-clean-version>
```

e.g:

```
make force-schema url="postgres://user:password@host:port/dbname" version=4
```

`latest-clean-version` is the latest migration that was success. Usually, it is the number before the number of your schema migration. For example, if your migration is `000008_module.up.sql`, then the latest clean version is 7.

### Migration

After creating schema, it's time to define your table migration.

To create migration, run this command:

```
$ make migration name=<migration-name> module=<module-name>
```

e.g:

```
$ make migration name=create_table_user module=user
```

Then, three commands: UP, DOWN, and UP migration. These three commands must all success.

```
$ make migrate url=<database-url> module=<module-name>
$ make rollback url=<database-url> module=<module-name>
$ make migrate url=<database-url> module=<module-name>
```

e.g:

```
$ make migrate url="postgres://user:password@host:port/dbname" module=user
$ make rollback url="postgres://user:password@host:port/dbname" module=user
$ make migrate url="postgres://user:password@host:port/dbname" module=user
```

If one command fails, the database will be dirty and can be cleaned by running this command:

```
make force-migrate url=<database-url> module=<module-name> version=<latest-clean-version>
```

e.g:

```
make force-migrate url="postgres://user:password@host:port/dbname" module=user version=4
```

`latest-clean-version` is the latest migration that was success. Usually, it is the number before the number of your schema migration. For example, if your migration is `000008_create_table_user.up.sql`, then the latest clean version is 7. 