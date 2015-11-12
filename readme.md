# One Proxy to Integrate Them All

-

## Configuration

Specifying which site to record/playback with regular expression:
./GenProxy --destination="^.*:80$"

By default proxy is always in playback mode. To switch to record mode, add "--record" flag during startup:
./GenProxy --record


Do a curl request with proxy details: 
+ curl http://mirage.readthedocs.org --proxy http://localhost:8500/

### Playback

Start proxy in playback mode:
./GenProxy

## API

Access admin panel under default port 8888:

* Recorded requests: http://proxy_hostname:8888/records
