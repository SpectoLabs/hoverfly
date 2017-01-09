.. image:: logo-large.png

|
|

What is Hoverfly?
=================

Hoverfly is a lightweight, open source tool for creating simulations of HTTP(S) APIs for use in development and testing. This technique is sometimes referred to as `service virtualization <https://en.wikipedia.org/wiki/Service_virtualization>`_.

Hoverfly allows you to:

- Capture traffic between your application and an external API to create a simulation of the API
- Export the simulation to JSON, edit it and re-import it
- Create a simulation from scratch by writing a JSON file
- Simulate unexpected API behaviour, such as random latency or failure

Features:

- CLI and native language bindings for Java and Python
- Single binary file with no dependencies (optional CLI is a separate binary file)
- Extend and customize behaviour with any language
- High performance
- Official Docker image
- REST API
- Apache 2.0 license

Hoverfly is developed and maintained by `SpectoLabs <https://specto.io>`_.

.. toctree::
    :maxdepth: 2

    pages/introduction/introduction
    pages/keyconcepts/keyconcepts
    pages/bindings/bindings
    pages/tutorials/tutorials
    pages/reference/reference
