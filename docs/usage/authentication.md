# Authentication
Hoverfly uses a combination of Basic Auth and [JWT](https://jwt.io/) (JSON Web Tokens) to authenticate users. Authentication is disabled by default.

## Enabling authentication

If you enable authentication, and you haven't created a user using flags or environment variables (see below), you will be prompted to create a new user when you start Hoverfly.

To enable authentication, you can use the `-auth` flag on startup:

    ./hoverfly -auth
    
Or you can use the `HoverflyAuthDisabled` environment variable:

    export HoverflyAuthEnabled=true
    
If the `-auth` flag is supplied **or** the `HoverflyAuthEnabled` environment variable is set to `true`, authentication will be enabled. 

When authentication is disabled, **any username and password combination** can be used to access the Admin UI. 

    
## Adding users

You can add a user using the `-add`, `-username` and `-password` flags at startup:

    ./hoverfly -add -username <username> -password <password>
    
This will add an admin user. To add a non-admin user, use the `-admin` flag:

    ./hoverfly -add -username <username> -password <password> -admin false

You can also add an initial super user using environment variables. This is useful if you are using Hoverfly in Docker, for example:

    export HoverflyAdmin="username"
    export HoverflyAdminPass="password"
    
## Token usage for API authentication

To get the token for a user, make an API call:

    curl -H "Content-Type application/json" -X POST -d '{"Username": "<username>", "Password": "<password>"}' http://${HOVERFLY_HOST}:8888/api/token-auth

To use the token in an API call:

    curl -H "Authorization: Bearer <token>" http://${HOVERFLY_HOST}:8888/api/v2/simulation

By default, tokens expire after one day. You can override this by setting the `HoverflyTokenExpiration` environment variable in seconds:

    export HoverflyTokenExpiration=3600
    
## Setting the Hoverfly secret

By default, a new random secret will be generated every time you launch Hoverfly. However, you can specify a secret using the `HoverflySecret` environment variable:

    export HoverflySecret=<my_secret>

