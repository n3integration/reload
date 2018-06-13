Reload [ ![Codeship Status for n3integration/reload](https://app.codeship.com/projects/89707150-4e72-0136-7326-322e9f850b54/status?branch=master)](https://app.codeship.com/projects/293485)
[![codecov](https://codecov.io/gh/n3integration/reload/branch/master/graph/badge.svg)](https://codecov.io/gh/n3integration/reload)
========

`reload` is a command line utility for live-reloading Go web applications. It is
a fork of the [gin](https://github.com/codegangsta/gin) project originally written by [Jeremy Saenz](https://github.com/codegangsta).

Just run `reload` in your app directory and your web app will be served with
`reload` as a proxy. `reload` will automatically recompile your code when it
detects a change. Your app will be restarted the next time it receives an
HTTP request (unless the `--immediate` flag is passed).

`reload` adheres to the "silence is golden" principle, so it will only complain
if there was a compiler error or if you successfully compile after an error.

## Installation

Assuming you have a working Go environment and `GOPATH/bin` is in your
`PATH`, `reload` is a breeze to install:

```shell
go get -u github.com/n3integration/reload
```

Then, verify that `reload` was installed correctly:

```shell
reload -h
```
## Basic usage
```shell
reload [options] run main.go
```
Options
```
   --laddr value, -l value       listening address for the proxy server
   --port value, -p value        port for the proxy server (default: 3000)
   --appPort value, -a value     port for the Go web server (default: 3001)
   --bin value, -b value         name of generated binary file (default: "gin-bin")
   --path value, -t value        Path to watch files from (default: ".")
   --build value, -d value       Path to build files from (defaults to same value as --path)
   --excludeDir value, -x value  Relative directories to exclude
   --immediate, -i               run the server immediately after it's built
   --all                         reloads whenever any file changes, as opposed to reloading only on .go file change
   --buildArgs value             Additional go build arguments
   --certFile value              TLS Certificate
   --keyFile value               TLS Certificate Key
   --logPrefix value             Setup custom log prefix
   --notifications               enable desktop notifications
   --help, -h                    show help
   --version, -v                 print the version
```

## Supporting Reload in Your Web App
`reload` assumes that your web app binds itself to the `PORT` environment
variable so it can properly proxy requests to your app.

## Using flags?
When you normally start your server with [flags](https://godoc.org/flag)
if you want to override any of them when running `reload` we suggest you
instead use [github.com/namsral/flag](https://github.com/namsral/flag)
as explained in [this post](http://stackoverflow.com/questions/24873883/organizing-environment-variables-golang/28160665#28160665)
