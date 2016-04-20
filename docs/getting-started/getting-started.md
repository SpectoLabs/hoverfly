# Getting started

## Get Hoverfly

### Pre-built binary

Hoverfly binaries for all major OSes are [available here.](https://github.com/SpectoLabs/hoverfly/releases)

To run Hoverfly, ensure that the correct permissions are set, then execute the binary file. Unless you have disabled authentication, or created a user using [flags or environment variables](../reference/flags_environment_variables.md), you will be prompted to create a new user.

### Docker image

To get the Hoverfly Docker image:

    docker pull spectolabs/hoverfly

To run the image:

    docker run -d \
               -p 8888:8888 \
               -p 8500:8500 \
               spectolabs/hoverfly:latest \
               -add -username <my_username> -password <my_password>

The port mapping is for the Hoverfly AdminUI/API and proxy respectively.

The Docker image supports [Hoverfly flags](../reference/flags_environment_variables.md). In the command above, the `-add`, `-username` and `-password` flags are used to create a user.


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

Homebrew formula and .deb/.rpm packages are on the roadmap.


## Hoverfly is an HTTP(s) proxy

To use Hoverfly in your development or test environment, you need to ensure that your application is using it as a proxy. This can be done at the OS level, or at the application level, depending on the environment.

### Linux, Unix, OSX

Set the HTTP_PROXY or HTTPS_PROXY environment variable to point to Hoverfly:

    export HTTP_PROXY=http://<HOVERFLY_HOST>:8500  
    export HTTPS_PROXY=https://<HOVERFLY_HOST>:8500

Hoverfly uses port 8500 for the proxy by default, although this is configurable.

If you just want to experiment with Hoverfly, you can tell cURL to use Hoverfly as a proxy when you execute cURL command:

    curl http://mirage.readthedocs.org --proxy http://<HOVERFLY_HOST>:8500/

### Windows

TODO

## Admin UI

The Hoverfly Admin UI is available on

    http://<HOVERFLY_HOST>:8888

by default. This is configurable.

The Admin UI can be used to change the Hoverfly mode, and to view basic analytic information. It uses the [Hoverfly API](../reference/api.md).
