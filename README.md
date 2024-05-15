# OAuth2 Proxy Service

A simple webservice that proxies requests to HTTP endpoints that are protected by OAuth2 authorization. I.e. this service handles the OAuth2 part (currently supporting only `client_credentials` with `default` scope.

## How to use

Just fill out `proxy.ini` (see below) and run (either compile yourself or run with Docker). If you change the bind port and want to use the Docker image, do not forget to update the Dockerfile to expose correct port.

## Configuration

Configuration is stored in `proxy.ini`. It has the following format:

```
[Webservice]
Bind = :8080                                              // the service will listen for requests on this address

[Proxy]
Server = https://gateway.example.com/secured/v1           // the target server

[OAuth2]
TokenEndpoint = https://gateway.example.com/oauth2/token  // OAuth2 token endpoint
ClientId = qwxyz123                                       // your OAuth2 client id
ClientSecret = yxcvb123                                   // your OAuth2 client secret
```