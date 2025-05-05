## Request sample for Client Credentials Flow

```
curl -X POST http://localhost:8080/token \
  -u "my-client-id:my-client-secret" \
  -H "Content-Type: application/json" \
  -d '{
    "grant_type": "client_credentials",
    "scope": "read write"
  }'
```
