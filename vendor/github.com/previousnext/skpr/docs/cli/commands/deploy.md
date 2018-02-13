# deploy

Triggers deployment of a packaged version to the specified environment.

<pre>
<code>{{ "skpr deploy --help 2>&1" | exec }}</code>
</pre>


This command is asynchronous by default. Use the `--wait` flag if you need the deployment to complete before proceeding (ie in a bash script), or if you want to review the output.

## Examples

Deploy `v1.0.3` to prod environment and wait for deployment to complete.
```
$ skpr deploy prod v1.0.3 --wait

Deploying previousnext/pnx-d8:v1.0.3 to prod...
>> prod...

/data/bin/drush -r /data/app -l http://127.0.0.1 updb -y
 [success] No database updates required.
/data/bin/drush -r /data/app -l http://127.0.0.1 entity:updates -y
 [success] No entity schema updates required
 [success] Cache rebuild complete.
 [success] Finished performing updates.
/data/bin/drush -r /data/app -l http://127.0.0.1 en drush_cmi_tools
 [notice] Already enabled: drush_cmi_tools
/data/bin/drush -r /data/app -l http://127.0.0.1 cimy -y --source=/data/config-export --install=/data/config-install --delete-list=/data/drush/config-delete.yml
/data/bin/drush -r /data/app -l http://127.0.0.1 cr
 [success] Cache rebuild complete.
```
