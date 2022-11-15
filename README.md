# go-blog-cors (github.com/antonio-alexander/go-blog-cors)

This is a companion repository for an article describing CORS (cross origin resource sharing) with an emphasis on Golang. CORS is a protection enforced by browsers to ensure that calls between browsers and server(s) (especially between web servers) are explicitly allowed. In addition to this document itself, there will also be example source code to demonstrate troubleshooting and development. After review of this repository, you should know:

- How to create a CORS issue
- How to verify CORS configuration
- How to troubleshoot CORS problems
- How CORS problems are affected by proxies

## TLDR; Too Long Didn't read

CORS is a kind of protection against certain attacks when a browser attempts to get information from a domain/host that's different from the original request. CORS creates a default deny behavior and allows whitelisting of certain domains to enable functionality but also enhance security. In terms of design, understanding the appropriate CORS configuration and which domains you need to whitelist should be short and easy.

IF you find that your whitelist is too long or you have to be too permissive, it likely points that the scope of your application is too big or that the way your domains are architected is poor. The need for complicated CORS configuration could also indicate that your application as a whole isn't properly encapsulated. For example, CORS specifically has to do with browsers that have to hop between servers and maintain the headers (meaning that the destination server is aware that the request came from another entity they don't talk to directly); while there are some situations where this isn't avoidable, those are incredibly specific (e.g. you proxy a request to the token/authorize endpoint for an oauth2 server).

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

Each of these endpoints can be used to validate the functionality of CORS and understand how to configure CORS with the least amount of permissions

### Configuration

It should go without saying that the configuration is the lifeblood of this entire example; below i'll provide an exhautive list of the configuration and some default values. The section(s) below should do a much better job of describing how these different configurations affect CORS (and your experience with the API).

The configuration below can be configured directly through the environment stanza within the docker compose or more elegantly through the .env file:

- CORS():
- CORS_DEBUG():
- ALLOW_CREDENTIALS():
- ALLOWED_ORIGINS():
- ALLOWED_METHODS():
- ALLOWED_HEADERS():
- USERNAME():
- PASSWORD():

Although we can use [gorilla/cors](https://github.com/gorilla/handlers/blob/master/cors.go), the examples will exclusively use [rs/cors](https://github.com/rs/cors) because the logging is so much better.

## CORS (Cross Origin Resource Sharing)

This purpose of this section is to provide you the ability to actively create specific CORS issues and then resolve them. I'll be providing examples that are both using [curl](https://curl.se/) and [swagger](http://localhost:8080/swagger). Throughout these examples, we'll be modifying the values in the [.env](.env) file.

> You can bring the compose up/down with make run, ctrl+c and then make run, alternatively, you can forceibly remove the go-blog-cors container and execute make run again to re-create the container (with the new configuration)

To experience our first CORS issue, we'll need to set the configuration appropriately; modify the .env file such that:

- CORS_DEBUG="true"
- CORS="cors"
- ALLOWED_ORIGINS="http://example"
- ALLOWED_METHODS="GET,OPTIONS"
- ALLOWED_HEADERS=""
- PROXY_CORS_ALLOWED_ORIGINS="https://example_proxy"
- PROXY_CORS_ALLOWED_METHODS="GET,OPTIONS"
- PROXY_CORS_ALLOWED_HEADERS=""

Once the file is updated and saved, bring the compose up (make run).

> Prior to getting into the nitty gritty, I want to to make sure that you're keeping track of _perspective_ when it comes to interacting with swagger and using curl. Perspective makes or breaks the entire idea/paradigm of CORS. I'll do my best to note the perspective given the situation

The __MOST__ confusing thing about CORS is that it's something that's _specifically_ implemented in browsers to mitigate some request's ability to interact with web servers that exist on another server. It's not a replacement for [server-side access control](https://github.com/rs/cors/issues/129#issuecomment-1264434797) and even some requests made by browsers that are simple will make it through. As a result, some of these examples are inherently flawed as they may not be 1:1 with your eventual implementation but always keep in mind that CORS is not a 100% mitigation strategy.

> You _CAN_ check the headers to confirm that they're present within your own endpoints to perform a 100% mitigation strategy, just keep in mind that if you have non-browser clients, they may not provide the origin

Within a browser, for most requests, an OPTIONS request is initially sent to confirm that the Origin and Method of the request are allowed, this OPTIONS request will return one or more headers communicating if it's allowed or not; if those headers __aren't__ received, the browser will return a CORS error rather than allowing the request to go through (you'll notice that curl doesn't care).

This is a failing request, note that the origin is <https://example> which doesn't match the configuration above:

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

Keep in mind that without any additional code/middle-ware, even though the CORS preflight checks failed, the requests would still be completed without issue.

### CORS: Proxy Requests

This example also provides the ability to proxy a request from one web server to another (they're instances of the same application); the .env file covers configuration from both instances of the example (example and example_proxy). Although this is a small part of a much larger conversation regarding architecture, a proxy request is unique because it _forwards_ the request rather than encapsulating the request; from the perspective of the server being proxied to, the request comes __from__ the user rather than the server proxying the request.

> If you can extrapolate, there's a scalability issue with CORS (by design I think) such that you'd have to coordinate the CORS configuration if you had to proxy through multiple web servers and __HAD__ to worry about CORS/browsers. One situation where proxying is necessary/common would be the Authorize/token endpoints for OAuth2

Lets start with the following configuration:

- PROXY_ADDRESS="example_foobar"
- PROXY_PORT="8080"
- CORS="cors"
- ALLOW_CREDENTIALS="true"
- ALLOWED_ORIGINS="http://example"
- ALLOWED_METHODS="POST,GET,OPTIONS"
- ALLOWED_HEADERS="Authorization"
- PROXY_ALLOW_CREDENTIALS="true"
- PROXY_ALLOWED_ORIGINS="https://example_proxy"
- PROXY_ALLOWED_METHODS="POST,GET,OPTIONS"
- PROXY_ALLOWED_HEADERS=""

We won't go through the trouble of re-creating all of the above situations (I'll leave that to you), but I do want to show how proxying a request has a couple of differences, below we'll send both requests, one for the OPTIONS call (that would come from the browser) and another for making the actual request.

```sh
curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: GET" \
  -H "Origin: http://example" \
  http://localhost:8080/proxy

curl --include -X OPTIONS \
  -H "Access-Control-Request-Method: GET" \
  -H "Origin: http://example" \
  http://localhost:8080/proxy
```

These headers returned make sense and are expected (from the perspective of example)

```log
HTTP/1.1 204 No Content
Server: nginx/1.21.6
Date: Sat, 26 Nov 2022 19:02:55 GMT
Connection: keep-alive
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: GET
Access-Control-Allow-Origin: http://example
Vary: Origin
Vary: Access-Control-Request-Method
Vary: Access-Control-Request-Headers
```

Here, we try to execute the actual request (and we're expecting a 200 OK response):

```sh
curl --include -X 'GET' \
  -H "Origin: http://example_proxy" \
  -H 'accept: application/json' \
  http://localhost:8080/proxy
```

We instead get a 502 bad gateway (because the proxy address/port doesn't exist).

```log
HTTP/1.1 502 Bad Gateway
Server: nginx/1.21.6
Date: Sat, 26 Nov 2022 19:03:12 GMT
Content-Length: 0
Connection: keep-alive
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: http://example
Vary: Origin
```

<!-- REVIEW: this doesn't seem to be true, or to be half true -->
<!-- This tells us that the OPTIONS request is NOT proxied to the destination (I tried this); it's intercepted by the CORS middleware every time, so the CORS configuration being used is specifically the initial entrypoint and NOT the proxied entrypoint. -->

Let's update the configuration to "fix" this problem:

- PROXY_ADDRESS="example_proxy"
- PROXY_PORT="8080"
- CORS="cors"
- ALLOWED_ORIGINS="http://example"
- ALLOWED_METHODS="POST,GET,OPTIONS"
- PROXY_ALLOWED_ORIGINS="https://example_proxy"
- PROXY_ALLOWED_METHODS="POST,GET,OPTIONS"

```sh
curl --include -X 'GET' \
  -H "Origin: http://example_proxy" \
  -H 'accept: application/json' \
  http://localhost:8080/proxy
```

```log
HTTP/1.1 200 OK
Server: nginx/1.21.6
Date: Sat, 26 Nov 2022 23:51:12 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 27
Connection: keep-alive
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: http://example_proxy
Vary: Origin
Vary: Origin

{"message":"Hello, World!"}
```

This is the log from example

```log
[cors] 2022/11/26 23:51:12 Handler: Actual request
[cors] 2022/11/26 23:51:12   Actual request no headers added: origin 'http://example_proxy' not allowed
attempting to proxy to: http://example_proxy:8080
```

This is the log from example_proxy

```log
[cors] 2022/11/26 23:51:12 Handler: Actual request
[cors] 2022/11/26 23:51:12   Actual response added headers: map[Access-Control-Allow-Credentials:[true] Access-Control-Allow-Origin:[http://example_proxy] Vary:[Origin]]
```

Again, note that although the initial request (to example) fails the CORS test, the request makes it to example_proxy and the expected JSON string is returned; you can "fix" this problem by updating the configuration as follows:

- PROXY_ADDRESS="example_proxy"
- PROXY_PORT="8080"
- CORS="cors"
- ALLOWED_ORIGINS="http://example"
- ALLOWED_METHODS="POST,GET,OPTIONS"
- PROXY_ALLOWED_ORIGINS="https://example_proxy,http://example"
- PROXY_ALLOWED_METHODS="POST,GET,OPTIONS"

```sh
curl --include -X 'GET' \
  -H "Origin: http://example_proxy" \
  -H 'accept: application/json' \
  http://localhost:8080/proxy
```

Although you've seen this plenty of times, something interesting to note is that there are __TWO__ sets of headers for all of the CORS headers (Access-Control-Allow-Credentials, Access-Control-Allow-Origin, Vary).

```log
HTTP/1.1 200 OK
Server: nginx/1.21.6
Date: Sun, 27 Nov 2022 01:54:17 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 27
Connection: keep-alive
Access-Control-Allow-Credentials: true
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: http://example_proxy
Access-Control-Allow-Origin: http://example_proxy
Vary: Origin
Vary: Origin

{"message":"Hello, World!"}
```

```log
[cors] 2022/11/27 01:52:32 Handler: Actual request
[cors] 2022/11/27 01:52:32   Actual response added headers: map[Access-Control-Allow-Credentials:[true] Access-Control-Allow-Origin:[http://example_proxy] Vary:[Origin]]
attempting to proxy to: http://example_proxy:8080
```

```log
[cors] 2022/11/27 01:52:32 Handler: Actual request
[cors] 2022/11/27 01:52:32   Actual response added headers: map[Access-Control-Allow-Credentials:[true] Access-Control-Allow-Origin:[http://example_proxy] Vary:[Origin]]
```

### CORS: Swagger

Swagger provides an interesting data point in that it's a simple way to demo how CORS works (specifically the how Chrome handles CORS). I've used the [swagger-ui](https://swagger.io/tools/swagger-ui/) image and injected two swagger files that expose the API for the example on port 8080 and 8083 (dynamic host for swagger is a pain). This allows us to expose swagger via [http://localhost:8083/swagger](http://localhost:8083/swagger); it also does so through [NGINX](#cors-nginx).

[http://localhost:8083](http://localhost:8083) connects directly to the swagger container (rather than being proxied through NGINX); this has a DIRECT effect on CORS, let's start with this configuration:

- ALLOWED_ORIGINS="http://example_proxy,http://example"
- ALLOWED_METHODS="GET,OPTIONS"

Let's connect to swagger on [http://localhost:8083](http://localhost:8083), use the go-blog-cors(example) definitio, use the  go to View > Developer > Inspect Element and then navigate to the Network tab. Scroll down to the cors / GET endpoint. Clearing the items prior to executing the swagger endpoint will make it easier to see success/failure.

Click "Try it out" and then "Execute"; this is the log from example:

```logs
[cors] 2022/11/28 04:04:48 Handler: Actual request
[cors] 2022/11/28 04:04:48   Actual request no headers added: origin 'http://localhost:8083' not allowed
```

This should be pretty familiar, if you look at the request headers, you'll see that swagger automatically adds an Origin header of <http://localhost:8083> and since that address isn't in the ALLOWED_ORIGINS configuration, we get the CORS error. If we modify the configuration as such:

- ALLOWED_ORIGINS="http://example_proxy,http://example,http://localhost:8083"
- PROXY_ALLOWED_ORIGINS="http://example_proxy,http://example,http://localhost:8083"

And try again; the request should be successful.

```log
[cors] 2022/11/28 04:06:01   Actual response added headers: map[Access-Control-Allow-Credentials:[true] Access-Control-Allow-Origin:[http://localhost:8083] Vary:[Origin]]
```

Feel free to try some of the other examples from this perspective of swagger; also note that this example has zero interaction with NGINX.

### CORS: NGINX

Although NGINX isn't the focus of this repository/example; NGINX is a common use case when trying to put different microservice APIs together (e.g. a [facade pattern](https://docs.firstdecode.com/microservices-architecture-style/design-patterns-for-microservices/facade-pattern/)) and CORS definitely comes into play when trying to integate a front end. NGINX can be a proxy to have a single address/port to communicate with other microservices with the path being different. In terms of CORS, NGINX can have two configurations:

  1. Bypass: configure NGINX to replace the response headers to ensure the most permissive CORS configuration; this is done without the knowledge of underlying services and although is a simple solution for development; this should almost NEVER be a practice in production
  2. Passthrough: configure NGINX to pass through the provided headers and _REPLACE_ the Origin header with a static value (even if it's set for some reason)

I'll demonstrate both solutions using swagger (through NGINX). To start, use the following configuration:

- ALLOWED_ORIGINS="http://example_proxy,http://example"
- ALLOWED_METHODS="GET,POST,OPTIONS"
- NGINX_CONFIG_FILE="./cmd/nginx/config_origin.conf"

This configuration will allow us to test CORS for swagger by passing through the Origin header (in this case we're actually using the CORS in the example microservice). To start, connect to swagger using [http://localhost:8080/swagger](http://localhost:8080/swagger) and use the go-blog-cors(nginx) definition. Attempt to execute the cors / api.

This should fail due to a CORS error and the example microservice should have the following output:

```log
[cors] 2022/11/28 04:43:55 Handler: Actual request
[cors] 2022/11/28 04:43:55   Actual request no headers added: origin 'http://localhost:8080' not allowed
```

Again, we can _fix_ this problem by updating the configuration to include "http://localhost:8080" in the ALLOWED_ORIGINS:

- ALLOWED_ORIGINS="http://example_proxy,http://example,http://localhost:8080"
- ALLOWED_METHODS="GET,POST,OPTIONS"
- NGINX_CONFIG_FILE="./cmd/nginx/config_origin.conf"

Alternatively, you may want to bypass CORS in the backend microservice altogether by using the alternate configuration below:

- ALLOWED_ORIGINS="http://example_proxy,http://example"
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
