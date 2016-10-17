# Flags and environment variables

Hoverfly can be configured using flags on startup, or using environment variables.

## Authentication
### Enable/disable authentication
#### Flag
    -auth <string>
#### Environment variable
    export HoverflyAuthEnabled=<string>
Supply `true` to enable authentication. Defaults to `false`.
### Add a new user
#### Flags

    -add -username <string> -password <string> -admin <string>

Supply '-admin false' to make this a non-admin user (defaults to 'true').

For example:

    ./hoverfly -add -username username -password password -admin false   	

This creates a new non-admin user with the username 'username' and the password 'password'.

#### Environment variables
    export HoverflyAdmin="username"
    export HoverflyAdminPass="password"

Setting these environment variables will create a new admin user when Hoverfly starts. 

### Set Hoverfly secret
By default, a random secret is generated every time Hoverfly starts.
#### Environment variable
    export HoverflySecret=<string>

### Set API token expiration (in seconds)
Set to one day by default.
#### Environment variable
    export HoverflyTokenExpiration=<string>

## Port selection
### Set the Admin UI/API port

Defaults to 8888.

#### Flag
    -ap <string>

#### Environment variable

    export AdminPort=<string>
### Set the proxy port
Defaults to 8500.
#### Flag

    -pp <string>

#### Environment variable
    export ProxyPort=<string>

## Mode selection, import & middleware
By default, Hoverfly starts in *simulate mode*.
### Set capture mode
#### Flag

    -capture

### Set synthesize mode
Requires middleware to be specified.
#### Flag

    -synthesize
### Set modify mode
Requires middleware to be specified.
#### Flag
    -modify

### Specify middleware
#### Flag

    -middleware <string>

Supply the path to the middleware script.

For example:

    ./hoverfly -synthesize -middleware "scripts/gen_response.py"

### Import service data
#### Flag

    -import <string>

Import a service data JSON file from file system or URL. For example:

    ./hoverfly -import http://mypage.com/service_x.json

    ./hoverfly -import path/to/my/service_x.json      
#### Environment variable
    export HoverflyImport=<string>

For example:
     
    export HoverflyImport="http://mypage.com/service_x.json"

## Webserver
### Turn Hoverfly into a simulation webserver
####
```
-webserver
```

## Destination

### Specify which hosts to process
#### Flag
```
-dest <string>
```

For example:

    ./hoverfly -dest fooservice.org -dest barservice.org -dest catservice.org

This will start Hoverfly in *simulate mode*, and only simulate requests that are sent to fooservice.org, barservice.org and catservice.org. Requests to all other hosts will pass through.

### Specify host URI
Use regular expression. Defaults to "."
#### Flag

    -destination <string>

## Persistence
### Specify BoltDB or in-memory
####Flag

    -db <string>

By default, Hoverfly uses BoltDB to store data in a file on disk. Use `-db memory` to disable this and use in-memory persistence only.
### Set BoltDB file
By default, a `requests.db` file is created in the Hoverfly directory.
#### Flag

    -db-dir <string>
#### Environment variable
    export HoverflyDB=<string>
The file will be created if it doesn't exist.	    	

## TLS & Certificate management
### Generate certificate
Hoverfly will generate private and public keys in the current directory.
#### Flags
    -generate-ca-cert -cert-name <string> -cert-org <string>

Certificate name defaults to "hoverfly.proxy". Organization name defaults to "Hoverfly Authority".
### Use certificate and key
Supply paths to certificate and key file.
#### Flags
    -cert <string> -key <string>

### Turn off TLS verification
Defaults to "true".
#### Flag

    -tls-verification=<string>

#### Environment variable

    export HoverflyTlsVerification=<string>
    
## Logging & metrics
### Enable verbose mode
Logs every proxy request to STDOUT.
#### Flag
    -v

### Enable metrics logging
Logs metrics to STDOUT.
#### Flag

    -metrics

## Misc
### Use uncompiled static files
Serve Admin UI static files directly from ./static/dist instead from statik binary.
####Flag
    -dev
    
### Get version
Get the version of Hoverfly
####Flag
    -version    
