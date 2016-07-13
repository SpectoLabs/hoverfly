[![Circle CI][CircleCI-Image]][CircleCI-Url]

![Hoverfly](core/static/img/hoverfly_logo.png)
## Dependencies without the sting

Hoverfly is a light-weight open source tool for creating simulations of external services for use in development and testing. 
This technique is sometimes referred to as [service virtualization](https://en.wikipedia.org/wiki/Service_virtualization).

Hoverfly was designed to provide you with the means to create your own "dependency sandbox": a simulated development and test environment that you control.

## Hoverfly + hoverctl

Hoverctl is a command line interface for Hoverfly. It uses the Hoverfly REST API to control local or remote Hoverfly instances and manage Hoverfly data.

## Quickstart

Run the install script to install Hoverfly and hoverctl:

    curl -o install.sh https://storage.googleapis.com/specto-binaries/install.sh && bash install.sh

Or [download the Hoverfly and hoverctl binaries](https://github.com/SpectoLabs/hoverfly/releases), set the correct permissions and copy them to a directory on your path.

Or get the docker image (Hoverfly only):

    docker pull spectolabs/hoverfly
    docker run -d -p 8888:8888 -p 8500:8500 spectolabs/hoverfly

## Documentation

The Hoverfly documentation is [available here](https://www.gitbook.com/book/spectolabs/hoverfly/details).

For convenience:

### Getting started
* [Use cases](https://spectolabs.gitbooks.io/hoverfly/content/getting_started/use_cases.html)
* ["Modes" and middleware](https://spectolabs.gitbooks.io/hoverfly/content/getting_started/modes_and_middleware.html)
* [Installation and setup](https://spectolabs.gitbooks.io/hoverfly/content/getting_started/installation_and_setup.html)

### Usage
* [Capturing traffic](https://spectolabs.gitbooks.io/hoverfly/content/usage/capturing_traffic.html)
* [Simulating services](https://spectolabs.gitbooks.io/hoverfly/content/usage/simulating_services.html)
* [Exporting and importing](https://spectolabs.gitbooks.io/hoverfly/content/usage/exporting_and_importing.html)
* [Using middleware](https://spectolabs.gitbooks.io/hoverfly/content/usage/using_middleware.html)
* [Creating synthetic services](https://spectolabs.gitbooks.io/hoverfly/content/usage/creating_synthetic_services.html)
* [Modifying traffic](https://spectolabs.gitbooks.io/hoverfly/content/usage/modifying_traffic.html)
* [Filtering destination URLs and hosts](https://spectolabs.gitbooks.io/hoverfly/content/usage/filtering_destination_urls_and_hosts.html)
* [Using the metadata API](https://spectolabs.gitbooks.io/hoverfly/content/usage/using_the_metadata_api.html)
* [Admin UI](https://spectolabs.gitbooks.io/hoverfly/content/usage/admin_ui.html)
* [HTTPS & certificate management](https://spectolabs.gitbooks.io/hoverfly/content/usage/certificate_management.html)
* [Authentication](https://spectolabs.gitbooks.io/hoverfly/content/usage/authentication.html)

### Reference
* [API](https://spectolabs.gitbooks.io/hoverfly/content/reference/api.html)
* [Flags and environment variables](https://spectolabs.gitbooks.io/hoverfly/content/reference/flags_and_environment_variables.html)
* [Hoverctl](https://spectolabs.gitbooks.io/hoverfly/content/reference/hoverctl.html)


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
