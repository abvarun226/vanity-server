# Go Vanity Server

This is a server for vanity URLs for Go packages, written in Go.


## How to build?
```
$ go build -o vanity-server .
```

## How to run?
```
$ ./vanity-server -h
Usage of ./vanity-server:
  -config string
    	Contains go-import mapping rules (default "/etc/go-vanity/config.json")
  -godocurl string
    	The godoc URL
  -listen string
    	Address where this server listens (default ":8080")

$ ./vanity-server -config "./example-config.json" -listen ":80"
```

Make sure to put this server behind Nginx or Traefik, and use SSL certificates to secure the endpoints.

There is an `example-config.json` config file that provides an example of the go-import mapping rules.