.. _getting_started:

Getting Started
===============

.. sidebar:: Note

    It is recommended that you keep Hoverfly and hoverctl in the same directory. However if they are not in the same directory, hoverctl will look in the current directory for Hoverfly, then in other directories on the PATH.


Hoverfly is composed of two binaries: **Hoverfly** and **hoverctl**.

hoverctl is a command line tool that can be used to configure and control Hoverfly. It allows you to run Hoverfly as a daemon.

Hoverfly is the application that does the bulk of the work. It provides the proxy server or webserver, and the API endpoints.

Once you have extracted both Hoverfly and hoverctl into a directory on your PATH, you can run hoverctl and Hoverfly.

.. code:: bash

    hoverctl version
    hoverfly -version

Both of these commands should return a version number. Now you can run an instance of Hoverfly:

.. code:: bash

    hoverctl start

Check whether Hoverfly is running with the following command:

.. code:: bash

    hoverctl logs

The logs should contain the string ``serving proxy``. This indicates that Hoverfly is running.

Finally, stop Hoverfly with:

.. code:: bash

    hoverctl stop
