![M8s](/logo/small.png "Logo")
=========================

[![CircleCI](https://circleci.com/gh/previousnext/m8s.svg?style=svg&circle-token=cd74c538bea3d8ae9d0de9b008fedf35b7f00ad8)](https://circleci.com/gh/previousnext/m8s)

**Maintainer**: Nick Schuch

## Overview

M8s is a CLI and API for building temporary environments in Kubernetes.

Often in your CI/CD workflows, you want a real environment to run automated or manual tests. For example, you might want to preview changes you are making in a branch or Pull Request. m8s provides a simple tool for acheiving this. It takes a docker compose file, and translates that into a pod definition Kubernetes understands, and deploys it on Kubernetes. The pod is _ephemeral_, meaning it's not meant to stick around for long, and any data will be deleted once the pod is removed.

![Diagram](/docs/diagram.png "Diagram")

## Getting Started

To get started you will need 1 of each of the following:

**Cluster**

Kubernetes and the M8s API server running.

* GKE ([Easy](/docs/cluster/gcp/easy.md)|[Comprehensive](/docs/cluster/gcp/comprehensive.md))
* Openshift - Coming soon...
* Kops - Coming soon...

**Pipelines**

A continuous integration service / setup which will send information to the M8s API.

* [CircleCI](/docs/pipeline/circleci/circleci.md)
* Bitbucket Pipelines - Coming soon...
* Jenkins - Coming soon...
* TravisCI - Coming soon...

**Projects**

Example implementations for applications.

* [Drupal](/docs/project/drupal/drupal.md)

## Documentation

* [User Interface](/docs/ui/ui.md)
* [Cache](/docs/cache.md)
* [Solr](/docs/solr.md)
* [Mailhog](/docs/mailhog.md)

## Acknowledgements

Built in partnership with:

* Transport for NSW - https://www.transport.nsw.gov.au

## Developing m8s

### Roadmap

Our product roadmap can be found [here](/issues)

### Getting Started

If you wish to work on m8s or any of its built-in systems, you will first need Go installed on your machine.

#### Manual Setup

For local development of m8s, first make sure Go is properly installed and that a GOPATH has been set. You will also need to add $GOPATH/bin to your $PATH.

Next, using Git, clone this repository into $GOPATH/src/github.com/previousnext/m8s. All the necessary dependencies are either vendored or automatically installed, so you just need to type `make test`. This will run the tests and compile the binary. If this exits with exit status 0, then everything is working!

```bash
$ cd "$GOPATH/src/github.com/previousnext/m8s"
$ make test
```

To compile a development version of m8s, run `make build`. This will build everything using gox and put binaries in the bin and $GOPATH/bin folders:

```bash
$ make build
...

# Linux:
$ bin/m8s_linux_amd64 --help

# OSX:
$ bin/m8s_darwin_amd64 --help
```

#### Easy Setup

Alternatively, you can use the [Docker Compose](docker-compose.yml) stack in the root of this repo to stand up a container with the appropriate dev tooling already set up for you.

Using Git, clone this repo on your local machine. Run the test suite to ensure the tooling works.

```bash
$ docker-compose run --rm dev make test
```

To compile a development version of m8s, run `make build`. This will build everything using gox and put binaries in the bin and $GOPATH/bin folders:

```bash
$ docker-compose run --rm dev make build

...

$ docker-compose run --rm dev bin/m8s_linux_amd64 --help
```

### Dependencies

m8s stores its dependencies under `vendor/`, which [Go 1.6+ will automatically recognize and load](https://golang.org/cmd/go/#hdr-Vendor_Directories). We use [`dep`](https://github.com/golang/dep) to manage the vendored dependencies.

If you're developing m8s, there are a few tasks you might need to perform.

For details, see:

* [Adding a dependency](#adding-a-dependency)
* [Updating a dependency](#updating-a-dependency)

### Tooling

* **Dependency management** - https://github.com/golang/dep
* **Build** - https://github.com/mitchellh/gox
* **Linting** - https://github.com/golang/lint
* **GitHub Releases** - https://github.com/tcnksm/ghr

### Common Tasks

#### Adding a dependency

If you're adding a dependency, you'll need to vendor it in the same Pull Request as the code that depends on it. You should do this in a separate commit from your code, as makes PR review easier and Git history simpler to read in the future.

To add a dependency:

Assuming your work is on a branch called `my-feature-branch`, the steps look like this:

1. Vendor the new dependency.

    ```bash
    dep ensure -add github.com/foo/bar
    ```

2. Review the changes in git and commit them.

#### Updating a dependency

To update a dependency:

1. Update the dependency.

    ```bash
    dep ensure -update github.com/foo/bar
    ```

2. Review the changes in git and commit them.

#### Running quality checks

```bash
make lint test
```

#### Building binaries

```bash
make build
```

#### Release

Release artifacts are pushed to the [github releases page](https://github.com/previousnext/m8s/releases) when tagged
properly. Use [semantic versioning](http://semver.org/) prefixed with `v` for version scheme. Examples:

- `v1.0.0`
- `v1.1.0-beta1`
