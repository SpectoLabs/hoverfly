.. _webserver:

Hoverfly as a Webserver
-----------------------

Sometimes you may not be able to configure your client to use a proxy, or you may simply want to point to Hoverfly as a webserver.

.. figure:: webserver.mermaid.png

.. note::
    
    In this mode, Hoverfly does not capture traffic, it can only replay it. For this reason, when you use Hoverfly as a webserver, you should have simulations ready to be loaded. 

.. In webserver mode, Hoverfly strips the domain from the endpoint's URL. So for example, if during capture phase you made requests to:

..    .. code::
        
        http://echo.jsontest.com/key/value

..    And Hoverfly is running on:

..    .. code::
        
        http://localhost:8888

..    Then the URL that would retrieve the data back would be:

..    .. code::
        
        http://localhost:8500/key/value

.. seealso::
    
    Please refer to the :ref:`webservertutorial` tutorial, for a step by step example.