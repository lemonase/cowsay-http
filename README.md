# cowsay-http

This is a HTTP server that wraps around the classic `fortune` and `cowsay`
commands and exposes them as an HTTP API that can be accessed via HTTP clients
like a web browser or curl.

## HTTP API

```
Welcome to the cowsay HTTP API!

GET /* -- This page (you are here!)

GET /cowsay -- Does 'fortune | cowsay' by default (customize with URL parameters)
  URL PARAMS
    s string -- Thing to say (defaults to fortune command)
    cf string -- Specify a cowfile (add l param to list available cowfiles)
    r bool -- Pick a random cowfile
    l bool -- List all cowfiles available

ALIASES for /cowsay include
  /say
  /cs

EXAMPLES:
  cows.rest/cowsay?r
  cows.rest/cs?s=moo%20world

GITHUB:
https://github.com/lemonase/cowsay-http
```

It is not technically a REST API, but I already registered the domain `cows.rest`
for the cheap so we are running with that.

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
