# consul-simple-notifier

A notification command for use with consul-alerts

[![Build Status](https://travis-ci.org/udzura/consul-simple-notifier.svg)](https://travis-ci.org/udzura/consul-simple-notifier)

## Settings

Just place `/etc/consul-simple-notifier.ini` (or `-c` option to pass custom location).

```conf
[email]
recipients = [
  "udzura+notify@udzura.jp",
  "udzura@example.com"
]

[ikachan]
url = "http://irc.example.com:4979"
channel = '#udzura_test'

```

Current available notifier is `/bin/mail` command and [ikachan](https://github.com/yappo/p5-App-Ikachan) hook.
