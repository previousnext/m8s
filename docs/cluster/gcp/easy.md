Google Container Engine (GKE)
=============================

## Overview

In the following document we will be deploying:

* Google Kubernetes Engine
* M8s API (+ components) using an inbuilt installation utility

Known limitations:

* We don't cache composer or yarn packages, we don't have a solution for "easy" network storage on GCP.

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

Given Google Kubernetes Engine automatically sets your Kubernetes config to the built cluster, we can
look at deploying our components onto the cluster.

## Installing M8s

M8s has the ability to install all the components it requires onto a Kubernetes cluster.

To run an install you will first need to download the latest M8s binary from the [releases](https://github.com/previousnext/m8s/releases) page.

To install M8s all you need to provide is:

* Random token string for authentication
* Domain you wish the API to respond on
* Email address of the operations team (for Lets Encrypt support)

```bash
$ m8s install --token=12345678 \
              --domain=api.m8sdemo.com \
              --email=nick@m8sdemo.com
Installing: Namespace
Installing Traefik: Deployment
Installing Traefik: Service
Installing M8s API: Deployment
Installing M8s API: PVC
Installing M8s API: Secret
Installing M8s API: Servce
Deployed!
Status: kubectl -n m8s get all
Entrypoints: kubectl -n m8s get svc
```

To verify our deploy we can now run the above commands:

```bash
$ kubectl -n m8s get pods
NAME                       READY     STATUS    RESTARTS   AGE
m8s-api-3914376045-szrdb   1/1       Running   0          3h
traefik-361531620-l69g5    1/1       Running   0          3h
```

```bash
$ kubectl -n m8s get svc
NAME      TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)         AGE
m8s-api   LoadBalancer   10.19.243.152   35.202.127.242   443:30535/TCP   3m
traefik   LoadBalancer   10.19.252.143   35.188.215.243   80:31475/TCP    3m
```

Finally you will need to setup some DNS entries:

* _api.m8sdemo.com_ -> _10.19.243.152_
* _*.m8sdemo.com_ -> _10.19.252.143_

## Congratulations

You are now ready to set up a Pipeline to create new environments!