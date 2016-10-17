# Using the metadata API

The Hoverfly metadata API provides simple key/value storage to make management of multiple Hoverfly instances easier. 

Currently, metadata is automatically added for each service data file that is imported (see the **Exporting and importing** section).

You can also add arbitrary key/value pairs yourself - this can be useful for naming a Hoverfly instance, for example.

## Adding some metadata
To give the Hoverfly instance a name, make an API call:

    curl -H "Content-Type application/json" -X PUT -d '{"key":"name", "value": "My Hoverfly"}' http://${HOVERFLY_HOST}:8888/api/metadata

You could also add a description for the Hoverfly instance:

    curl -H "Content-Type application/json" -X PUT -d '{"key":"description", "value": "Simulates keystone service, use user XXXX and password YYYYY to login"}' http://${HOVERFLY_HOST}:8888/api/metadata
    
## Retrieving metadata

To retrieve metadata, make an API call:

    curl http://${HOVERFLY_HOST}:8888/api/metadata

This will return the following JSON:

    {
        "data": {
            "description": "Simulates keystone service, use user XXXX and password YYYYY to login",
            "name": "My Hoverfly"
        }
    }

If you imported service data files (see the **Exporting and importing** section), you will see something like this:

    {
        "data": {
            "description": "Simulates keystone service, use user XXXX and password YYYYY to login",
            "import_1": "path/to/my/service_1.json",
            "import_2": "path/to/my/service_2.json",
            "name": "My Hoverfly"
        }
    }

## Deleting metadata

To delete all metadata from a Hoverfly instance, make an API call:

    curl -X DELETE http://${HOVERFLY_HOST}:8888/api/metadata 