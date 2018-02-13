# Command Line Interface

## Command Overview

* [config](commands/config.html)
    * [config delete](commands/config.md#delete)
    * [config get](commands/config.md#get)
    * [config set](commands/config.md#set)
    * [config list](commands/config.md#list)
* [delete](commands/delete.md)
* [deploy](commands/deploy.md)
* [env](commands/env.md)
* [exec](commands/exec.md)
* [jobs](commands/jobs.md)
* [logs](commands/logs.md)
* [package](commands/package.md)
* [rsync](commands/rsync.md)
* [scale](commands/scale.md)
* [self-update](commands/self-update.md)
* [shell](commands/shell.md)
* [sync](commands/sync.md)

## Installation

### OSX

```bash
wget http://bins.skpr.io/darwin-amd64-latest.tar.gz
sudo tar -zxf darwin-amd64-latest.tar.gz -C /usr/local/bin/
```

### Linux

```bash
wget http://bins.skpr.io/linux-amd64-latest.tar.gz
sudo tar -zxf linux-amd64-latest.tar.gz -C /usr/local/bin/
```

### Create a Github Personal Access Key

You will need to create a personal access key for the Skipper CLI.

Go to the Personal access tokens page on Github and click **Generate new token**.

Give it a name (e.g. Skpr) and grant it **Repo** access.

### Save your Access Key Locally

Create a new file in your home directory ***~/.skpr.yml*** and add the following:

```yaml
user: "kimpepper" # Your github username
token: "XXXXXXXXXXXXXXXXXXXXXXXXX" # The token created in the step above.
```
