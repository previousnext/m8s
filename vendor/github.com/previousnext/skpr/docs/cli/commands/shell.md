# shell

Opens an interactive SSH session to the specified environment.

```bash
skpr shell [environment]

Arguments:
  [environment] Environment as configured in skipper. Usually one of  prod/staging/dev.

Flags:
  --instance    The name of instance. If not provided one will be chosen for you.
  --container   The name of the container to exec on.
```

## Examples

SSH into `staging` environment.

```bash
$ skpr shell staging

root@staging-438573485:/data
```

SSH into `solr` container on `prod` environment.

```bash
$ skpr shell prod --container "solr"
```
