# scale

Manually adjust scaling settings for the specified environment.

```bash
skpr scale [environment]

Arguments:
  [environment] Environment as configured in skipper. Usually one of  prod/staging/dev.

Flags:
  --min         Autoscaler minimum.
  --max         Autoscaler maximum.
```

See the [Scaling](/config/sizing.md#scaling) section for details on the how the horizontal pod autoscaler works with the `--min` and `--max` flags.

> Note: --min must be equal or less than --max, and vice-versa.

Any overrides are reset to the defaults defined in `.skpr.yml` at the next deployment.

## Examples

Boost the minimum number of pods to 6.

> Use case: Client is launching a marketing campaign and needs more resources ready to go.

```bash
$ skpr scale prod --min 6 --max 20
```
