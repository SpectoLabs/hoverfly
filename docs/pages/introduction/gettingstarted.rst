Getting Started
---------------

.. sidebar:: Note

    You are strongly recommended to keep hoverfly and hoverctl in the same directory. However if they are not in the same directory, then hoverctl first looks in the current directory for hoverfly, then in other directories in the path.


Hoverfly is composed of two binaries: hoverfly, and hoverctl.

hoverctl can be used to spawn, configure, control, and stop hoverfly. It enables to run Hoverfly as a daemon.
hoverfly is the binary that does the bulk of the work, including the proxy, and api endpoints.

Once you have extracted both hoverfly and hoverctl into the correct directories in your PATH, are able to run hoverctl and hoverfly.

.. code::

    hoverctl --version
    hoverfly --version

If installed correctly, both of these binaries should return a version number. You are able to run an instance of Hoverfly:

.. code::

    hoverctl start

We can check whether or not Hoverfly is running with the hoverctl logs command:

.. code::

    hoverctl logs

The logs should contain the string “serving proxy” which indicates that hoverfly is up and running successfully.

Finally we can stop Hoverfly:

.. code::

    hoverctl stop

