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

#### Running

```
docker run -p 8091:8091 jamesdixon/cowsay-http:latest-amd64
```

#### Building

```
docker buildx build --platform linux/amd64 -t jamesdixon/cowsay-http:latest-amd64 .
```

#### Pushing

```
docker push jamesdixon/cowsay-http:latest-amd64
```

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

Equivalent to `fortune | cowsay`

```
curl 'http://localhost:8091/cowsay'
```

With a random cow

```
curl 'http://localhost:8091/cowsay?r'
```

With your own text

```
curl http://localhost:8091/cs?r&s=yoohoo
```

You will have to url encode characters like spaces (browsers do this automatically)

```
http://localhost:8091/cs?r&s=give%20me%20some%20space
```
