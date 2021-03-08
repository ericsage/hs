HS
==

Displays a hardcoded list of cards defined by the table handler in `main.go`.
The listening address is also defined in main.go

Build
-----

Install the latest version of Go from (the offical Go site)[https://golang.org/doc/install].

Run go build in the project directory:
```
go build
```

Running
-------

Add a client id and secret in `main.go`:

```
const(
  ...
  clientID     = "{your_id}"
  clientSecret = "{your_secret}"
)
```

Run the binary:
```
./hs
```
