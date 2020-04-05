.. _contributing:

Contributing
============

Contributions are welcome! To contribute, please:

1. Fork the repository
2. Create a feature branch on your fork
3. Commit your changes, and create a pull request against Hoverfly's master branch

In your pull request, please include details regarding your change (why you made the change, how to test it etc).

Learn more about the `forking workflow here <https://www.atlassian.com/git/tutorials/comparing-workflows/forking-workflow>`_.

Building, running & testing
---------------------------

You will need `Go 1.14 <https://golang.org>`_ . Instructions on how to set up your Go environment can be `found here <https://golang.org/doc/install>`_.

.. code:: bash

    cd $GOPATH/src
    mkdir -p github.com/SpectoLabs/
    cd github.com/SpectoLabs/
    git clone https://github.com/SpectoLabs/hoverfly.git
    # or: git clone https://github.com/<your_username>/hoverfly.git
    cd hoverfly
    make build


Notice the binaries are in the ``target`` directory.

Finally, to test your build:

.. code:: bash

    make test
