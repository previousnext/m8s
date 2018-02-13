# exec

Executes a bash command on the specified environment.

```bash
skpr exec [environment] [command]

Arguments:
  [environment] Environment as configured in skipper. Usually one of  prod/staging/dev.
  [command]     Command to execute on remote environment.
```

## Examples

Get drush status for `prod` environment.

```bash
$ skpr exec prod drush status

Drupal version   : 8.4.2
Site URI         : default
DB driver        : mysql
Database         : Connected
Drupal bootstrap : Successful
Admin theme      : seven
PHP binary       : /usr/local/bin/php
PHP config       : /usr/local/etc/php/php.ini
PHP OS           : Linux
Drush script     : /data/bin/drush
Drush version    : 9.0.0-beta9
Drush temp       : /tmp
Drush configs    : /data/vendor/drush/drush/drush.yml
                   /data/drush/drush.yml
Drupal root      : /data/app
Site path        : sites/default
Files, Public    : sites/default/files
Files, Private   : /private
Files, Temp      : /tmp
```
