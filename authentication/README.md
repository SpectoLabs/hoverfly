# Authentication to Hoverfly

Hoverfly uses a combination of basic auth and JWT (JSON Web Tokens) to authenticate users

### Authentication (currently disabled by default)

To enable admin interface authentication you can pass '-auth' flag during startup:

    ./hoverfly -auth
    
or supply environment variable:

    export HoverflyAuthEnabled=true
    
If environment variable __or__ flag is given to enable authentication - it will be enabled (if you set flag to 'false' 
but leave environment variable set to 'true', or vice versa - auth will be enabled). 

Export Hoverfly secret:

    export HoverflySecret=VeryVerySecret
    

If you skip this step - a new random secret will be generated every single time when you launch Hoverfly. This can be useful
if you are deploying it in cloud but it can also be annoying if you are working with Hoverfly where it is constantly restarted.

You can also specify token expiration time (defaults to 72):

    export HoverflyTokenExpiration=200

### Adding users

Then, add your first admin user:

    ./hoverfly -v -add -username hfadmin -password hfadminpass
     
You can also create non-admin users by supplying 'admin' flag as follows:

    ./hoverfly -v -add -username hfadmin -password hfadminpass -admin false

Getting token:

    curl -H "Content-Type application/json" -X POST -d '{"Username": "hoverfly", "Password": "testing"}' http://localhost:8888/token-auth

Using token:

    curl -H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTYxNTY3ODMsImlhdCI6MTQ1NTg5NzU4Mywic3ViIjoiIn0.Iu_xBKzBWlrO70kDAo5hE4lXydu3bQxDZKriYJ4exg3FfZXCqgYH9zm7SVKailIib9ESn_T4zU-2UtFT5iYhw_fzhnXtQoBn5HIhGfUb7mkx0tZh1TJBkLCv6y5ViPw5waAnFBRcygh9OdeiEqnJgzHKrxsR87EellXSdMn2M8wVIhjIhS3KiDjUwuqQl-ClBDaQGlsLZ7eC9OHrJIQXJLqW7LSwrkV3rstCZkTKrEZCdq6F4uAK0mgagTFmuyaBHDEccaivkgYDcaBb7n-Vmyh-jUnDOnwtFnrOv_myXlqqkvtezfm06MBl4PzZE6ZtEA5XADdobLfVarbvB9tFbA" http://localhost:8888/records
