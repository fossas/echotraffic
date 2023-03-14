### `echotraffic`: echo your network traffic for FOSSA CLI

This small program allows users to view traffic sent by `fossa-cli` to the FOSSA service.

### 30 second setup

0. Install Go: https://go.dev/doc/install
1. Run `go run main.go`.
2. Run `fossa analyze -e 'http://localhost:3000'` in a different terminal.
3. Observe `fossa-cli` traffic written to the terminal.

Run `go run main.go -h` for more usage information.

### Building

If desired, you can build and distribute this with `go build`.

### How to read output

Each discrete API call is denoted by the line `🚀 Forward '{path}' to '{destination}'`:
```
✨ Serving on ':3000', forwarding to 'https://app.fossa.com'
🚀 Forward '/api/cli/organization' to 'https://app.fossa.com/api/cli/organization'
GET /api/cli/organization HTTP/1.1
Host: localhost:3000
User-Agent: Go-http-client/1.1
Accept: application/json
Accept-Encoding: gzip
Authorization: Bearer <snip>
```

If a body is attached to the request, it's listed after a blank line:
```
🚀 Forward '/api/proxy/sherlock/scans/bde59fdc-923b-4429-a871-45fbb24bc2ba/files' to 'https://app.fossa.com/api/proxy/sherlock/scans/bde59fdc-923b-4429-a871-45fbb24bc2ba/files'
POST /api/proxy/sherlock/scans/bde59fdc-923b-4429-a871-45fbb24bc2ba/files HTTP/1.1
Host: localhost:3000
User-Agent: Go-http-client/1.1
Content-Length: 261729
Accept-Encoding: gzip
Authorization: Bearer <snip>
Content-Type: application/json

{"ScanData":{"...
```
