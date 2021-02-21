
## Start the local HTTP Server:
Every HTTP-Request is send to the specified nats subject.
`$./tunats -serve -sub bla -port 8090 -nurl "nats://my.nats.host" -creds /path/to/user.creds`

## Start the 'forwarder':
The HTTP requests are restored and sent to the configured URL.
`$./tunats -forward -sub bla -to "https://de.wikipedia.org" -nurl "nats://my.nats.host" -creds /path/to/user.creds`
