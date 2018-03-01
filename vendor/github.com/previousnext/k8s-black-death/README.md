Kubernetes: Black Death
=======================

[![CircleCI](https://circleci.com/gh/previousnext/k8s-black-death.svg?style=svg)](https://circleci.com/gh/previousnext/k8s-black-death)

**Maintainer**: Nick Schuch

This tool comes from a long line of "container killer" projects.

* Docker Cleanup - Deletes old containers on a single Docker host (Interal POC)
* [ECS Reaper](https://github.com/previousnext/ecs-reaper) - Cleans up old ECS Tasks 

and now....

**K8s Black Death**

## How it works

Black death is a "cluster deployed" daemon which checks for pods which have been
annotated (marked) with "black-death".

The annotation should contain a unix timestamp X days in the future. Once its time
has come, our daemon will delete that resource.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: victim
  annotations:
    black-death: 1502666026
spec:
  containers:
  - image: nginx
```

## Supported K8s Objects

* Pod
* ReplicaSet
* Deployment
* Service
* Ingress

## Resources

* [Dave Cheney - Reproducible Builds](https://www.youtube.com/watch?v=c3dW80eO88I)

## Development

### Tools

* **Dependency management** - https://github.com/golang/dep
* **Build** - https://github.com/mitchellh/gox
* **Linting** - https://github.com/golang/lint

### Workflow

(While in the `workspace` directory)

**Installing a new dependency**

After adding to imports in codebase run:

```bash
dep ensure
```

**Running quality checks**

```bash
make lint test
```

**Building binaries**

```bash
make build
```
