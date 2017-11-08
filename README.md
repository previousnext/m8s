![M8s](/logo/small.png "Logo")
=========================

[![CircleCI](https://circleci.com/gh/previousnext/m8s.svg?style=svg&circle-token=cd74c538bea3d8ae9d0de9b008fedf35b7f00ad8)](https://circleci.com/gh/previousnext/m8s)

**Maintainer**: Nick Schuch

## Overview

m8s is a CLI and API for building temporary environments in Kubernetes.

Often in your CI/CD workflows, you want a real environment to run automated or manual tests. For example, you might want to preview changes you are making in a branch or Pull Request. m8s provides a simple tool for acheiving this. It takes a docker compose file, and translates that into a pod definition Kubernetes understands, and deploys it on Kubernetes. The pod is _ephemeral_, meaning it's not meant to stick around for long, and any data will be deleted once the pod is removed.

![Diagram](/docs/diagram.png "Diagram")

## How this works

### Caches

We can the following directories by default:

* **Composer** - /root/.composer
* **Yarn** - /usr/local/share/.cache/yarn

## Installation

### Prerequisites

- A Kubernetes cluster running v1.6 or later.
  - [Google Container Engine (GKE)](https://cloud.google.com/container-engine/) is a managed kubernetes-as-a-service procided by Google Cloud Platform (there is a free tier!).
  - [Kops](https://github.com/kubernetes/kops) is a tool for simplifying the management of DIY kubernetes clusters.
  - [Kubernetes the Hard Way](https://github.com/kelseyhightower/kubernetes-the-hard-way) is for those who want to manually configure every component of the cluster.
- SSH key-pair for fetching dependencies from private repositories.
- TBC
  
### Install m8s API components

- TBC

## Development

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
