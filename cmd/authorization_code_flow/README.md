## GET /.well-known/openid-configuration

```
curl http://localhost:8080/.well-known/openid-configuration
```

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

## POST /signup

```
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securePassword123"}'
```

## POST /authorize

```
curl -G http://localhost:8080/authorize \
  --data-urlencode "client_id=my-client-id" \
  --data-urlencode "redirect_uri=http://localhost:3000/callback" \
  --data-urlencode "scope=openid"
```

## POST /login

```
curl -vL -X POST http://localhost:8080/login \
  -d "email=user@example.com" \
  -d "password=securePassword123" \
  -d "client_id=my-client-id" \
  -d "redirect_uri=http://localhost:3000/callback" \
  -d "scope=openid"
```


## POST /token

```
curl -X POST http://localhost:8080/token \
  -u my-client:my-secret \
  -d grant_type=authorization_code \
  -d code=abc123 \
  -d redirect_uri=http://localhost:3000/callback
```
