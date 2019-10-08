# :dog: bantay

Lightweight uptime monitor for web services, with support for sending alerts through Slack and email (using Mailgun)

## Getting started

On both Docker and your local machine, you need to first make a `checks.yml` ([see section below](#example-checksyml)) to define the uptime checks you want to run, along with settings and reporters.

### Running with :whale: Docker

The simplest way to run `bantay` is with Docker.

```console
$ make package
$ docker run -v "$(pwd)/checks.yml":/opt/bantay/bin/checks.yml fipanganiban/bantay:local bantay check
```


### Running on your local machine

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Then, running it then should be as simple as:

```console
$ go get .
$ make build
$ ./bin/bantay check
```
to run the checks once, or
```console
$ ./bin/bantay server
```
to run checks over and over, on the interval specified in `checks.yml` as `poll_interval`.

## Example `checks.yml`

```yaml
---
server:
  poll_interval: 10
checks:
  - name: Google
    url: https://www.google.com/
    valid_status: 200
  - name: Hacker News
    url: https://news.ycombinator.com/
    valid_status: 200
    body_match: Hacker News
reporters:
  - type: log
  - type: slack
    options:
      slack_channel: YOUR-SLACK-CHANNEL-HERE
      slack_token: YOUR-SLACK-TOKEN-HERE
  - type: mailgun
    options:
      mailgun_domain: YOUR-MAILGUN-DOMAIN
      mailgun_private_key: YOUR-MAILGUN-PRIVATE-API-KEY
      mailgun_sender: bantay@yourdomain.io
      mailgun_recipients: [webmaster@yourdomain.io]
      mailgun_exclude:
        - Hacker News
```

## `checks.yml` Options

### `server` section

Settings used when running bantay in server mode, ie `./bantay server`

- `poll_interval`: How long to wait for each check of all microservices (in seconds)

### `checks` section

List of items that bantay will run, with their check settings.

- `name`: Unique identifier for each check (case sensitive), used for reporting and alerts
- `url`: Absolute URL that bantay will poll each time a check is run
- `valid_status`: HTTP status code to expect from the HTTP response
- `body_match` (optional): String to search for in the HTTP response

### `reporters` section:

List of reporters that bantay will use to report check results, each with their own set of options.

- `type`: Type of reporter to use. Currently supported: `log` (stdout/stderr), `slack`, `mailgun`
- `options`: Options specific for each reporter type (more below)

#### Options:

- `slack`
  - `slack_channel`: ID of the Slack channel to send reports to
  - `slack_token`: Slack token of the bot/user to use
- `mailgun`:
  - `mailgun_domain`: Your Mailgun registered domain
  - `mailgun_private_key`: Private API key to access the Mailgun V3 API
  - `mailgun_sender`: Email address to show as sender of alerts
  - `mailgun_recipients`: List of email addresses to send emails to
  - `mailgun_exclude`: List of unique check `name`s to exclude from sending email alerts
