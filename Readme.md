This program is a NATS-HTTP Proxy that allows you to forward HTTP requests to NATS, process them, and send them back as HTTP responses. It can also act as an HTTP server, listening for incoming requests and forwarding them as NATS messages.

## Features

-   Forward HTTP requests to NATS and back
-   Act as an HTTP server, listening for incoming requests
-   Support for TLS and customizing timeout settings
-   Ability to use NATS credentials for authentication
-   Configurable NATS and HTTP server settings

## Dependencies

-   [nats.go](https://github.com/nats-io/nats.go)
-   [msgpack](https://github.com/vmihailenco/msgpack)

## Usage

To run the program, first compile it with `go build`. You can customize its behavior using command-line flags.

### Flags

-   `-port`: The HTTP server port (default: 8090)
-   `-serve`: Run the HTTP server (default: false)
-   `-sub`: The NATS subject to listen on/send to
-   `-forward`: Forward NATS messages to an HTTP server (default: false)
-   `-insecureSkipVerify`: Allow insecure TLS connections (default: false)
-   `-to`: The remote HTTP server URL
-   `-nurl`: The NATS cluster URL (default: nats://localhost:4222)
-   `-creds`: The path to the NATS credentials file
-   `-nats_name`: NATS connection name

### Forwarding HTTP requests to NATS

To forward HTTP requests to a NATS subject, run the program with `-forward` flag and specify the NATS subject and remote HTTP server URL using `-sub` and `-to` flags respectively:

```bash
./tunats -forward -sub "http_subject" -to "http://remote-http-server.com"
```


### Running the HTTP server

To run the HTTP server and forward incoming requests to a NATS subject, use the `-serve` flag and specify the NATS subject:

```bash
./nats-http-proxy -serve -sub "http_subject"
```

### Connecting to a secure NATS cluster

To connect to a NATS cluster that requires authentication, use the `-creds` flag to provide the path to the NATS credentials file:

```bash
./nats-http-proxy -serve -sub "http_subject" -creds "/path/to/your/nats.creds"
```

### Ignoring TLS certificate verification

To disable TLS certificate verification (not recommended for production environments), use the `-insecureSkipVerify` flag:

```
./nats-http-proxy -serve -sub "http_subject" -insecureSkipVerify
```

## License

This project is released under the MIT License. See the `LICENSE` file for more information.