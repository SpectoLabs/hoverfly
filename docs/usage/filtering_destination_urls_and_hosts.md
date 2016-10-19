# Filtering destination URLs and hosts
You may wish to control what Hoverfly captures or simulates. By default, Hoverfly will process everything. 

To specify which URLs Hoverfly processes you can either use the `-destination` flag on startup to supply a regular expression:

    ./hoverfly -destination="<my_regex>"
    
Or you can supply multiple hosts using the `-dest` flag on startup:

    ./hoverfly -dest www.myservice1.org -dest www.myservice2.org -dest www.myservice3.org

You can also set the destination host using the API:

    curl -H "Content-Type application/json" -X POST -d '{"destination": "www.myservice1.org"}' http://${HOVERFLY_HOST}:8888/api/state

Hoverfly will then process only those requests which match the specified destination. 

All other requests will be passed through. This allows you to start by simulating just a few services.

