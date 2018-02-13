# Sizing

This document describes the sizing and scaling

## Scaling

Skipper utilises the Kubernetes [horizontal pod autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) to adjust the number of pods deployed to the cluster.

We define the minimum and maximum pod boundaries, and the HPA adjusts the actual amount based on configurable rules. By default, it is when the average CPU across all the pods reaches 80%.

This can be configured with the following parameters.

**.skpr.yml**
```yaml
environments:
  - name: prod
    spec:
      ...
      # Change CPU scaling threshold to 50%.
      scaleCPU: 50
```
