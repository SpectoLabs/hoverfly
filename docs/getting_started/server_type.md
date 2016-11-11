# Server type
While Hoverfly is primarily a proxy, there are some situations in which you may need to run it as a webserver - for example if setting the host OS or application to use a proxy is not possible or desirable.

Currently, when running as a webserver, Hoverfly can only be set to *simulate* mode. This is useful if you have a simulation that has been created by capturing traffic using a Hoverfly instance running as a proxy, and you want to import and run it in an environment which cannot be configured to use a proxy. 

## Starting Hoverfly as a webserver
If you are using hoverctl to manage your instance of Hoverfly, you can start Hoverfly as a webserver using hoverctl.
```
hoverctl start webserver
```

If you are running the Hoverfly binary, you can specify the webserver flag which will start Hoverfly as a webserver.
```
./hoverfly -webserver
```

**NOTE:** Currently HTTPS is not supported when running Hoverfly as a webserver. HTTPS support when running as a webserver is on the roadmap.

## Simulations and request matching
When running as a webserver, although Hoverfly functionality is limited to simulate mode, Hoverfly still uses the standard simulations.

When Hoverfly is running as a webserver, the server is available on the same port as the proxy. When you make requests to the webserver, responses will be matched in the same way as with the proxy, except the host will be disregarded. This means that if you have a simulation with multiple hosts, they will all be served from the same host.