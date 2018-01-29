UI
==

## Overview

The m8s user interface was created to facilitate the following requirements:

* Ensure that developers don't need to install any CLI tools
* Central hub to view environments
* Command line access
* Logs access

![UI](ui.png "UI")

## Why not the Kubernetes Dashboard 

The Kubernetes Dashboard is targetted at viewing the entire Kubernetes cluster.

The m8s user interface is designed to run in a single Namespace.

## How to deploy

The following will deploy:

* Deployment
* Service
* Ingress

```bash
$ kubectl create -f ui.yaml
```

Note: We highly recommend you use authentication eg. oauth2_proxy.