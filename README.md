# cowsay-http

This is a HTTP server that wraps around the classic `fortune` and `cowsay`
commands and exposes them as an HTTP API that can be accessed via HTTP clients
like a web browser or curl.

## HTTP API

```
GET /

GET /cowsay -- Does 'fortune | cowsay' by default (customize with URL parameters)
  URL PARAMS
    s string  // Thing to say (defaults to fortune command)
    cf string // Specify a cowfile (add l param to list available cowfiles)
    r bool    // Pick a random cowfile
    l bool    // List all cowfiles available
    // Additional cows flags
    b bool    // Cow appears bored
    d bool    // Cow appears dead
    g bool    // Cow appears greedy
    p bool    // Cow appears paranoia
    s bool    // Cow appears st0ned
    t bool    // Cow appears tired
    w bool    // Cow appears wired (not tired)
    y bool    // Cow appears youthful

ALIASES for /cowsay include
  /say
  /cs

EXAMPLES:
  cows.rest/cowsay
  cows.rest/cowsay?r
  cows.rest/cowsay?d&s=0xDEADBEEF
  cows.rest/cs?s=moo%20world

TIP:
  # URL escape strings with perl or python:
  perl -nE 'use URI::Escape; chomp $_; print(uri_escape($_))' <<< "some long random text"
  python -c 'import urllib.parse; print(urllib.parse.quote(input()))' <<< "some long random text"

  curl "cows.rest/cs?r&s=$(<url encoded string>)"

GITHUB:
https://github.com/lemonase/cowsay-http
```

It is not technically a REST API, but I already registered the domain `cows.rest`
for the cheap so we are running with that.

## Running

You can use Docker or run and build locally with Go

### Docker

```
docker run -p 8091:8091 ghcr.io/lemonase/cowsay-http:master
```

### Go

```
git clone https://github.com/lemonase/cowsay-http.git
cd cowsay-http
go build -o cowsay-http
./cowsay-http
```

## Acknowledgements

cow{say,think} version 3.03, (c) 1999 Tony Monroe
GPLv3 / Artistic License


## Misc

Github has this 'cowsay' like API:

- https://api.github.com/octocat?s=octocow

Other cowsay projects

- [apjanke's fork of the classic cowsay project](https://github.com/cowsay-org/cowsay)
- [cowsay-files](https://github.com/paulkaefer/cowsay-files)
- [Neo-cowsay](https://github.com/Code-Hex/Neo-cowsay)
