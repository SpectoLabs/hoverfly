# Installation and setup

## Get Hoverfly
Hoverfly is a single binary file. It comes with an optional command line interface tool called hoverctl.

Zip archives containing the Hoverfly and hoverctl binaries for Windows, MacOS and Linux are available on the GitHub releases page:

[Hoverfly & hoverctl zip archives](https://github.com/SpectoLabs/hoverfly/releases/latest)

Download the archive for your OS, extract the Hoverfly and hoverctl binaries and move them to a directory on your [PATH](https://www.java.com/en/download/help/path.xml).

### Homebrew (MacOS)

```
brew install SpectoLabs/tap/hoverfly
```


## Run Hoverfly

To capture traffic between your application and an external service, you will need to configure your OS, browser or application to use Hoverfly as a proxy.

### MacOS & Linux

Run Hoverfly using hoverctl:
```
hoverctl start
```

By default, the Hoverfly proxy runs on localhost:8500. Switch Hoverfly to "capture" mode and make a request with cURL, using Hoverfly as a proxy:
```
hoverctl mode capture
curl --proxy http://localhost:8500 http://hoverfly.io/
```

Hoverfly has captured the request and the response. View the Hoverfly logs:

```
hoverctl logs
```

Switch Hoverfly to "simulate" mode" and make the same request:
```
hoverctl mode simulate
curl --proxy http://localhost:8500 http://hoverfly.io/
```
Hoverfly has returned the captured response.

### Windows

Open a command prompt and run Hoverfly using hoverctl:
```
hoverctl start
```

Configure your application, browser or OS to use the Hoverfly proxy (http://localhost:8500). Switch Hoverfly to "capture" mode:

```
hoverctl mode capture
```

Make some requests from your application, browser or OS, then view the Hoverfly logs:

```
hoverctl logs
```

Switch Hoverfly to "simulate" mode:

```
hoverctl mode simulate
```

Make the same requests from your browser, OS or application. Hoverfly is returning the captured responses.

More information on proxy settings:

* [Windows proxy settings explains](http://blog.raido.be/?p=426)
* [Firefox proxy settings](https://support.mozilla.org/en-US/kb/advanced-panel-settings-in-firefox#w_connection)
* [Java Networking and Proxies](https://docs.oracle.com/javase/6/docs/technotes/guides/net/proxies.html)

## Docker image

The Hoverfly docker image contains only the Hoverfly binary.

    docker run -d -p 8888:8888 -p 8500:8500 spectolabs/hoverfly

The port mapping is for the Hoverfly AdminUI/API and proxy respectively. The Docker image supports Hoverfly flags (see **Flags and environment variables** in the **Usage** section).

## JUnit rule

The Hoverfly JUnit rule is available as a Maven dependency. To get started, add the following to your pom.xml:

    <groupId>io.specto</groupId>
    <artifactId>hoverfly-junit</artifactId>
    <version>0.1.4</version>

This will download the JUnit rule and the Hoverfly binary. More information on how to use the Hoverfly JUnit rule is available here:

[Easy API Simulation with the Hoverfly JUnit Rule](https://specto.io/blog/hoverfly-junit-api-simulation.html)         


## Setting the HOVERFLY_HOST environment variable

Throughout the documentation, `${HOVERFLY_HOST}` is used in API examples and Admin UI guides. If you are running a binary on your local machine, the Hoverfly host will be `localhost`. However, if you are running Hoverfly on a remote machine, or using Docker Machine for example, it will be different.

To make things easier when following the documentation, it is recommended that you set the HOVERFLY_HOST environment variable. For example:

    export HOVERFLY_HOST=localhost

## Hoverfly as an HTTP(S) proxy

Hoverfly is primarily a proxy - although it can run (with limitations) as a webserver. To use Hoverfly in your development or test environment to capture traffic, you need to ensure that your application is using it as a proxy. This can be done at the OS level, or at the application level, depending on the environment.


## Admin UI

The Hoverfly Admin UI is available on port 8888 by default:

    http://${HOVERFLY_HOST}:8888

The port is configurable (see the **Flags and environment variables** section).

When authentication is disabled (which is the default), you can use **any username and password combination** to access the Admin UI.

The Admin UI can be used to change the Hoverfly mode, and to view basic analytic information. It uses the Hoverfly API.
