# :dog: bantay

Lightweight uptime monitor for web services

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Write a `checks.yml` to define the uptime checks you want to run, along with settings and reporters. Then, running it then should be as simple as:

```console
$ make
$ vim checks.yml
$ ./bin/bantay
```
