# cowsay-http

This is a HTTP server that wraps around the classic `fortune` and `cowsay`
commands and exposes them as an HTTP API that can be accessed via browser
or curl.

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

```
docker run -p 8091:8091 jamesdixon/cowsay-http
```

## API

```
 ________________________________ 
< Hi, Welcome to Cowsay HTTP API >
 -------------------------------- 
       \   ,__,
        \  (oo)____
           (__)    )\
              ||--|| *

GET / -- Returns this page


GET /fortune -- Returns a fortune with an optional pipe to cowsay
  URL PARAMS
    cowsay bool -- Toggle cowsay
      randomCow bool -- Toggle random cowfile
      cowfile string -- Specify a cowfile
      borg bool
      dead bool
      greedy bool
      paranoia bool
      stoned bool
      tired bool
      wired bool
      youthful bool
    time bool -- Print time in response


GET /listCows -- Returns a list of available cows
```

## Examples

Equivalent to `fortune | cowsay`

```
curl 'http://localhost:8091/fortune?cowsay'
```

With a random cow

```
curl 'http://localhost:8091/fortune?cowsay&randomCow'
```
