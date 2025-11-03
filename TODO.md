# Things to do

This works to show that CORS consoluted for every redirect

```yaml
# CORS Configuration for example
PROXY_ADDRESS="example_proxy"
PROXY_PORT="8080"
REDIRECT_ADDRESS="localhost"
REDIRECT_PORT="8082"
ALLOW_CREDENTIALS="true"
ALLOWED_ORIGINS="http://localhost:8080"
ALLOWED_METHODS="POST,GET,OPTIONS"
ALLOWED_HEADERS="Authorization"

# CORS configuration for example_proxy
PROXY_ALLOW_CREDENTIALS="true"
PROXY_ALLOWED_ORIGINS="http://localhost:8080"
PROXY_ALLOWED_METHODS="POST,GET,OPTIONS"
PROXY_ALLOWED_HEADERS="Authorization"
```
