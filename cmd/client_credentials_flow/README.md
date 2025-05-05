## POST /register

reference: https://datatracker.ietf.org/doc/html/rfc7591

```
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "Sample app",
    "redirect_uris": ["http://localhost:3000/callback"]
  }'
```

## POST /token

```
curl -X POST http://localhost:8080/token \
  -u "my-client-id:my-client-secret" \
  -H "Content-Type: application/json" \
  -d '{
    "grant_type": "client_credentials",
    "scope": "read write"
  }'
```
