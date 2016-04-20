# Flags and environment variables

Hoverfly can be configured using flags on startup, or using environment variables.

## Admin UI Authentication

    -no-auth

Disable authentication, currently it is enabled by default.

    -add

Add a new user.

    -username <string>

Username for new user.

    -password <string>

Password for new user.

    -admin <bool>

Supply '-admin false' to make this a non-admin user (defaults to 'true').

For example:

    ./hoverfly -add -username hfadmin -password hfpass -admin false   	

This creates a new non-admin user with the username 'hfadmin' and the password 'hfpass'.

## Port selection

    -ap <string>

Sets the Admin UI port.

    -pp <string>

Sets the proxy port. For example:

    ./hoverfly -ap 1234 -pp 4567         

This starts Hoverfly with the Admin UI on port 1234 and the proxy on 4567.    	

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

For example:

    ./hoverfly -synthesize -middleware "scripts/gen_response.py"

This starts Hoverfly in synthesize mode with a middleware script that generates responses on the fly.

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

Persistent storage to use. By default, Hoverfly uses BoltDB to store data in a file on disk. Specify 'memory'
to disable this and use in-memory cache only. 
Note: 'memory' should only be used for small test cases, if your machine runs out of RAM while Hoverfly is in 'memory' mode - 
bad things can happen.

    -db-path <string>

Path to BoltDB data file. By default, a "requests.db" file will be created in the directory from which Hoverfly is executed. 
Supply a custom path with filename to use a different file or location. The file will be created if it doesn't exist.

For example:

	 ./hoverfly -db-path new_name.db

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

## Logging & metrics

    -v

Verbose mode. Logs every proxy and admin UI/API request to stdout.

    -metrics

Logs metrics to stdout every 5 seconds.


## Misc

    -dev

Supply -dev flag to serve Admin UI static files directly from ./static/admin/dist instead from statik binary. Useful when 
developing UI.

# Environment variables

You can configure Hoverfly through environment variables, this is a standard approach when application is running in a 
Docker container or hosted via platform. 

## Admin UI authentication

    HoverflyAuthDisabled  
    
Hoverfly authentication setting. Defaults to false. Set it to 'true' to disable authentication.
    
For example:
    
    export HoverflyAuthDisabled="true"
    
    HoverflySecret
    
Secret used to generate authentication tokens. If this variable is not set - Hoverfly generates one during startup.
Note: if you leave this value unset - you will have to authenticate in the admin UI after each Hoverfly restart.

    HoverflyTokenExpiration
    
Token expiration time in seconds. Defaults to one day (24 * 60 * 60).

### Setting admin user through env variables

You can provide credentials for initial user by setting these variables:

    HoverflyAdmin
    
Admin user name.

    HoverflyAdminPass
    
Admin password.
Note: if you do not set initial user through environment variables with authentication enabled - Hoverfly will ask you to input
username and password for the first user during startup. This could result in a 'stuck' container.

## Port selection 

    AdminPort
    
Admin UI port, defaults to 8500.
    
    ProxyPort
    
Proxy port, defaults to 8888.
    
## Persistence
 
    HoverflyDB

Path to BoltDB data file. By default, a "requests.db" file will be created in the directory from which Hoverfly is executed. 
Supply a custom path with filename to use a different file or location. The file will be created if it doesn't exist.

## TLS

    HoverflyTlsVerification
    
TLS verification. Since Hoverfly is making requests on behalf of it's users - it has to either trust remote servers or not.
This setting defaults to 'true' which means that untrusted hosts won't be accepted if request is done via HTTPs. You can turn it
off by providing "false" value to this environment variable, for example:

    export HoverflyTlsVerification="false"
    
    
    