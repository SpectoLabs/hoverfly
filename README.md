![Hoverfly](https://raw.githubusercontent.com/SpectoLabs/hoverfly/master/hoverfly_logo.png)

[![Circle CI][CircleCI-Image]][CircleCI-Url] 
[![Documentation Status](https://readthedocs.org/projects/hoverfly/badge/?version=latest)](http://hoverfly.readthedocs.io/en/latest/?badge=latest)
[![Join the chat at https://gitter.im/SpectoLabs/hoverfly](https://badges.gitter.im/SpectoLabs/hoverfly.svg)](https://gitter.im/SpectoLabs/hoverfly?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## API simulations for development and testing

Hoverfly is a lightweight, open source API simulation tool. Using Hoverfly, you can create realistic simulations of the APIs your application depends on.

* Replace slow, flaky API dependencies with realistic, re-usable simulations
* Simulate network latency, random failures or rate limits to test edge-cases
* Extend and customize with any programming language
* Export, share, edit and import API simulations
* CLI and native language bindings for [Java](https://hoverfly-java.readthedocs.io/en/latest/)
* REST API
* Lightweight, high-performance, run anywhere
* Apache 2 license

Hoverfly is developed and maintained by [SpectoLabs](https://specto.io).

## Quickstart

* [Download and installation](https://hoverfly.readthedocs.io/en/latest/pages/introduction/downloadinstallation.html)
* [Read the docs](https://hoverfly.readthedocs.io)
* [Join the mailing list](https://groups.google.com/a/specto.io/forum/#!forum/hoverfly)

## Hoverfly Cloud

Hoverfly Cloud - API simulations as a service. [Sign up here.](https://hoverfly.io/hoverfly-cloud)

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


## Setup for local development

### Build
1. Install Go. You should follow the download and installation instructions for your OS from [official Go website](https://golang.org/dl). If you have installed Go before 
using either `apt-get` or `homebrew`, you should uninstall those first.
1. Clone the hoverfly project:
    ```bash
    git clone https://github.com/SpectoLabs/hoverfly.git
    # or: git clone https://github.com/<your_username>/hoverfly.git
    ```
1. Finally you can run the build

    ```bash
    cd hoverfly
    make build 
    ```

Notice the binaries are in the ``target`` directory.


### Test

You should be able to execute all unit and functional tests with: 

```bash
make test
```

Some middleware tests may fail if you don't have `ruby` and `python` setup in your environment. If you are using Mac, you can install them with [Homebrew](https://brew.sh/):
```bash
brew install ruby
brew install python
``` 


## License

Apache License version 2.0 [See LICENSE for details](https://github.com/SpectoLabs/hoverfly/blob/master/LICENSE).

(c) [SpectoLabs](https://specto.io) 2017.

[CircleCI-Image]: https://circleci.com/gh/SpectoLabs/hoverfly.svg?style=shield
[CircleCI-Url]: https://circleci.com/gh/SpectoLabs/hoverfly
