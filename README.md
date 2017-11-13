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

* [GKE](/docs/cluster/gcp/gcp.md)
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

* [Drupal](/docs/projects/drupal/drupal.md)

## Documentation

* [Cache](/docs/cache.md)
* [Solr](/docs/solr.md)
* [Mailhog](/docs/mailhog.md)

## Acknowledgements

Built in partnership with:

* Transport for NSW - https://www.transport.nsw.gov.au

## Development

### Roadmap

Our product roadmap can be found [here](/ROADMAP.md)

### Tools

* **Dependency management** - https://github.com/golang/dep
* **Build** - https://github.com/mitchellh/gox
* **Linting** - https://github.com/golang/lint
* **GitHub Releases** - https://github.com/tcnksm/ghr

### Workflow

(While in the `workspace` directory)

**Installing a new dependency**

```bash
dep ensure -add github.com/foo/bar
```

**Running quality checks**

```bash
make lint test
```

**Building binaries**

```bash
make build
```

**Release**

Release artifacts are pushed to the [github releases page](https://github.com/previousnext/m8s/releases) when tagged
properly. Use [semantic versioning](http://semver.org/) prefixed with `v` for version scheme. Examples:

- `v1.0.0`
- `v1.1.0-beta1`
