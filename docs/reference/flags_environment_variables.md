# Flags and environment variables

Hoverfly can be configured using flags on startup, or using environment variables.


## Authentication

Hoverfly uses a combination of basic auth and JWT (JSON Web Tokens) to authenticate users. [Read more about authentication](#).

    -no-auth

Disable authentication. Currently it is enabled by default. If you disable authentication, you can use any username and password combination to authenticate.

Authentication can also be disabled using the `HoverflyAuthDisabled` environment variable. For example: `HoverflyAuthDisabled=true`.

Note: if Hoverfly is set to use in-memory persistence only (see `-db` flag below), AdminUI authentication will be disabled.

    -add

Add a new user. 

    -username <string>

Username for new user.

    -password <string>
      	
Password for new user.
    	
    -admin <string>
    	
Supply '-admin false' to make this a non-admin user (defaults to 'true').
    	
Example:

    ./hoverfly -add -username hfadmin -password hfpass -admin false   	

This creates a new non-admin user with the username 'hfadmin' and the password 'hfpass'.

By default, Hoverfly will generate a random secret key every time it is started.

You can specify a secret key using the `HoverflySecret` environment variable. For example: `HoverflySecret=<my_secret>`. 

By default, token expiration is set to 1 day. You can specify the token expiration (in seconds) using the `HoverflyTokenExpiration` environment variable. For example: `HoverflyTokenExpiration=3600`.

## Port selection

    -ap <string>

Sets the Admin UI port. 
    
    -pp <string>
         
Sets the proxy port.

Example:
         
    ./hoverfly -ap 1234 -pp 4567         

This starts Hoverfly with the Admin UI on port 1234 and the proxy on 5678. 
   	
Ports can also be set with the `AdminPort` and `ProxyPort` environment variables.   	
    	
## Mode selection, import & middleware

By default, Hoverfly starts in "virtualize" mode.

    -capture
    	
Start Hoverfly in capture mode - transparently intercepts and saves requests/response.
   	
    -synthesize

Start Hoverfly in synthesize mode. Middleware is required to generate responses for incoming requests.
       	
    -modify
      	
Start Hoverfly in modify mode. Middleware is required. Middleware is applied to both outgoing and incoming HTTP traffic.

    -middleware <string>

Set proxy to use middleware. Supply the path to the middleware script.

Example:

    ./hoverfly -synthesize -middleware "scripts/gen_response.py"

This starts Hoverfly in synthesize mode with a middleware script that generates responses on the fly.

Middleware can also be specified using the `HoverflyMiddleware` environment variable. For example: `HoverflyMiddleware="scripts/gen_response.py"`.

 
    -import <string>

Import a virtual service JSON file from file system or URL. For example:
     
    ./hoverfly -import http://mypage.com/service_x.json
     
    ./hoverfly -import path/to/my/service_x.json      

This starts Hoverfly in virtualize mode (default mode), and imports a virtual service JSON file from either a URL or the local filesystem.
      	    	
Note: currently, service JSON can only be exported from Hoverfly via the API.  

## Destination URL
    	    	
    -dest <string>
    
Specify which hosts to process. For example: 

    ./hoverfly -dest fooservice.org -dest barservice.org -dest catservice.org

This will start Hoverfly in virtualize mode, and only virtualize requests that are sent to fooservice.org, barservice.org and catservice.org. Requests to all other hosts will pass through.
   
    -destination <string>
    
Specify which URI to catch using regluar expression. (Defaults to ".").
 	    	
## Persistence
 	    	
    -db <string>
    	
Persistent storage to use. By default, Hoverfly uses BoltDB to store data in a file on disk. Specify 'memory' to disable this and use in-memory persistence only. 

Note: If 'memory' is specified, AdminUI authentication will be disabled.

    -db-dir <string>

Path to BoltDB data file. By default, a "requests.db" file will be created in the Hoverfly directory. Supply a custom path and/or filename to use a different file or location. The file will be created if it doesn't exist.	
    	
The database file/path can also be set using the `HoverflyDB` environment variable.   
      	    	
## TLS & Certificate management

    -generate-ca-cert

Hoverfly will generate private and public keys in the current directory.

    -cert-name <string>

Certificate name (defaults to "hoverfly.proxy")

    -cert-org <string>

Organization name for new certificate (defaults to "Hoverfly Authority"). For example:

    ./hoverfly -generate-ca-cert -cert-name my_certificate -cert-org my_organization

This will create a certificate with the name "my_certificate" for the organization "my_organization" in the current directory.

    -cert <string>

Path to the certificate file to use.
    
    -key <string>
        
Path to the key file to use.  

    -tls-verification=<string>

Turn on/off TLS verification for outgoing requests (Hoverfly will not try to verify certificates) - defaults to true.
	      
TLS verification can also be turned on/off with the `HoverflyTlsVerification` environment variable. For example `HoverflyTlsVerification=false`.	      

## Logging & metrics

    -v	

Verbose mode. Logs every proxy request to stdout.

    -metrics

Logs metrics to stdout.


## Misc

    -dev

Supply -dev flag to serve Admin UI static files directly from ./static/dist instead from statik binary.
	