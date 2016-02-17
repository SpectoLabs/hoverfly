# Exports

This is an example of how data looks after exporting it.

To export data you can use _curl_:

    curl http://localhost:8888/records > requests.json
    
To import it back:

    curl --data "@requests.json" http://localhost:8888/records


You can also use flags to import records on startup from disk or from given url:

    ./hoverfly -import "https://gist.githubusercontent.com/rusenask/a43caf7a87e2b066d3dd/raw/d75c6c3532d631c7816a79a9e2767e6995f2e16f/requests.json"
    
or:
    
    ./hoverfly -import "requests.json"
    