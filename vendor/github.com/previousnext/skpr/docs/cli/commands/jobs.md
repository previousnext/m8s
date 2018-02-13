# jobs

Lists the *jobs* that have run for a specified environment.

```
skpr jobs [environment]

Arguments:
  [environment]   Environment as configured in skipper. Usually one of  prod/staging/dev.

Flags:
  --format        Output format - supported formats are "text" (default) and "json".
  --job           Name of the cronjob type to return, eg "drush", "queue" or "mail"
```

For each environment the following metadata is shown:

* Kubernetes job identifier
* Status
    * `Running` - Job is in progress.
    * `Successful` - Job command completed.
    * `Failed` - Job command encountered an error.
* Start datetime
* Duration of execution
* Command executed

Only the last 10 executions of each job are retained, so the length of job history depends on the frequency of the cronjob.

## Examples

List the completed jobs for `prod` environment.

```bash
$ skpr jobs prod

NAME                  	STATUS     	STARTED                          	DURATION                 	COMMAND                                                                         
----------------------	-----------	---------------------------------	-----------              	---------------------------------                                               
prod-backup-1511096400	Successful 	Mon, 20 Nov 2017 00:00:13 AEDT   	6m5s                     	/bin/bash -c drush sql-dump --result-file=/tmp/drupal.sql && tar -czf           
                      	           	                                 	                         	/tmp/drupal.tar.gz /tmp/drupal.sql /data ; backup /tmp/drupal.tar.gz pnx-backups
prod-drush-1512010800 	Successful 	Thu, 30 Nov 2017 14:09:54 AEDT   	35s                      	/bin/bash -c drush cron
prod-mail-1512043800  	Running    	Thu, 30 Nov 2017 23:10:16 AEDT   	-2562047h47m16.854775808s	/bin/bash -c drush queue-run --time-limit=30 queue_mail
```

List the latest "mail" jobs for `prod` environment in json format.

```bash
$ skpr jobs prod --job mail --format json

[{
  "name": "prod-mail-1511096400",
  "status": "successful",
  "started": "Mon, 20 Nov 2017 00:00:13 AEDT",
  "duration": "6m5s",
  "command": "/bin/bash -c drush queue-run --time-limit=30 queue_mail"
},
{
  "name": "prod-mail-1512043800",
  "status": "running",
  "started": "Thu, 30 Nov 2017 23:10:16 AEDT",
  "duration": "",
  "command": "/bin/bash -c drush queue-run --time-limit=30 queue_mail"
}]
```
