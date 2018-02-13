# sync

Synchronise the DB and files between 2 skipper environments.

```bash
skpr sync [source] [destination]

Arguments:
  [source]        Source environment as configured in skipper.
  [destination]   Source environment as configured in skipper.

Flags:
  --wait          Waits for the sync process to complete before exiting.
  --resource      Resource(s) to sync. Options are:
                  - all (default)
                  - db
                  - db-[db name]
                  - files
                  - files-[mount name]
```

This command is asynchronous by default. Use the `--wait` flag if you need the sync to complete before proceeding (ie in a bash script).

## Examples

Sync all resources from `prod` to `staging`.

```bash
$ skpr sync prod staging
```

Sync _only_ db from `prod` to `staging`, waiting for the process to complete.

```bash
$ skpr sync prod staging --resource db --wait
```

Sync _only_ public files from `staging` to `dev`.

```bash
$ skpr sync staging dev --resource files-public
```
