.. _download_and_installation:

Download and installation
=========================

Hoverfly comes with a command line interface called **hoverctl**. Archives containing the Hoverfly and hoverctl binaries are available for the major operating systems and architectures.

- :zip_bundle_os_arch:`MacOS 64bit <OSX_amd64>`
- :zip_bundle_os_arch:`Linux 32bit <linux_386>`
- :zip_bundle_os_arch:`Linux 64bit <linux_amd64>`
- :zip_bundle_os_arch:`Windows 32bit <windows_386>`
- :zip_bundle_os_arch:`Windows 64bit <windows_amd64>`

Download the correct archive, extract the binaries and place them in a directory on your PATH.

Homebrew (MacOS)
~~~~~~~~~~~~~~~~

If you have `homebrew <http://brew.sh/>`_, you can install Hoverfly using the ``brew`` command.

.. code:: bash

    brew install SpectoLabs/tap/hoverfly

To upgrade your existing hoverfly to the latest release:

.. code:: bash

    brew upgrade hoverfly

To show which versions are installed in your machine:

.. code:: bash

    brew list --versions hoverfly

You can switch to a previously installed version as well:

.. code:: bash

    brew switch hoverfly <version>

To remove old versions :

.. code:: bash

    brew cleanup hoverfly

Docker
~~~~~~

If you have `Docker <https://www.docker.com/>`_, you can run Hoverfly using the ``docker`` command.

.. code:: bash

    docker run -d -p 8888:8888 -p 8500:8500 spectolabs/hoverfly:latest

This will run the latest version of the `Hoverfly Docker image <https://hub.docker.com/r/spectolabs/hoverfly/>`_.

You can also pass Hoverfly configuration flags when starting with Docker. For example if you need to run Hoverfly in webserver mode:

.. code:: bash

    docker run -d -p 8888:8888 -p 8500:8500 spectolabs/hoverfly:latest -webserver

This Docker image does not contain hoverctl. Our recommendation is to have hoverctl on your host machine and then 
configure hoverctl to use the newly started Hoverfly Docker instance as a new target.

.. seealso::

    For a tutorial of creating a new target in hoverctl, see :ref:`remotehoverfly`.


Kubernetes
~~~~~~~~~~

You can use `Helm <https://helm.sh/>`_ to install Hoverfly directly to your Kubernetes cluster. Hoverfly chart is available from the
official Helm incubator repo.

Use ``helm repo add`` to add the incubator repo if you haven't done so:

.. code:: bash

    helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/


Here is the command for a basic Hoverfly installation with the default settings:

.. code:: bash

    helm install incubator/hoverfly

The default installation create a ``ClusterIP`` type service, which makes it only reachable from within the cluster.

After the installation, you can use port forwarding to access the application on localhost:

.. code:: bash

    kubectl port-forward $HOVERFLY_POD_NAME 8888 8500

You can find the ``HOVERFLY_POD_NAME`` by doing ``kubectl get pod``

See more details on `Helm Charts project <https://github.com/helm/charts/tree/master/incubator/hoverfly>`_.

