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

* [Mocking APIs for development and testing](http://www.specto.io/blog/api-mocking-for-dev-and-test-part-1.html)
* [Using Hoverfly to build Spring Boot microservices alongside a Java monolith](http://www.specto.io/blog/using-api-simulation-to-build-microservices.html)
* [Easy API simulation with the Hoverfly JUnit rule](https://specto.io/blog/hoverfly-junit-api-simulation.html)

## Hoverfly modes


### Capture mode

In this mode, Hoverfly sits between a client application and a server application and transparently intercepts incoming requests and returned responses. These request/response pairs are stored in memory and can be persisted on disk.

This is how you capture a real service or API so you can virtualize it for use in development or test. You could capture real traffic, export it as JSON, modify it to your requirements, then import it back into Hoverfly.

[Read more about capture mode](#)

### Virtualize mode

Once Hoverfly has either captured real traffic or has been populated with pre-made (or pre-captured) virtual services, it can be put into "virtualize" mode. In this mode, Hoverfly will return a response for each request it receives.
   
This will be useful if you have "captured" a real service that you don't have reliable access to, and you want to build a client application that uses it. Or if a team-mate is working on a different service that your application needs to use, and they have shared a Hoverfly virtual service with you.   

[Read more about virtualize mode](#)


### Synthesize mode

In this mode, Hoverfly applies "middleware" (see below) to each incoming request to generate a response on the fly.

This will be useful if you want fine-grained control over what responses Hoverfly returns, and/or if you don't want to mess around capturing traffic and editing JSON files.

[Read more about synthesize mode](#)

### Modify mode

In this mode, Hoverfly sits between the client and the server application, and applies middleware to each request and response.

This could be used for all kinds of things. For example, adding authentication headers to requests from a client application to an external service.

[Read more about modify mode](#)

## Middleware

Hoverfly supports middleware scripts that can be applied to the outgoing requests and/or the incoming responses, depending on the mode. Middleware can be written in any language - as long as it is supported by the machine Hoverfly is running on. You could write middleware scripts in Python, Ruby or Javascript (if you have NodeJS installed).
 
Middleware works in the following ways, depending on the mode:  
 
* Capture mode: middleware affects only outgoing requests.
* Virtualize mode: middleware affects only responses (cache contents remain untouched).
* Synthesize mode: middleware creates responses.
* Modify mode: middleware affects requests and responses. 

[Read more about middleware](#)

