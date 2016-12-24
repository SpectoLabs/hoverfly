[![Circle CI][CircleCI-Image]][CircleCI-Url] [![Join the chat at https://gitter.im/SpectoLabs/hoverfly](https://badges.gitter.im/SpectoLabs/hoverfly.svg)](https://gitter.im/SpectoLabs/hoverfly?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![Hoverfly](core/static/img/hoverfly_logo.png)
## Dependencies without the sting

Hoverfly is a lightweight, open source [service virtualization](https://en.wikipedia.org/wiki/Service_virtualization) tool. Using Hoverfly, you can virtualize your application dependencies to create a self-contained development or test environment.

Hoverfly is a proxy written in [Go](https://golang.org/). It can capture HTTP(s) traffic between an application under test and external services, and then replace the external services. It can also generate synthetic responses on the fly.

* Capture traffic between a client application and an external service
* Use captured traffic to create a simulated service
* Export and import captured traffic
* Extend and customize with any programming language
* Dynamically create responses to requests on the fly
* Manipulate data in requests and responses
* Simulate network latency, random failures, rate limits...

## Quickstart

* [Get Hoverfly](http://hoverfly.io/#get-hoverfly)
* [Read the docs](http://hoverfly.io/introduction/)


## Contributing

Contributions are welcome!

To contribute, please:

1. Fork the repository
1. Create a feature branch on your fork
1. Commit your changes, and create a pull request against Hoverfly's master branch
1. In your pull request, include details regarding your change, i.e
    1. why you made it
    1. how to test it
    1. any information about testing you have performed

To read more about forking model, check out this link: [forking workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/forking-workflow).

## Building, running & testing

```bash
cd $GOPATH/src
mkdir -p github.com/SpectoLabs/
cd github.com/SpectoLabs/
git clone https://github.com/SpectoLabs/hoverfly.git
# or: git clone https://github.com/<your_username>/hoverfly.git
cd hoverfly
make build
```

Notice the binaries are in the ``target`` directory.

Finally to test your build:

```bash
make test
```

## License

Apache License version 2.0 [See LICENSE for details](https://github.com/SpectoLabs/hoverfly/blob/master/LICENSE).

(c) [SpectoLabs](https://specto.io) 2016.

[CircleCI-Image]: https://circleci.com/gh/SpectoLabs/hoverfly.svg?style=shield
[CircleCI-Url]: https://circleci.com/gh/SpectoLabs/hoverfly
