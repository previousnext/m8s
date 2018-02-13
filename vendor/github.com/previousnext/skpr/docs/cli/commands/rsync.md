# rsync

Rsync files between local and remote environments.

```bash
skpr rsync [source] [destination]

Arguments:
  [source]        File/directory source location.
                  See format documentation for details.
  [destination]   File/directory destination.
                  See format documentation for details.
```

> Note: This command can not sync between 2 remote environments. Only local -> remote and vice-versa.

## Format

As `rsync` is used under the hood to perform the sync, the same source/destination formatting rules apply. Skipper adds a thin layer which allows you to use the environment name in the place of a hostname.

## Examples

Sync public files directory from `prod` to local.

```bash
$ skpr rsync prod:app/sites/default/files/ app/sites/default/files
```

Sync private files from local to `staging`.

```bash
$ skpr rsync ./private/ staging:/private
```
