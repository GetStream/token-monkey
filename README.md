# Token Monkey

A simple test http server that you can use to simulate a real-lie auth token service for Stream Chat.

The http server will randomly make fail such as returning errors, invalid tokens or already expired tokens.

## How to run it

1. Clone this repo
2. Make sure you have a recent Golang compiler installed
3. Run the HTTP server
```shell
go run main.go -secret STREAM_API_KEY
```

This will run an HTTP server on :8090 and return a JSON response with a token for request on the `/token` path

```shell
curl http://localhost:8090/token?userID=jack
```

## Outcomes

By default, the server will pick a random outcome from this list:

- `serverTimesOut` server replies with empty body after 1 hour
- `serverErrors` server errors with a 500 HTTP error
- `returnsExpiredToken` server returns a valid token but with the `exp` claim in the past
- `returnsValidToken` server returns a valid token with 60s expiration
- `returnsGarbageToken` server returns a valid JWT token but with invalid claims

You can pick different outcomes when you start the server:

```shell
// will only return 500 errors and timeouts
go run main.go -secret STREAM_API_KEY -outcomes serverErrors,serverTimesOut
```
