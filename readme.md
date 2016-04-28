[![Circle CI][CircleCI-Image]][CircleCI-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]

![Hoverfly](static/img/hoverfly_logo.png)
## Dependencies without the sting

Hoverfly is a light-weight open source tool for creating simulations of external services for use in development and testing. 
This technique is sometimes referred to as [service virtualization](https://en.wikipedia.org/wiki/Service_virtualization).

Hoverfly was designed to provide you with the means to create your own "dependency sandbox": a simulated development and test environment that you control.

## Documentation

The Hoverfly documentation is available here.

For convenience:

### Getting started
* Basic concepts
* Installation and setup

### Usage
* Capturing traffic
* Simulating services
* Exporting and importing
* Using middleware
* Creating synthetic services
* Modifying traffic
* Filtering destination URLs and hosts
* Using the metadata API
* Admin UI
* Certificate management
* Authentication

### Reference
* API
* Flags and environment variables


## Further reading

Articles and blog posts with step-by-step Hoverfly tutorials:

* [Creating fast versions of slow dependencies](http://www.specto.io/blog/speeding-up-your-slow-dependencies.html)
* [Modifying traffic on the fly](http://www.specto.io/blog/service-virtualization-is-so-last-year.html)
* [Mocking APIs for development and testing](http://www.specto.io/blog/api-mocking-for-dev-and-test-part-1.html)
* [Virtualizing the Meetup API](http://www.specto.io/blog/hoverfly-meetup-api.html)
* [Using Hoverfly to build Spring Boot microservices alongside a Java monolith](http://www.specto.io/blog/using-api-simulation-to-build-microservices.html)
* [Easy API simulation with the Hoverfly JUnit rule](https://specto.io/blog/hoverfly-junit-api-simulation.html)

## Contributing

Contributions are welcome!

To submit a pull request you should fork the Hoverfly repository, and make your change on a feature branch of your fork.
Then generate a pull request from your branch against master of the Hoverfly repository. Include in your pull request
details of your change (why and how, as well as the testing you have performed). To read more about forking model, check out
this link: [forking workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/forking-workflow).

Hoverfly is a new project, we will soon provide detailed roadmap.

## License

Apache License version 2.0 [See LICENSE for details](https://github.com/SpectoLabs/hoverfly/blob/master/LICENSE).

(c) [SpectoLabs](https://specto.io) 2016.

[CircleCI-Image]: https://circleci.com/gh/SpectoLabs/hoverfly.svg?style=shield
[CircleCI-Url]: https://circleci.com/gh/SpectoLabs/hoverfly
[ReportCard-Url]: http://goreportcard.com/report/spectolabs/hoverfly
[ReportCard-Image]: http://goreportcard.com/badge/spectolabs/hoverfly
