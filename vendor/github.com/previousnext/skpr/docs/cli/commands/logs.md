# logs

Stream logs to your terminal. The command supports both start/end dates, and ability to stream logs in real-time.

```bash
skpr logs [environment]

Arguments:
  [environment]   Environment to pull logs from

Flags:
  --start         Time of first of logs to retrieve - See supported formats below.
                  Defaults to "15m" (or "now" if --follow enabled).
  --end           Time of final logs to retrieve - See supported formats below.
                  Defaults to "now".
                  Conflicts with --follow.
  --source        Filter logs to a specific system source. Comma-separated multi values supported.
                  Supports "app" (default).
                  Future support for "cron", "mysql-slow", "cloudfront".
  --follow        Streams logs in real-time until interrupted (ctrl+c).
                  Conflicts with --end.
  --timezone      Timezone to use for --start and --end flags. See formats below.
                  Defaults to local system timezone.
                  Has no effect when relative dates provided.

Date/Time Formats:
  Relative: 30s, 15m, 2h, 1d (hopefully suffixes are self-explanatory)
  Absolute:
    '19:05': time with no date, defaults to current date (use 24 hour syntax)
    '2017-6-15': date with no time - assumes midnight of specified date
    '2017-6-15 19:05': date and time specified.
  Timezones: Use the IANA time zone format, eg "Australia/Sydney".
```

## Examples

Get `prod` logs generated between two times on a **specific date**.

```bash
$ skpr logs prod --start '2017-10-11 9:00' --end '2017-10-11 10:30'
```

Get `staging` logs for the **last hour**.

```bash
$ skpr logs staging --start 1h
```

Stream `dev` logs in **real time**.

```bash
$ skpr logs dev --follow
```

Get `prod` logs generated between two times **today** (Sydney time).

```bash
$ skpr logs prod --start '2017-10-11 11:55' --end '2017-10-11 12:05' --timezone "Australia/Sydney"
```

Get `dev` cloudfront logs generated between two **relative times**.

```bash
$ skpr logs dev --start 2h --end 30m --source cloudfront
```
