Metrics
=======

Metrics should be exposed to Prometheus on:

* Port `9000`
* Path `/metrics`

We should use the package `github.com/prometheus/client_golang/prometheus/promhttp` and implemented with the following line:

```
package main

import (
  "log"
  "net/http"

  "github.com/alecthomas/kingpin"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
  cliPrometheusPort = kingpin.Flag("prometheus-port", "Prometheus metrics port").Default(":9000").OverrideDefaultFromEnvar("METRICS_PORT").String()
  cliPrometheusPath = kingpin.Flag("prometheus-path", "Prometheus metrics path").Default("/metrics").OverrideDefaultFromEnvar("METRICS_PATH").String()
)

func main() {
  go metrics(*cliPrometheusPort, *cliPrometheusPath)
}

// Helper function for serving Prometheus metrics.
func metrics(port, path string) {
  http.Handle(path, promhttp.Handler())
  log.Fatal(http.ListenAndServe(port, nil))
}
```