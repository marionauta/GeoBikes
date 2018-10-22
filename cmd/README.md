## How to use the Library

This is a working example on how to use the library.

## Requirements

To use this library you need an API key from [JCDecaux][1], and store it and the
contract in your source code. One way to do it may be:

```go
package main

const (
	contract = "CONTRACT"
	token    = "APIKEY"
)
```

## Use the binary

Just compile the code with `go build -o geobikes` and run the program:

    ./geobikes

This will open the `80` port so you can access from a browser. If you don't have
permissions to open that port, you can use the `--port` flag:

    ./geobikes --port 8080

[1]: https://developer.jcdecaux.com/#/opendata/vls?page=getstarted
