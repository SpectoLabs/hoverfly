# API

## GET /api/v2/simulation
Gets all simulation data being used in the running instance of Hoverfly. This includes recordings, templates, delays and meta data.
```
{
  "data": {
    "pairs": [
      {
        "response": {
          "status": 200,
          "body": "<h1>Matched on recording</h1>",
          "encodedBody": false,
          "headers": {
            "Content-Type": [
              "text/html; charset=utf-8"
            ]
          }
        },
        "request": {
          "requestType": "recording",
          "path": "/",
          "method": "GET",
          "destination": "myhost.io",
          "scheme": "https",
          "query": "",
          "body": "",
          "headers": {
            "Accept": [
              "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
            ],
            "Content-Type": [
              "text/plain; charset=utf-8"
            ],
            "User-Agent": [
              "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"
            ]
          }
        }
      },
      {
        "response": {
          "status": 200,
          "body": "<h1>Matched on template</h1>",
          "encodedBody": false,
          "headers": {
            "Content-Type": [
              "text/html; charset=utf-8"
            ]
          }
        },
        "request": {
          "requestType": "template",
          "path": "/template",
          "method": null,
          "destination": null,
          "scheme": null,
          "query": null,
          "body": null,
          "headers": null
        }
      }
    ],
    "globalActions": {
      "delays": []
    }
  },
  "meta": {
    "schemaVersion": "v1",
    "hoverflyVersion": "v0.9.0",
    "timeExported": "2016-11-11T11:53:52Z"
  }
```

## PUT /api/v2/simulation
Puts the simulation into Hoverfly and replaces the previous set of data.
```
{
  "data": {
    "pairs": [
      {
        "response": {
          "status": 200,
          "body": "<h1>Matched on recording</h1>",
          "encodedBody": false,
          "headers": {
            "Content-Type": [
              "text/html; charset=utf-8"
            ]
          }
        },
        "request": {
          "requestType": "recording",
          "path": "/",
          "method": "GET",
          "destination": "myhost.io",
          "scheme": "https",
          "query": "",
          "body": "",
          "headers": {
            "Accept": [
              "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
            ],
            "Content-Type": [
              "text/plain; charset=utf-8"
            ],
            "User-Agent": [
              "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"
            ]
          }
        }
      },
      {
        "response": {
          "status": 200,
          "body": "<h1>Matched on template</h1>",
          "encodedBody": false,
          "headers": {
            "Content-Type": [
              "text/html; charset=utf-8"
            ]
          }
        },
        "request": {
          "requestType": "template",
          "path": "/template",
          "method": null,
          "destination": null,
          "scheme": null,
          "query": null,
          "body": null,
          "headers": null
        }
      }
    ],
    "globalActions": {
      "delays": []
    }
  },
  "meta": {
    "schemaVersion": "v1",
    "hoverflyVersion": "v0.9.0",
    "timeExported": "2016-11-11T11:53:52Z"
  }
```

## GET /api/v2/hoverfly
Gets all configuration for the running instance of Hoverfly.
```
{
    destination: ".",
    middleware: "",
    mode: "simulate",
    usage: {
        counters: {
            capture: 0,
            modify: 0,
            simulate: 0,
            synthesize: 0
        }
    }
}
```

## GET /api/v2/hoverfly/destination
Gets the current destination for the running instance of Hoverfly.
```
{
    destination: "."
}
```

## PUT /api/v2/hoverfly/destination	
Puts the new destination and overwrites current destination for the running instance of Hoverfly. This requires a JSON body on the request.

```
{
    destination: "new-destination"
}
```


## GET /api/v2/hoverfly/middleware
Gets the current middleware value for the running instance of Hoverfly. This is likely to be an executable command and path to middleware script being used.
```
{
	"middleware": "python ~/middleware.py"
}
```

## PUT /api/v2/hoverfly/middleware
Puts the new middleware value and overwrites current middleware value for the running instance of Hoverfly. This requires a JSON body on the request. The value you send should be a command that runs on the machine Hoverfly is running. You may need to specify the command to run the script as well as the path to the script.
{
	"middleware": "python ~/new-middleware.py"
}

## GET /api/v2/hoverfly/mode
Gets the current mode for the running instance of Hoverfly.
```
{
    mode: "simulate"
}
```

## PUT /api/v2/hoverfly/mode
Puts the new mode and overwrites current mode of the running instace of Hoverfly. This requires a JSON body on the request.
```
{
    mode: "simulate"
}
```
                                       
## GET /api/v2/hoverfly/usage
Gets the metrics for the running instance of Hoverfly.
```
{
	"metrics": {
		"counters": {
			"capture": 0,
			"modify": 0,
			"simulate": 0,
			"synthesize": 0
		}
	}
}
```
