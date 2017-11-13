Cache
=====

Cached directories are simply directories which are shared between builds.

Given a cached directory directly related to a Kubernetes PVC, we define these as a flag on the M8s API:

```
m8s server --cache-dirs=composer:/root/.composer,yarn:/usr/local/share/.cache/yarn
```

OR via the environment variable:

```
M8S_CACHE_DIRS
```