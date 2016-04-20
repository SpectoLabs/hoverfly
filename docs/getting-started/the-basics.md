# The Basics

## Hoverfly use-cases

Its important to understand the different Hoverfly use-cases and modes.

First of all, Hoverfly is intended to cater for two high-level use cases.

### Capturing real HTTP(s) traffic between an application and an external service for re-use in testing or development.

If the external service you want to simulate already exists, you can put Hoverfly in between the client application and the external service and set it to capture every request from the client application and every matching response from the external service (Capture mode).

These request/response pairs are persisted in Hoverfly and can be exported to a JSON file. The JSON can be stored elsewhere (a Git repository, for example), modified as required, then imported into another Hoverfly instance.

You can then configure Hoverfly to act as a surrogate for the external service, returning a matched response for every request it received (Virtualize mode).

This can be useful if you want to create, say, a portable, self-contained version of the Twitter API to develop or test against. With Hoverfly, you can get around the problem of rate-limiting (which can be frustrating when working with a public API), you can write Hoverfly extensions to manipulate the data in pre-recorded responses, you can work while offline, or you can speed up your workflow by replacing a slow dependency with an extremely fast Hoverfly surrogate.

More information on these use-cases is available here:

* [Creating fast versions of slow dependencies](http://www.specto.io/blog/speeding-up-your-slow-dependencies.html)
* [Virtualizing the Meetup API](http://www.specto.io/blog/hoverfly-meetup-api.html)


### Manually creating simulated services for use in a testing or development.

In some cases, the external service you want to simulate might not exist yet. In this case, you can create service simulations manually by writing Hoverfly JSON files. This use-case is in line with the principle of "contract-first" development, where developers provide a "contract" for their service so that other developers can start developing against the service before it is complete.

Hoverfly JSON files can be created manually by each developer, then stored in a Git repository. Other developers can then import the simulated service JSON directly from the repository URL, providing them with a Hoverfly surrogate of the service to work with.

Alternatively, instead of writing a JSON file, you can write a script that will make Hoverfly return a response based on the request it receives (Synthesize mode). This logic is implemented using Hoverfly middleware.  

More information on this use-case is available here:

* [Synthetic service example](https://github.com/SpectoLabs/hoverfly/tree/master/examples/middleware/synthetic_flight_search)
* [Easy API simulation with the Hoverfly JUnit rule](https://specto.io/blog/hoverfly-junit-api-simulation.html)

## Hoverfly modes

### Capture mode

In this mode, Hoverfly, acting as a proxy, sits in between the client application and the external service. It transparently intercepts and stores out-going requests from the client and matching incoming responses from the external service.

This is how you capture real traffic for use in development or testing.

[Read more about capture mode](../usage/capture.md)

### Virtualize mode

In this mode, Hoverfly uses either previously captured traffic, or imported JSON files, to mimic the external service.

This is useful if you are developing or testing an application that needs to talk to an external service that you don't have reliable access to. You can use the Hoverfly simulated version of the service instead of the real service.

[Read more about virtualize mode](../usage/virtualize.md)

### Synthesize mode

In this mode, Hoverfly doesn't use any stored request/response pairs. Instead, it generates responses to incoming requests on the fly and returns them to the client. This mode is dependent on *middleware* (see below) to generate the responses.

This is useful if you can't (or don't want to) capture real traffic, or if you don't want to mess around creating JSON files.

[Read more about synthesize mode](../usage/synthesize.md)

### Modify mode

In this mode, Hoverfly passes requests through from the client to the server, and passes the responses back. However, it also executes middleware on the requests and responses.

This is useful for all kinds of things. For example, manipulating the data in requests and/or responses on the fly.

[Read more about modify mode](../usage/modify.md)

## Middleware

Middleware can be written in any language, provided that language is supported by the Hoverfly host. For example, you could write middleware in Go, Python or JavaScript (if you have NodeJS installed).

Middleware is applied to the requests and/or the responses depending on the mode:

* Capture Mode: middleware affects only outgoing requests
* Virtualize Mode: middleware affects only responses (cache contents remain untouched)
* Synthesize Mode: middleware creates responses
* Modify Mode: middleware affects requests and responses

Middleware can be used to do many useful things, such as simulating network latency or failure, rate limits or controlling data in requests and responses.

[Read more about middleware](../usage/middleware.md)
