# Lambda 503
This is an example of a AWS Lambda function written in Python which can be used as middleware. This example will will randomly delay your request for up to 2 seconds and then return an Nginx 503 internal server error page.

You can configure Hoverfly to run this middleware using hoverctl:
```
hoverctl middleware <link to AWS Lambda endpoint>
```

To find out more, please check the documentation regarding [using middleware](https://spectolabs.gitbooks.io/hoverfly/content/usage/using_middleware.html).

If you want to learn more about AWS Lambda, please check their documentation regarding [getting started](http://docs.aws.amazon.com/lambda/latest/dg/getting-started.html).
