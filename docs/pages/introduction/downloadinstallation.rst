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

    brew list --version hoverfly

You can switch to a previously installed version as well:

.. code:: bash

    brew switch hoverfly <version>

To remove old versions :

.. code:: bash

    brew cleanup hoverfly

