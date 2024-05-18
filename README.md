# cowsay-http

This is a HTTP server that wraps around the classic `fortune` and `cowsay`
commands and exposes them as an HTTP API that can be accessed via browser
or curl.

It is not technically a REST API, but I did get the domain for `cows.rest`
so we are running with that.

## HTTP API

```
GET / -- Returns this page

GET /cowsay -- Does cowsay (customize with URL parameters)
GET /cs
  URL PARAMS
    s string -- Thing to say (defaults to fortune command)
    cf string -- Specify a cowfile (see /list or add l param to request)
    r bool -- Pick a random cowfile
    l bool -- List all cowfiles available
```

## Examples

Get `fortune | cowsay`

```
curl cows.rest/cs
```

With random cowfile

```
curl cows.rest/cs?r
```

With text

```
curl cows.rest/cs?r&s=hi%20fren
```

## Running

Clone this repo and run locally with Go, or use docker

### Go

```
git clone https://github.com/lemonase/cowsay-http.git
cd cowsay-http
go build -o cowsay-http
./cowsay-http
```

### Docker

#### Running

```
docker run -p 8091:8091 jamesdixon/cowsay-http:latest-amd64
```

## Acknowledgements

cow{say,think} version 3.03, (c) 1999 Tony Monroe
GPLv3 / Artistic License

## Misc

Github has this 'cowsay' like API:

- https://api.github.com/octocat?s=octocow

Other cowsay projects

- [apjanke's fork of the classic cowsay project](https://github.com/cowsay-org/cowsay)
- [Neo-cowsay](https://github.com/Code-Hex/Neo-cowsay)
- [cowsay-files](https://github.com/paulkaefer/cowsay-files)
