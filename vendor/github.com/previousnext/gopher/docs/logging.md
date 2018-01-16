Logging
=======

Logging should use the package `github.com/prometheus/common/log`

This results in a log like:

```
INFO[0000] Serving Prometheus metrics endpoint           source="main.go:23"
```

Notice the `source="main.go:23"`, this allows us to track down the line of code.