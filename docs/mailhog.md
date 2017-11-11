Mailhog
=======

Given your have a Mailhog container running on port "8025" and on path "/mailhog"...

eg.

```yaml
search:
  image: previousnext/mailhog:latest
```

You will be able to access Mailhog on the ephemeral environment on path:

```
http://example.com/mailhog
```