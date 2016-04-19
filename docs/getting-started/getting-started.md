# Getting started

## Get Hoverfly

### Pre-built binary

Hoverfly binaries for all major OSes are [available here.](https://github.com/SpectoLabs/hoverfly/releases)

To run the binary, ensure that the correct permissions are set, then execute the binary file.

### Docker image

To get the Hoverfly Docker image:

    docker pull spectolabs/hoverfly

To run the image:

    docker run -d -p 8888:8888 -p 8500:8500 spectolabs/hoverfly:latest


### Maven and JUnit integration

The Hoverfly JUnit rule is available as a Maven dependency. To get started, add the following to your pom.xml:

    <groupId>io.specto</groupId>
    <artifactId>hoverfly-junit</artifactId>
    <version>0.1.1</version>

This will download the JUnit rule and the Hoverfly binary. More inforation on how to use the Hoverfly JUnit rule is available here:

[Easy API Simulation with the Hoverfly JUnit Rule](https://specto.io/blog/hoverfly-junit-api-simulation.html)         

### Build from source

You will need Go 1.6 installed. Instructions on how to set up your Go environment are [available here](https://golang.org/doc/code.html).

    mkdir -p "$GOPATH/src/github.com/SpectoLabs/"
    git clone https://github.com/SpectoLabs/hoverfly.git "$GOPATH/src/github.com/SpectoLabs/hoverfly"
    cd "$GOPATH/src/github.com/SpectoLabs/hoverfly"
    make build

To run the binary:

    ./cmd/hoverfly/hoverfly
        

### Package management?

Homebrew formula and .deb packages are on the roadmap.


## Hoverfly is an HTTP(s) proxy

To use Hoverfly in your development or test environment, you need to ensure that your application is using it as a proxy. This can be done at the OS level, or at the application level, depending on the environment.

### Linux, Unix, OSX

Set the HTTP_PROXY or HTTPS_PROXY environment variable to point to Hoverfly:

    export HTTP_PROXY=http://<HOVERFLY_HOST>:8500  
    export HTTPS_PROXY=https://<HOVERFLY_HOST>:8500

Hoverfly uses port 8500 for the proxy by default, although this is configurable.

If you just want to experiment with Hoverfly, you can tell cURL to use Hoverfly as a proxy when you execute cURL command:

    curl blah TODO

### Windows

TODO

## Admin UI

The Hoverfly Admin UI is available on 

    http://<HOVERFLY_HOST>:8888 
    
by default. This is configurable.

The Admin UI can be used to change the Hoverfly mode, and to view basic analytic information. uses the Hoverfly API (link TODO).
