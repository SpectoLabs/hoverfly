# Authentication to Hoverfly

Hoverfly uses a combination of basic auth and JWT (JSON Web Tokens) to authenticate users

## Usage

Add new user:

    ./hoverfly -v -add -username hfadmin -password hfadminpass 

Getting token:

    curl -H "Content-Type application/json" -X POST -d '{"Username": "hoverfly", "Password": "testing"}' http://localhost:8888/token-auth

Using token:

    curl -H "Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NTYxNTY3ODMsImlhdCI6MTQ1NTg5NzU4Mywic3ViIjoiIn0.Iu_xBKzBWlrO70kDAo5hE4lXydu3bQxDZKriYJ4exg3FfZXCqgYH9zm7SVKailIib9ESn_T4zU-2UtFT5iYhw_fzhnXtQoBn5HIhGfUb7mkx0tZh1TJBkLCv6y5ViPw5waAnFBRcygh9OdeiEqnJgzHKrxsR87EellXSdMn2M8wVIhjIhS3KiDjUwuqQl-ClBDaQGlsLZ7eC9OHrJIQXJLqW7LSwrkV3rstCZkTKrEZCdq6F4uAK0mgagTFmuyaBHDEccaivkgYDcaBb7n-Vmyh-jUnDOnwtFnrOv_myXlqqkvtezfm06MBl4PzZE6ZtEA5XADdobLfVarbvB9tFbA" http://localhost:8888/records
