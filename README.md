# :dog: bantay

Lightweight uptime monitor for web services

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Write a `checks.yml` ([see section below](#example-checks.yml)) to define the uptime checks you want to run, along with settings and reporters. Then, running it then should be as simple as:

```console
$ make
$ vim checks.yml
$ ./bin/bantay
```

## Example `checks.yml`

```yaml
---
server:
  poll_interval: 10
checks:
  - name: Google
    url: https://www.google.com/
    valid_status: 200
  - Hacker News
    url: https://news.ycombinator.com/
    valid_status: 200
    body_match: Hacker News
reporters:
  - type: log
  - type: Slack
    options:
      slack_channel: YOUR-SLACK-CHANNEL-HERE
      slack_token: YOUR-SLACK-TOKEN-HERE
```
