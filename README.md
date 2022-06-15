# Token Monkey

A simple test http server that you can use to simulate a real-lie auth token service for Stream Chat.

The http server will randomly make fail such as returning errors, invalid tokens or already expired tokens.

## How to run it

1. Clone this repo
2. Make sure you have a recent Golang compiler installed
3. Run the HTTP server
```bash
go run main.go -secret $STREAM_API_KEY
```

This will run an HTTP server on :8090 and return a JSON response with a token for request on the `/token` path

```bash
curl http://localhost:8090/token?userID=jack
```
