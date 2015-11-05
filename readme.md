# Generic E2E testing proxy

-

## Configuration

Specifying which site to mock:
/GenProxy --target="http://mirage.readthedocs.org/"

By default proxy is always in playback mode. To switch to record mode, add "--record" flag during startup:
/GenProxy --record