# Packaging

Packaging containers is handled using the project `Dockerfile`. Multistage builds are supported.

By default, `skpr package` will use the `Dockerfile` in the project root.

> See [`skpr package`](../cli/commands/package.md)

## Example Dockerfile

This example demonstrates using a multi-stage dockerfile to fetch the dependencies using the `-dev` container image (which has dev tooling like yarn / composer etc..) and then copies the compiled artefact into the `-apache` image which is used for skipper environments.

```dockerfile
FROM previousnext/php:7.2-dev
COPY . /data
RUN make init-package

FROM previousnext/php:7.2-apache
COPY --from=0 /data /data
```
