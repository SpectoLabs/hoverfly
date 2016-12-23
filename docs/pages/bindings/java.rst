Java
----

Simply add the following rule to your class, giving it the path to your Simulation on the ClassPath.

.. code:: java

    public HoverflyRule hoverflyRule = HoverflyRule.inSimulationMode("test-service.json");

As Hoverfly is a proxy, all requests (assuming they respect the JVM proxy settings) will be routed through it. Hoverfly will then respond instead of the real service.

The rule will attempt to detect the operating system and architecture type of the host, and then extract and execute the correct hoverfly binary. It will import the json into it’s database and then and destroy the process at the end of the tests.

To read more about Hoverfly's junit binding, please go to https://github.com/SpectoLabs/hoverfly-junit

