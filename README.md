# go-blog-cors (github.com/antonio-alexander/go-blog-cors)

This is a companion repository for an article describing CORS (cross origin resource sharing) with an emphasis on Golang. CORS is a protection enforced by browsers to ensure that calls between browsers and server(s) (especially between web servers) are explicitly allowed. In addition to this document itself, there will also be example source code to demonstrate troubleshooting and development. After review of this repository, you should know:

- How to create a CORS issue
- How to verify CORS configuration
- How to troubleshoot CORS problems
- How CORS problems are affected by proxies
- How to simulate browser CORS interaction with an http client

## Bibliography

- [https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [https://en.wikipedia.org/wiki/Cross-origin_resource_sharing](https://en.wikipedia.org/wiki/Cross-origin_resource_sharing)
- [https://serverfault.com/questions/562756/how-to-remove-the-path-with-an-nginx-proxy-pass](https://serverfault.com/questions/562756/how-to-remove-the-path-with-an-nginx-proxy-pass)
- [https://enable-cors.org/server_nginx.html](https://enable-cors.org/server_nginx.html)
- [https://stackoverflow.com/questions/12173990/how-can-you-debug-a-cors-request-with-curl](https://stackoverflow.com/questions/12173990/how-can-you-debug-a-cors-request-with-curl)
- [https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials)
- [https://daniel.haxx.se/blog/2015/09/11/unnecessary-use-of-curl-x/](https://daniel.haxx.se/blog/2015/09/11/unnecessary-use-of-curl-x/)
- [https://gist.github.com/Stanback/7145487](https://gist.github.com/Stanback/7145487)
- [https://swagger.io/docs/open-source-tools/swagger-ui/usage/cors/](https://swagger.io/docs/open-source-tools/swagger-ui/usage/cors/)
- [https://ron-liu.medium.com/what-canonical-http-header-mean-in-golang-2e97f854316d](https://ron-liu.medium.com/what-canonical-http-header-mean-in-golang-2e97f854316d)
- [https://www.sneppets.com/software/connection-0-to-host-localhost-left-intact/](https://www.sneppets.com/software/connection-0-to-host-localhost-left-intact/)

## TLDR; Too Long Didn't read

CORS is a kind of protection against certain attacks when a browser attempts to get information from a domain/host that's different from the original request. CORS creates a default deny behavior and allows whitelisting of certain domains to enable functionality but also enhance security. In terms of design, understanding the appropriate CORS configuration and which domains you need to whitelist should be short and easy.

IF you find that your whitelist is too long or you have to be too permissive, it likely points that the scope of your application is too big or that the way your domains are architected is poor. The need for complicated CORS configuration could also indicate that your application as a whole isn't properly encapsulated. For example, CORS specifically has to do with browsers that have to hop between servers and maintain the headers (meaning that the destination server is aware that the request came from another entity they don't talk to directly); while there are some situations where this isn't avoidable, those are incredibly specific (e.g. you proxy a request to the token/authorize endpoint for an oauth2 server).

> You may also find that you have an API that must be reachable by an unknown number of [unknown] origins, and in that case, may be it makes sense to allow all origins, but limit the methods and headers that can be used

## Getting Started

To get started, execute the make run command; it should bring up the docker compose, build any dependencies it can't find and get everything up and running.

```sh
make run
```

You can verify that everything is up an running using the following command:

```sh
docker ps
```

```log
CONTAINER ID   IMAGE                                                   COMMAND                  CREATED         STATUS         PORTS                              NAMES
e59a9d046321   ghcr.io/antonio-alexander/go-blog-cors-nginx:latest     "/docker-entrypoint.…"   5 minutes ago   Up 5 minutes   0.0.0.0:8080->80/tcp               nginx
f3c79b107c66   ghcr.io/antonio-alexander/go-blog-cors-swagger:latest   "/docker-entrypoint.…"   5 minutes ago   Up 5 minutes   80/tcp, 0.0.0.0:8083->8080/tcp     swagger
e74d5d226748   ghcr.io/antonio-alexander/go-blog-cors:latest           "/bin/sh -c 'tar -xz…"   5 minutes ago   Up 5 minutes   2345/tcp, 0.0.0.0:8081->8080/tcp   go-blog-cors
feca60e012b8   ghcr.io/antonio-alexander/go-blog-cors:latest           "/bin/sh -c 'tar -xz…"   5 minutes ago   Up 5 minutes   2345/tcp, 0.0.0.0:8082->8080/tcp   go-blog-cors_proxy
Vidarr:go-blog-cors noobius$ 
```

The default configuration will be with CORS enabled but the least permissive (it only allows localhost); what you'll notice is that although you can access the swagger docs, attempts to try them out will fail with a CORS issue. Keep reading if you want to understand how to update the configuration to resolve this issue.

### Example(s)

There are two containers running the same code (example and example_proxy); each provides four endpoints:

- "/" (GET): this simply returns a string that says hello-world
- "/authroize (POST)": this is an authorize endpoint that authorizes a given username/password
- "/proxy" (GET): this proxies a request to example_proxy appending the path as a query parameter using the GET method
- "/proxy" (POST): this proxies a request to example_proxy appending the path as a query parameter using the POST method
- "/redirect" (GET): this redirects (301 temporary redirect) a request to example_proxy appending the path as a query parameter using the GET method
- "/redirect" (POST): this redirects (301 temporary redirect) a request to example_proxy appending the path as a query parameter using the POST method

Each of these endpoints can be used to validate the functionality of CORS and understand how to configure CORS with the least amount of permissions

### Configuration

It should go without saying that the configuration is the lifeblood of this entire example; below i'll provide an exhautive list of the configuration and some default values. The section(s) below should do a much better job of describing how these different configurations affect CORS (and your experience with the API).

The configuration below can be configured directly through the environment stanza within the docker compose or more elegantly through the .env file:

- CORS(cors): this configures which CORS implementation to use; valid values are cors and gorilla
- CORS_DEBUG(false): this configures whether to enable debugging for CORS
- ALLOW_CREDENTIALS: this automatically allows basic credentials via the Authorization Header
- ALLOWED_ORIGINS: the origins to allow (comma separated addresses with ports)
- ALLOWED_METHODS: the methods to allow
- ALLOWED_HEADERS: the headers to allow; must be comma separated and are automatically formatted to be canonical
- USERNAME: the username to use with basic auth
- PASSWORD: the password to use with basic auth
- CORS_ORIGIN: if cors simulation is enabled for the client, this is the origin that's injected into the OPTIONS request

Although we can use [gorilla/cors](https://github.com/gorilla/handlers/blob/master/cors.go), the examples will exclusively use [rs/cors](https://github.com/rs/cors) because the logging is so much better.

## Proxy vs Redirection

One of the more confusing things for me when doing this was understanding the functional difference between proxying a request and performing a redirect. Although the difference between the two is stark, as you pile on layers of complexity (e.g., CORS, NGINX, reverse proxies, DNS, etc...), it can be confusing.

When you proxy a request, the server takes whatever data you sent it (e.g., the header contents and request body) and sends it to another location known by the server, but most likely not by you. When the server proxies the request, they have full access to the request and can modify it as needed for better or worse.

> In a way, the ability to proxy is completely server-side and for the most part is bereft of any client-side protection, the server acts as a true man-in-the-middle and although it sends it to another server that information is obfuscated from the perspective of the client

Contrast proxied requests with redirection (e.g. 301 permanent redirect or 302 temporary redirect); although there is a difference between the two when it comes to SEO (search engine optimization), for purposes of this repository, we can simplify it to mean that it tells the client __where__ to forward its request to. Redirection _involves_ the client and gives it the opportunity to protect itself; in this, the destination isn't obfuscated and the client has the option to NOT follow the redirect.

When a redirect is done in a browser, CORS is triggered, while for a proxied request, CORS is ignored; this is more because there's no implicit OPTIONS query sent rather than CORS not "running".

> CORS is enforced by browsers to mitigate cross-site-scripting and that doesn't cover proxied requests, only interactions with the browser and by extension http redirects

This is somewhat easy to confirm using the example server. I'll do the same request, one using curl to hit example_proxy through example with a proxy request and the other using redirection. I'l be communicating with the servers directly rather than through NGINX:

```sh
curl --verbose \
  'http://localhost:8081/proxy' \
  -H 'accept: application/json'
```

```log
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> GET /proxy HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.88.1
> accept: application/json
> 
< HTTP/1.1 200 OK
< Content-Length: 27
< Content-Type: application/json; charset=utf-8
< Date: Wed, 26 Nov 2025 22:30:41 GMT
< Vary: Origin
< Vary: Origin
< 
* Connection #0 to host localhost left intact
{"message":"Hello, World!"}
```

Note that the log above only shows localhost:8081 when the request is proxied to example_proxy:8080 (from the perspective of the container). Compare with doing the same using redirection:

```sh
curl --verbose \
  'http://localhost:8081/redirect' \
  -H 'accept: application/json'
```

```log
*   Trying 127.0.0.1:8081...
* Connected to localhost (127.0.0.1) port 8081 (#0)
> GET /redirect HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.88.1
> accept: application/json
> 
< HTTP/1.1 302 Found
< Content-Type: text/html; charset=utf-8
< Location: http://localhost:8080
< Vary: Origin
< Date: Wed, 26 Nov 2025 22:34:14 GMT
< Content-Length: 44
< 
<a href="http://localhost:8080">Found</a>.

* Connection #0 to host localhost left intact
```

You'll see in this log, that after you connection to localhost:8081, it redirects you to localhost:8080; in this case you (the client) are made aware of where your request went.

The second example is closer to what happens if you do this from swagger/browser, you'll notice different behavior on the redirection endpoint (depending on your CORS configuration) which will generally be unaffected by your proxy configuration (assuming you don't configure the proxy to point to a destination that's unreachable by the container).

## CORS (Cross Origin Resource Sharing)

This purpose of this section is to provide you the ability to actively create specific CORS issues and then resolve them. I'll be providing examples that are both using [curl](https://curl.se/) and [swagger](http://localhost:8080/swagger). Throughout these examples, we'll be modifying the values in the [.env](.env) file.

> You can bring the compose up/down with make run, ctrl+c and then make run, alternatively, you can forceibly remove the go-blog-cors container and execute make run again to re-create the container (with the new configuration)

To experience our first CORS issue, we'll need to set the configuration appropriately; modify the .env file such that:

- CORS_DEBUG="true"
- CORS="cors"
- ALLOWED_ORIGINS="`http://example`"
- ALLOWED_METHODS="GET,OPTIONS"
- ALLOWED_HEADERS=""
- PROXY_CORS_ALLOWED_ORIGINS="https://example_proxy"
- PROXY_CORS_ALLOWED_METHODS="GET,OPTIONS"
- PROXY_CORS_ALLOWED_HEADERS=""

Once the file is updated and saved, bring the compose up (make run).

> Prior to getting into the nitty gritty, I want to to make sure that you're keeping track of _perspective_ when it comes to interacting with swagger and using curl. Perspective makes or breaks the entire idea/paradigm of CORS. I'll do my best to note the perspective given the situation

The __MOST__ confusing thing about CORS is that it's something that's _specifically_ implemented in browsers to mitigate some request's ability to interact with web servers that exist on another server. It's not a replacement for [server-side access control](https://github.com/rs/cors/issues/129#issuecomment-1264434797) and even some requests made by browsers that are simple will make it through. As a result, some of these examples are inherently flawed as they may not be 1:1 with your  implementation but always keep in mind that CORS is not a 100% mitigation strategy.

> You _CAN_ check the headers to confirm that they're present within your own endpoints (in _addition_ to the CORS middleware) to perform a 100% mitigation strategy, just keep in mind that if you have non-browser clients, they may not provide the origin. Additionally, you may not have routes for certain methods and without CORS, you'd get a 405 instead of a 403; this could be an entry way for DoS or a kind of security leak

Within a browser, for most requests, an OPTIONS request is initially sent to confirm that the origin, headers, and method of the request are allowed, this OPTIONS request will return one or more headers communicating if it's allowed or not; if any those headers __aren't__ received, the browser will return a CORS error rather than allowing the request to go through (you'll notice that curl doesn't care).

This is a failing request, note that the origin is `https://example` which doesn't match the configuration above:

```sh
curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: GET" \
  -H "Origin: https://example" \
  http://localhost:8080
```

This is the response from the above command, although you get a 204 back, you'll also notice that there are no headers returned (this is easier when you see a successful OPTIONS request).

```log
HTTP/1.1 204 No Content
Server: nginx/1.21.6
Date: Fri, 25 Nov 2022 21:03:45 GMT
Connection: keep-alive
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
```

This is the response from the CORS logs (docker logs example -f); this will at the very least communicate that the origin isn't allowed and cause a browser to stop (but not curl)

```log
[cors] 2022/11/25 21:03:45 Handler: Preflight request
[cors] 2022/11/25 21:03:45   Preflight aborted: origin 'https://example' not allowed
```

Alternatively, the commands and responses below, show successful responses

```sh
curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: GET" \
  -H "Origin: http://example" \
  http://localhost:8080
```

Note that there are headers returned for Access-Control-Allow-Methods and Access-Control-Allow-Origin

```log
HTTP/1.1 204 No Content
Server: nginx/1.21.6
Date: Fri, 25 Nov 2022 21:06:03 GMT
Connection: keep-alive
Access-Control-Allow-Methods: GET
Access-Control-Allow-Origin: http://example
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
```

Also note that headers are provided below to indicate that the origins are OK

```log
[cors] 2022/11/25 21:06:03 Handler: Preflight request
[cors] 2022/11/25 21:06:03   Preflight response headers: map[Access-Control-Allow-Methods:[GET] Access-Control-Allow-Origin:[http://example] Vary:[Origin Access-Control-Request-Method Access-Control-Request-Headers]]
```

With the given configuration, if we attempt to use an allowed origin, but a method that's not allowed, we'd get the following feedback:

```sh
curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: POST" \
  -H "Origin: http://example" \
  http://localhost:8080
```

```log
[cors] 2022/11/25 23:40:41 Handler: Preflight request
[cors] 2022/11/25 23:40:41   Preflight aborted: method 'POST' not allowed
```

```log
HTTP/1.1 204 No Content
Server: nginx/1.21.6
Date: Fri, 25 Nov 2022 23:40:41 GMT
Connection: keep-alive
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
```

We can fix this problem by modifying the configuration to the following (remember to make stop/make run):

- ALLOWED_METHODS="GET,POST,OPTIONS"

Now when we attempt the options call, it'll be successful:

```sh
curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: POST" \
  -H "Origin: http://example" \
  http://localhost:8080
```

```log
HTTP/1.1 204 No Content
Server: nginx/1.21.6
Date: Sat, 26 Nov 2022 00:17:36 GMT
Connection: keep-alive
Access-Control-Allow-Methods: POST
Access-Control-Allow-Origin: http://example
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
```

```log
[cors] 2022/11/26 00:17:35   Preflight response headers: map[Access-Control-Allow-Methods:[POST] Access-Control-Allow-Origin:[http://example] Vary:[Origin Access-Control-Request-Method Access-Control-Request-Headers]]
```

### CORS: Authorization Requests

Along with your ability to restrict the origin, methods and headers, you can also specifically restrict whether or not authorization is allowed. Authorization works (generally) by including an authorization header, cookie or other method (TLS cert maybe?). Again it won't "prevent" the request unless you __add__ code, but the incomming request will lack the values in the Access-Control-Request-Headers and AllowCredentials Header.

To show the negative case for authorization, configure the example with the following:

- ALLOW_CREDENTIALS="false"
- ALLOWED_HEADERS=""

```sh
curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Authorization" \
  -H "X-Requested-With: POST" \
  -H "Origin: http://example" \
  http://localhost:8080
```

You'll notice that the ouput doesn't include anything being allowed, and the logs from CORS indicate that the authorization header specifically isn't allowed

```log
HTTP/1.1 204 No Content
Server: nginx/1.21.6
Date: Sat, 26 Nov 2022 07:28:10 GMT
Connection: keep-alive
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
```

```log
[cors] 2022/11/26 07:28:10 Handler: Preflight request
[cors] 2022/11/26 07:28:10   Preflight aborted: headers '[Authorization]' not allowed
```

Let's modify the configuration and try again:

- ALLOW_CREDENTIALS="false"
- ALLOWED_HEADERS="Authorization"

```sh
curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Authorization" \
  -H "X-Requested-With: POST" \
  -H "Origin: http://example" \
  http://localhost:8080
```

Again, you'll see that this was "successful", but there also isn't a mention of Access-Controll-Allow-Credentials, a browser may cause this to fail by REQUIRING that that header is returned from the OPTIONS call.

```log
HTTP/1.1 204 No Content
Server: nginx/1.21.6
Date: Sat, 26 Nov 2022 07:38:50 GMT
Connection: keep-alive
Access-Control-Allow-Headers: Authorization
Access-Control-Allow-Methods: POST
Access-Control-Allow-Origin: http://example
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
```

```log
[cors] 2022/11/26 06:55:53 Handler: Preflight request
[cors] 2022/11/26 06:55:53   Preflight response headers: map[Access-Control-Allow-Credentials:[true] Access-Control-Allow-Methods:[POST] Access-Control-Allow-Origin:[http://example] Vary:[Origin Access-Control-Request-Method Access-Control-Request-Headers]]
```

Let's update the configuration one more time to get a completely successful response:

- ALLOW_CREDENTIALS="true"
- ALLOWED_HEADERS="Authorization"

```sh
curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Authorization" \
  -H "X-Requested-With: POST" \
  -H "Origin: http://example" \
  http://localhost:8080
```

The only real difference with this request is that now there's an Access-Control-Allow-Credentials header that's set to true.

```log
HTTP/1.1 204 No Content
Server: nginx/1.21.6
Date: Sat, 26 Nov 2022 07:44:05 GMT
Connection: keep-alive
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Authorization
Access-Control-Allow-Methods: POST
Access-Control-Allow-Origin: http://example
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
```

```log
[cors] 2022/11/26 07:44:05 Handler: Preflight request
[cors] 2022/11/26 07:44:05   Preflight response headers: map[Access-Control-Allow-Credentials:[true] Access-Control-Allow-Headers:[Authorization] Access-Control-Allow-Methods:[POST] Access-Control-Allow-Origin:[http://example] Vary:[Origin Access-Control-Request-Method Access-Control-Request-Headers]]
```

Keep in mind that without any additional code/middle-ware, even though the CORS preflight checks failed, the requests could still be completed without issue.

### CORS: Swagger

Swagger provides an interesting data point in that it's a simple way to demo how CORS works (specifically the how Chrome handles CORS). I've used the [swagger-ui](https://swagger.io/tools/swagger-ui/) image and injected two swagger files that expose the API for the example on port 8080 and 8083 (dynamic host for swagger is a pain). This allows us to expose swagger via [http://localhost:8083/swagger](http://localhost:8083/swagger); it also does so through [NGINX](#cors-nginx).

[http://localhost:8083](http://localhost:8083) connects directly to the swagger container (rather than being proxied through NGINX); this has a DIRECT effect on CORS, let's start with this configuration:

- ALLOWED_ORIGINS="`http://example_proxy,http://example`"
- ALLOWED_METHODS="GET,OPTIONS"

Let's connect to swagger on [http://localhost:8083](http://localhost:8083), use the go-blog-cors(example) definitio, use the  go to View > Developer > Inspect Element and then navigate to the Network tab. Scroll down to the cors / GET endpoint. Clearing the items prior to executing the swagger endpoint will make it easier to see success/failure.

Click "Try it out" and then "Execute"; this is the log from example:

```logs
[cors] 2022/11/28 04:04:48 Handler: Actual request
[cors] 2022/11/28 04:04:48   Actual request no headers added: origin 'http://localhost:8083' not allowed
```

This should be pretty familiar, if you look at the request headers, you'll see that swagger automatically adds an Origin header of `http://localhost:8083` and since that address isn't in the ALLOWED_ORIGINS configuration, we get the CORS error. If we modify the configuration as such:

- ALLOWED_ORIGINS="`http://example_proxy,http://example,http://localhost:8083`"
- PROXY_ALLOWED_ORIGINS="`http://example_proxy,http://example,http://localhost:8083`"

And try again; the request should be successful.

```log
[cors] 2022/11/28 04:06:01   Actual response added headers: map[Access-Control-Allow-Credentials:[true] Access-Control-Allow-Origin:[http://localhost:8083] Vary:[Origin]]
```

Feel free to try some of the other examples from this perspective of swagger; also note that this example has zero interaction with NGINX.

### CORS: Redirection

CORS is engaged when making a request that involves redirection (e.g., going from one address to another); this is especially true with the forwarding address exists on a different host (e.g. from google.com to yahoo.com) and when the server may _update_ your request with additional headers (relatively common if not overly useful).

As described in the section [Proxy vs Redirection](#proxy-vs-redirection), redirection involves the client and triggers CORS while proxying doesn't involve the client and is largely opaque from the client's perspective.

Taking into account CORS, if you send a request and it's forwarded three times, it means that CORS is engaged three times for each of those servers (via the implicit OPTIONS request sent by the client). We can test this out using these configuration options:

```sh
# Common/WebServer Configuration
CORS="cors"
CORS_DEBUG="true"
USERNAME="username"
PASSWORD="password"

# redirect/proxy configuration
PROXY_ADDRESS="example_proxy"
PROXY_PORT="8080"
REDIRECT_ADDRESS="localhost" 
REDIRECT_PORT="8080"

# CORS Configuration for example
ALLOW_CREDENTIALS="true"
ALLOWED_ORIGINS="http://localhost:8080"
ALLOWED_METHODS="POST,GET,OPTIONS"
ALLOWED_HEADERS="Authorization"

# CORS configuration for example_proxy
PROXY_ALLOW_CREDENTIALS="true"
PROXY_ALLOWED_ORIGINS="http://localhost:8080"
PROXY_ALLOWED_METHODS="POST,GET,OPTIONS"
PROXY_ALLOWED_HEADERS=""

# NGINX Configuration
# use this when you want nginx to replace the response headers (bypass CORS)
NGINX_CONFIG_FILE="./config/config_origin.conf"
# use this when you want nginx to inject the origin/preserve headers 
# NGINX_CONFIG_FILE="./config/config_bypass.conf"
```

This should configure the example service to allow the AUTHORIZATION header while the proxy service won't allow it. With that configuration, on the assumption that CORS is triggered both for the initial call to the example service and again after the call is redirected, if we attempt to access a redirected endpoint via swagger, the initial OPTIONS and GET call should succeed, but after the REDIRECT, the second OPTIONS call should fail.

### CORS: NGINX

Although NGINX isn't the focus of this repository/example; NGINX is a common use case when trying to put different microservice APIs together (e.g. a [facade pattern](https://docs.firstdecode.com/microservices-architecture-style/design-patterns-for-microservices/facade-pattern/)) and CORS definitely comes into play when trying to integate a front end. NGINX can be a proxy to have a single address/port to communicate with other microservices with the path being different. In terms of CORS, NGINX can have two configurations:

  1. Bypass: configure NGINX to replace the response headers to ensure the most permissive CORS configuration; this is done without the knowledge of underlying services and although is a simple solution for development; this should almost NEVER be a practice in production
  2. Passthrough: configure NGINX to pass through the provided headers and _REPLACE_ the Origin header with a static value (even if it's set for some reason)

I'll demonstrate both solutions using swagger (through NGINX). To start, use the following configuration:

- ALLOWED_ORIGINS="`http://example_proxy,http://example`"
- ALLOWED_METHODS="GET,POST,OPTIONS"
- NGINX_CONFIG_FILE="./cmd/nginx/config_origin.conf"

This configuration will allow us to test CORS for swagger by passing through the Origin header (in this case we're actually using the CORS in the example microservice). To start, connect to swagger using [http://localhost:8080/swagger](http://localhost:8080/swagger) and use the go-blog-cors(nginx) definition. Attempt to execute the cors / api.

This should fail due to a CORS error and the example microservice should have the following output:

```log
[cors] 2022/11/28 04:43:55 Handler: Actual request
[cors] 2022/11/28 04:43:55   Actual request no headers added: origin 'http://localhost:8080' not allowed
```

Again, we can _fix_ this problem by updating the configuration to include "`http://localhost:8080`" in the ALLOWED_ORIGINS:

- ALLOWED_ORIGINS="`http://example_proxy,http://example,http://localhost:8080`"
- ALLOWED_METHODS="GET,POST,OPTIONS"
- NGINX_CONFIG_FILE="./cmd/nginx/config_origin.conf"

Alternatively, you may want to bypass CORS in the backend microservice altogether by using the alternate configuration below:

- ALLOWED_ORIGINS="`http://example_proxy,http://example`"
- ALLOWED_METHODS="GET,POST,OPTIONS"
- NGINX_CONFIG_FILE="./cmd/nginx/config_bypass.conf"

Keeping in mind that this configuration actually _BYPASSES_ the CORS configuration in the example microservice, let's start by
connecting to swagger using [http://localhost:8080/swagger](http://localhost:8080/swagger) and use the go-blog-cors(nginx) definition. Attempt to execute the cors / api.

You'll note that although the endpoint is successful, the logs for the example microservice still show:

```log
[cors] 2022/11/28 04:56:04 Handler: Actual request
[cors] 2022/11/28 04:56:04   Actual request no headers added: missing origin
```

Even though the CORS failed, NGINX will overwrite the response headers to include the expected CORS headers so that the browser doesn't experience any CORS issues.
