FreePort
========

Get a free open TCP port that is ready to use

```bash
# Ask the kernel to give us an open port.
export port=$(freeport)

# Start standalone httpd server for testing
httpd -X -c "Listen $port" &

# Curl local server on the selected port
curl localhost:$port
```

#### Binary Downloads
 - Mac:   https://phayes.github.io/bin/current/freeport/mac/freeport.gz
 - Linux: https://phayes.github.io/bin/current/freeport/linux/freeport.gz

#### Building From Source
```bash
sudo apt-get install golang                    # Download go. Alternativly build from source: https://golang.org/doc/install/source
mkdir ~/.gopath && export GOPATH=~/.gopath     # Replace with desired GOPATH
export PATH=$PATH:$GOPATH/bin                  # For convenience, add go's bin dir to your PATH
go get github.com/phayes/freeport/cmd/freeport
```
