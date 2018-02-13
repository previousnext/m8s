# config

`skpr config` has several sub commands.

* [config delete](#delete)
* [config get](#get)
* [config list](#list)
* [config set](#set)

## delete

Deletes a config value for the specified environment.

```bash
skpr config delete [environment] [key]

Arguments:
  [environment]   Name of environment.
  [key]           Name of key/value pair to delete.
```

### Examples

Delete `healthz.skip` config from `prod` environment.

```bash
$ skpr config delete prod healthz.skip
```

## get

Get a single config value for the specified environment.

```bash
skpr config get [environment] [key]

Arguments:
  [environment]   Name of environment.
  [key]           Name of key/value pair to retrieve.
```

### Examples

Get `mysql.db.name` config value for `staging` environment.

```bash
$ skpr config get staging mysql.db.name

foo_staging
```

## list

List all of the config key/value pairs for the specified environment.

```bash
skpr config list [environment]

Arguments:
  [environment]   Name of environment.

  Flags:
    --format      Output format - supported formats are "text" (default) and "json".
```

### Examples

Show all config key/value pairs for `prod` environment.

```bash
$ skpr config list prod

KEY                 VALUE
--------            --------
mysql.db.host       svc-project-prod.abc.ap-southeast-2.rds.amazonaws.com
mysql.db.name       project_prod
mysql.db.pass       [secret]
mysql.db.user       project_prod
nr.app.name         Project Name - Prod
smtp.user           AKIAJF4GYPDIQHYXIA5Q
smtp.pass           [secret]
```

Show all config key/value pairs for `prod` environment in json format.

```bash
$ skpr config list prod --format json
{
  "mysql.db.host": "svc-project-prod.abc.ap-southeast-2.rds.amazonaws.com",
  "mysql.db.name": "project_prod",
  "mysql.db.pass": "",
  "mysql.db.user": "project_prod",
  "nr.app.name": "Project Name - Prod",
  "smtp.user": "AKIAJF4GYPDIQHYXIA5Q",
  "smtp.pass": "",
}
```

## set

Set a config value for the specified environment.

```bash
skpr config set [environment] [key] [value]

Arguments:
  [environment]   Name of environment.
  [key]           Name of key/value pair to set.
  [value]         Value of the config.

Flags:
  --secret        Indicates the config value is a secret.
                  This option encrypts the value at rest, and obfuscates the value when displayed with the config list command.
```

### Examples

Set the New Relic app name for the `dev` environment.

```bash
$ skpr config set dev nr.app.name 'Project - Dev'
```

Set the SMTP password for the `prod` environment.

```bash
$ skpr config set --secret prod smtp.pass af94dd731834a174241cecbadd9b1d77bfa2f71a
```
