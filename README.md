Mates (m8s)
===========

[![CircleCI](https://circleci.com/gh/previousnext/m8s.svg?style=svg&circle-token=cd74c538bea3d8ae9d0de9b008fedf35b7f00ad8)](https://circleci.com/gh/previousnext/m8s)

**Maintainer**: Nick Schuch

CLI and API for building temporary environments.

![Diagram](/docs/diagram.png "Diagram")

## How this works

### Caches

We can the following directories by default:

* **Composer** - /root/.composer
* **Yarn** - /usr/local/share/.cache/yarn

## Development

### Principles

* Code lives in the `workspace` directory

### Tools

* **Dependency management** - https://getgb.io
* **Build** - https://github.com/mitchellh/gox
* **Linting** - https://github.com/golang/lint

### Workflow

(While in the `workspace` directory)

**Installing a new dependency**

```bash
gb vendor fetch github.com/foo/bar
```

**Running quality checks**

```bash
make lint test
```

**Building binaries**

```bash
make build
```
