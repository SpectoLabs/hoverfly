# modify_response

This is an example of middleware. This middleware will replace the response body with `"body was replaced by middleware"`.

You can test middleware by using Hoverfly in `modify` mode.

```
hoverctl start
hoverctl mode modify
```

### JavaScript
```
hoverctl middleware --binary node --script modify_response.js
```

### Python
```
hoverctl middleware --binary python --script modify_response.py
```

### Ruby
```
hoverctl middleware --binary ruby --script modify_response.rb
```

To find out more, please check the documentation regarding [middleware](https://docs.hoverfly.io/en/latest/pages/keyconcepts/middleware.html).

