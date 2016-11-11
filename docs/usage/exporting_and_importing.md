# Managing simulation data

Hoverfly can export and import service data in JSON format. This is useful if:

* You have captured some traffic and want to store it somewhere other than the Hoverfly `requests.db` file - in a Git repository for example.
* You are running Hoverfly in a Docker container - so persisting data on disk is not ideal
* You want to capture traffic, then modify it somehow before re-importing it
* You want to share your service data someone else   

## Exporting captured data
Using hoverctl, you can export all the simulation data from Hoverfly. Exported data will be written to a JSON file in your current working directory.

```
hoverctl export mysimulation.json
```

For more information about hoverctl, [check here](../reference/hoverctl.md).

## Simulation data JSON format

Hoverfly stores captured **Request Response Pairs** (i.e. "traffic") in the following JSON structure:

    {
    	"data": {
    		"pairs": [{
    			"response": {
    				"status": 200,
    				"body": "body here",
    				"encodedBody": false,
    				"headers": {
    					"Content-Type": ["text/html; charset=utf-8"]
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
    					"Accept": ["text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"],
    					"User-Agent": ["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"]
    				}
    			}
    		}]
    	}
    }

When you export a simulation that you have captured, the JSON file will look something like this. Notice that by default, the request `requestType` is `recording`.

### Base64 encoding binary data

As JSON does not support binary data, binary responses are base64 encoded.  This is denoted by the `encodedBody` field.  Hoverfly automatically encodes and decodes the data during the export and import processes.

### A note about matching

Hoverfly works by inspecting requests being made, extracting key pieces of information and then matching them against stored requests. Standard matching uses everything in the request apart from headers. (This because request headers often change depending on browser, HTTP client and or the time of day.)

In some cases, you may want to use partial matching. For example, you may want Hoverfly to return a specific response for **any** incoming request going to a specific path. This can be achieved using request templates.

### Request templates (for partial matching)

If no exact match is found for an incoming request, Hoverfly will attempt to match on request templates.

Request templates are defined in the JSON file by setting the `requestType` property for a request to `template` and including **only** the information in the request that you want Hoverfly to use in the match.

For example:

    {
        "data": {
            "pairs": [{
                "response": {
                    "status": 200,
                    "body": "<h1>Matched on template</h1>",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": ["text/html; charset=utf-8"]
                    }
                },
                "request": {
                    "requestType": "template",
                    "path": "/template"
                }
            }]
        }
    }

Here, any request with the path `/template` will return the same response.

For looser matching, it is possible to use a wildcard to substitute characters. This is achieved by using an `*` symbol. This will match any number of characters, and is case sensitive.

It is possible to combine the wildcard (`*`) with characters to substitute parts of a string. In the next example, we use a wildcard the replace part of a URL path. This allows us to match on either `/api/v1/template` or `/api/v2/template`.

    {
        "data": {
            "pairs": [{
                "response": {
                    "status": 200,
                    "body": "<h1>Matched on template</h1>",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": ["text/html; charset=utf-8"]
                    }
                },
                "request": {
                    "requestType": "template",
                    "path": "/api/*/template"
                }
            }]
        }
    }

The JSON file can contain both requests recordings and request templates:

    {
        "data": {
            "pairs": [{
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
            }, {
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
            }],
            "globalActions": {
                "delays": []
            }
        }
    }

A standard workflow might be:

1. Capture some traffic
2. Export it to JSON
3. Edit the JSON to set certain requests to templates, removing the properties for these requests that should be excluded from the match
4. Re-import the JSON to Hoverfly

If the `requestType` property is not defined or not recognized, Hoverfly will treat a request as a recording.

## Importing simulation data

Service data can be imported from the local file system, or from a URL. After you have edited captured data to include some request templates, you will probably want to import it back into Hoverfly. There are three ways to do this.

1. Use hoverctl to load a file into Hoverfly:

       hoverctl import simulation.json

For more about hoverctl, [check here](../reference/hoverctl.md).

2. Start Hoverfly in *simulate mode* with the `-import` flag:

       ./hoverfly -import path/to/data.json
       ./hoverfly -import https://<MY_HOST>/data.json

  Multiple service data files can be imported like this:

       ./hoverfly -import path/to/data_1.json -import path/to/data_2.json

   If the file you specified cannot be found, Hoverfly will not start.    

2. Use the `HoverflyImport` environment variable:

       export HoverflyImport="path/to/data.json"
       export HoverflyImport="https://<MY_HOST>/data.json"

    If the file you specified cannot be found, Hoverfly will not start.  

3. Make an API call:

       curl --data "@/path/to/data.json" http://${HOVERFLY_HOST}:8888/api/v2/simulation

For each service data file that has been imported, metadata containing the imported service data source will be stored in Hoverfly (see the **Using the metadata API** section).
