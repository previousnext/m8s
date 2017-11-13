Google Container Engine (GKE)
=============================

## Overview

The following document is for creating a M8s API backed by Kubernetes and Google Cloud Platform.

Known limitations:

* We don't cache composer or yarn packaages, we don't have a solution for "easy" network storage on GCP.

## Bootstrapping a cluster

* Google Cloud Platform Account
* Billing setup for GKE

### Create a new cluster

```bash
gcloud container clusters create demo --num-nodes=2
```

### Inspect the cluster

```bash
$ gcloud container clusters list
NAME  ZONE           MASTER_VERSION  MASTER_IP       MACHINE_TYPE   NODE_VERSION  NUM_NODES  STATUS
demo  us-central1-b  1.7.8-gke.0     35.202.108.243  n1-standard-1  1.7.8-gke.0   2          RUNNING
```

### Setup Commandline

#### Install the Kubernetes CLI

```bash
gcloud components install kubectl
```

#### Set the context/credentials

```bash
gcloud container clusters get-credentials demo
```

## Deploying cluster components

### Traefik

Traefik is a software loadbalancer with native Kubernetes integration.

This project provides us with access to our build environments.

#### Deploy

See the example k8s manifest here - [traefik.yaml](traefik.yaml).

```bash
kubectl create -f traefik.yaml
```

#### Setup DNS

Once the deployment is finished we should see that:

* The pod is running
* Our service was given an external IP address

```bash
$ kubectl -n kube-system get pods -o wide | grep traefik
traefik-3932826267-szx2b   1/1   Running  0  6m   10.16.1.3   gke-demo-default-pool-da95f638-j8n4
```

```bash
$ kubectl -n kube-system get svc traefik
NAME      TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)        AGE
traefik   LoadBalancer   10.19.250.216   35.202.167.21   80:31702/TCP   7m
```

To finish off this Traefik deployment should point a wildcard domain at the `EXTERNAL-IP`

eg.

```
*.m8sdemo.com -> 35.202.167.21
```

This will allow developers to pick any domain under `m8sdemo.com` for their ephemeral sites:

```
site1.d8sdemo.com
```

## Deploying M8s API

### Requirements

* **Token** - A random string used for M8s CLi -> API authentication
* **LetsEncrypt** - TLS is required for this endpoint
* **Docker Account** - Credentials used for pulling private images
* **SSH Keypair** - A new keypair for allowing environment to ssh remote endpoints

The above are values which will need to be changed in the `m8s.yaml` file.

#### Create a namespace

This allows us to group our API and built environments.

```bash
kubectl create ns m8s
```

#### Deploy

First you will need to update all the references of `CHANGE_ME` in the file: `m8s.yaml`.

When complete you will be able to deploy the M8s API with:

```bash
kubectl create -f m8s.yaml
```

This will deploy the following components:

* API
* LetsEncrypt cache volume
* Service to expose to the outside world
* Role for allowing our API to access the Kubernetes API
* Caching strategy for "composer" and "yarn"

#### Setup DNS

We now want to get the `EXTERNAL-IP` from the service and create DNS record for our API which matches the `LETS_ENCRYPT_DOMAIN` env variable set in the `m8s.yaml`.

```bash
$ kubectl -n m8s get svc
NAME      TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)         AGE
api       LoadBalancer   10.19.255.42   104.197.224.56   443:31321/TCP   2m
```

```
api.m8sdemo.com -> 104.197.224.56
```

## Congratulations

You are now ready to setup a Pipeline to create new environments!
