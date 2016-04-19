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

TODO

### Manually creating simulated services for use in a testing or development.

In some cases, the external service you want to simulate might not exist yet. In this case, you can create service simulations manually by writing Hoverfly JSON files. This use-case is in line with the principle of "contract-first" development, where developers provide a "contract" for their service so that other developers can start developing against the service before it is complete.

Hoverfly JSON files can be created manually by each developer, then stored in a Git repository. Other developers can then import the simulated service JSON directly from the repository URL, providing them with a Hoverfly surrogate of the service to work with.

Alternatively, instead of writing a JSON file, you can write a script that will make Hoverfly return a response based on the request it receives (Synthesize mode). This logic is implemented using Hoverfly middleware.  

More information on this use-case is available here:

TODO  

### Hoverfly modes

### Virtualize mode

### Capture mode

### Synthesize mode

### Modify mode

### Bypass

## Middleware
