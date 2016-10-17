# Installation and setup
***Note:*** Although Hoverfly runs on any OS (including Windows) the examples in these docs are currently aimed Linux, Unix, *BSD and OSX users. 

## Get Hoverfly
Hoverfly comes with a command line application called "hoverctl". Currently, Hoverfly can be used either with hoverctl or without.

### Quick install
#### Homebrew (OSX)
     brew install SpectoLabs/tap/hoverfly

#### Installation script (OSX & Linux) 

**NOTE**: *The installation script will prompt you for a **[SpectoLab](https://lab.specto.io) API key**. Since SpectoLab is currently in private beta, you will have to register for an invite to get an API key on the SpectoLab site. **If you don't have an API key, or you just want to get started without using SpectoLab, just enter a null value when prompted for the key.** Without an API key, all hoverctl functionality will work, with the exception of "pushing" and "pulling" Hoverfly simulations to and from SpectoLab.* 

    curl -o install.sh https://raw.githubusercontent.com/SpectoLabs/hoverfly/master/install.sh && bash install.sh

Once the install process is complete, use hoverctl to start an instance of Hoverfly locally:

    hoverctl start

To view the various hoverctl commands:

    hoverctl help

For more information on hoverctl, see **Hoverctl** in the **Reference** section.

### Pre-built binary (Hoverfly only)

Hoverfly binaries for every major OS are [available here.](https://github.com/SpectoLabs/hoverfly/releases) Hoverctl binaries will be available for download once the source is released.

You can control Hoverfly without using the hoverctl CLI via flags and environment variables (see **Flags and environment variables** in the **Usage** section), or via the API directly.

To run Hoverfly, ensure that the binary is executable by setting the correct permissions, then execute the binary file:

    chmod +x hoverfly*
    ./hoverfly*

Hoverfly authentication is disabled by default (see the **Authentication** section). The API examples in these docs assume that authentication is disabled.

### Docker image

To get the Hoverfly Docker image:

    docker pull spectolabs/hoverfly

To run the image:

    docker run -d \
               -p 8888:8888 \
               -p 8500:8500 \
               spectolabs/hoverfly:latest

The port mapping is for the Hoverfly AdminUI/API and proxy respectively.

The Docker image supports Hoverfly flags (see **Flags and environment variables** in the **Usage** section). For example, to put Hoverfly into *capture mode* (see the **Capturing traffic** section) when the container starts:

    docker run -d \
               -p 8888:8888 \
               -p 8500:8500 \
               spectolabs/hoverfly:latest \
               -capture


### Maven and JUnit integration

The Hoverfly JUnit rule is available as a Maven dependency. To get started, add the following to your pom.xml:

    <groupId>io.specto</groupId>
    <artifactId>hoverfly-junit</artifactId>
    <version>0.1.4</version>

This will download the JUnit rule and the Hoverfly binary. More information on how to use the Hoverfly JUnit rule is available here:

[Easy API Simulation with the Hoverfly JUnit Rule](https://specto.io/blog/hoverfly-junit-api-simulation.html)         

### Build from source

You will need Go 1.7 installed. Instructions on how to set up your Go environment are [available here](https://golang.org/doc/code.html).

    mkdir -p "$GOPATH/src/github.com/SpectoLabs/"
    git clone https://github.com/SpectoLabs/hoverfly.git "$GOPATH/src/github.com/SpectoLabs/hoverfly"
    cd "$GOPATH/src/github.com/SpectoLabs/hoverfly"
    make build
    

Then, to run Hoverfly (with authentication disabled):

    target/hoverfly
    
The source for hoverctl has not yet been released (as of 2016-07-07), although it will be available in the next Hoverfly release.

### Setting the HOVERFLY_HOST environment variable

Throughout the documentation, `${HOVERFLY_HOST}` is used in API examples and Admin UI guides. If you are running a binary on your local machine, the Hoverfly host will be `localhost`. However, if you are running Hoverfly on a remote machine, or using Docker on OSX for example, it will be different.

To make things easier when following the documentation, it is recommended that you set the HOVERFLY_HOST environment variable. For example:

    export HOVERFLY_HOST=localhost

### Package management?

Homebrew formula, deb and rpm packages are on the roadmap.


## Hoverfly as an HTTP(S) proxy

Hoverfly is primarily a proxy - although it can run (with limitations) as a webserver. To use Hoverfly in your development or test environment to capture traffic, you need to ensure that your application is using it as a proxy. This can be done at the OS level, or at the application level, depending on the environment.

### Linux, Unix, OSX, *BSD

Set the HTTP_PROXY and/or HTTPS_PROXY environment variable to point to Hoverfly:

    export HTTP_PROXY=http://${HOVERFLY_HOST}:8500/  
    export HTTPS_PROXY=https://${HOVERFLY_HOST}:8500/

Hoverfly uses port 8500 for the proxy by default, although this is configurable via flags or environment variables (see the **Flags and environment variables** section).

If you just want to experiment with Hoverfly, you can make cURL use Hoverfly as a proxy with the `--proxy` flag:

    curl http://mirage.readthedocs.org --proxy http://${HOVERFLY_HOST}:8500/

### Windows

TODO

## Admin UI

The Hoverfly Admin UI is available on port 8888 by default:

    http://${HOVERFLY_HOST}:8888

The port is configurable (see the **Flags and environment variables** section).

When authentication is disabled (which is the default), you can use **any username and password combination** to access the Admin UI.

The Admin UI can be used to change the Hoverfly mode, and to view basic analytic information. It uses the Hoverfly API.
