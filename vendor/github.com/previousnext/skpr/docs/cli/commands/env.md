# env

`skpr env` has several sub commands.

* [env delete](#delete)
* [env get](#get)
* [env list](#list)

## delete

Deletes the specified environment.

```
skpr env delete [environment]

Arguments:
  [environment] Environment as configured in skipper. Usually one of  prod/staging/dev.

Flags:
  --dry-run     Displays list of resources that would be deleted take place
  --yes         Skips confirmation prompt (except if environment is `prod`)
```

The following resources are deleted **immediately**.

* Application pod(s)
* Kubernetes config map (see [`skpr config`](config.md))
* Cron jobs
* Secrets
* Solr core(s)
* Elasticache instance(s)
* MySQL account(s)

The following resources are marked for deletion **within 24 hours**.

* EFS Volumes (public / private file mounts)
* MySQL database(s)
* CloudFront distribution

Attempting to delete the `prod` environment will *always* prompt for confirmation.

### Safeguards

There are a few safeguard mechanisms in place to mitigate accidental deletion of environments.

* The default behaviour prompts the user for confirmation after displaying the list of resources that would be deleted. The user must enter the name of the environment.
* The `--yes` flag bypasses the confirmation step, but is _always ignored_ for the prod environment.
* If there is accidental deletion and the mistake is noticed before the nightly cleanup job, the environment can be restored without data loss or a DNS change. There will be an outage.
* Our uptime monitoring will alert SysOps if the prod environment disappears unexpectedly.

### Examples

Delete environment called `dev01` without confirmation prompt.

```bash
$ skpr env delete dev01 --yes

Building list of resources to delete.

SYSTEM  RESOURCE TYPE   IDENTIFIER
k8s     replicaset      dev01-779797126
k8s     configmap       dev01
k8s     cron            dev01-drush
k8s     cron            dev01-backup
aws     cloudfront      EDFDVBD632BHDS5
aws     efs             dev01-public [fs-085eb930]
aws     efs             dev01-private [fs-1b5ec922]
mysql   database        foo_dev01
mysql   account         foo_dev01

Are you sure you want to PERMANENTLY DELETE these resources? Please enter 'dev01' to confirm.
  Skipping confirmation due to --yes flag.

Deleting environment dev01...
  - Deleting k8s/replicaset/dev01-779797126 ... done
  - Deleting k8s/configmap/dev01 ... done
  - Deleting k8s/cron/dev01-drush ... done
  - Deleting k8s/cron/dev01-backup ... done
  - Scheduling deletion of aws/cloudfront/EDFDVBD632BHDS5 ... done
  - Scheduling deletion of aws/efs/dev01-public ... done
  - Scheduling deletion of aws/efs/dev01-private ... done
  - Scheduling deletion of mysql/database/foo_dev01 ... done
  - Deleting mysql/account/foo_dev01 ... done
```

Preview what a delete operation would do on `staging` environment.

```bash
$ skpr env delete staging --dry-run

Building list of resources to delete.

SYSTEM  RESOURCE TYPE   IDENTIFIER
k8s     replicaset      staging-779797126
k8s     configmap       staging
k8s     cron            staging-drush
k8s     cron            staging-backup
aws     cloudfront      ADEDVBD632BHDK9
aws     efs             staging-public [fs-085eb930]
aws     efs             staging-private [fs-1b5ec922]
mysql   database        bar_stg01
mysql   account         bar_stg01

Are you sure you want to PERMANENTLY DELETE these resources? Please enter `staging` to confirm.
  Aborted due to --dry-run flag.
```

## get

Display information for skipper environments.

```bash
skpr env get [environment]

Arguments:
  [environment]   Environment to fetch information for.

Flags:
  --format        Output format - supported formats are "text" (default) and "json".
```

The following metadata is provided for each environment:

* Hostnames being routed to the environment
* Currently deployed version
* Infrastructure resource allocation and current utilisation.
* Number of app instances running.

### Examples

Display information for `prod` environment in `json` format.

```bash
$ skpr env get prod --format=json

{
  "environment": "prod",
  "domain": [
    "pnx-d8-prod.cd.pnx.com.au",
    "previousnext.com.au",
    "www.previousnext.com.au"
  ],
  "version": "1.1.1",
  "size": {
    "min": 256,
    "max": 512,
    "proc": 128
  },
  "status": {
    "desired": 2,
    "healthy": 2
  },
  "utilisation": {
    "cpu": 0
  }
}
```

Get hostname used for `staging` environment.

> Requires [jq](https://github.com/stedolan/jq) binary to parse the json response.

```bash
$ skpr env get staging --format=json | jq '.domain[0]'

"pnx-d8-stage.cd.pnx.com.au"
```

## list

List information for all environments.

```bash
skpr env list

Flags:
  --format        Output format - supported formats are "text" (default) and "json".
```

### Examples

List all environments.

```bash
$ skpr env list

ENVIRONMENT	DOMAIN                      	VERSION           	SIZE            	STATUS     	CPU (Avg)  
-----------	----------------------      	-----------       	-----------     	-----------	-----------
dev        	pnx-d8-qa.cd.pnx.com.au     	1.0.23-57-g5679131	256 / 512 / 128 	1 / 1      	0%
prod       	pnx-d8-prod.cd.pnx.com.au   	1.1.1             	256 / 1024 / 128	2 / 2      	0%
           	previousnext.com.au         	                  	                	           	           
           	www.previousnext.com.au     	                  	                	           	           
staging    	pnx-d8-staging.cd.pnx.com.au	1.1.1             	256 / 512 / 128 	1 / 1      	0%

```
