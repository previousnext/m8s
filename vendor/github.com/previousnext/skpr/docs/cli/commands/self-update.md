# self-update

Updates the local `skpr` binary to the latest version.

```bash
skpr self-update

Flags:
  --yes           Skips confirmation prompt
  --edge          Prefer beta releases over stable
```

## Examples

Update to latest stable version.

```bash
$ skpr self-update --yes

Checking for latest suitable release ... done

Do you want to update to version v2.0.3?
  Skipping confirmation due to --yes flag.

Downloading (42.2 MB) ... done

Self update complete. Run this command to confirm:
  skpr version
```
