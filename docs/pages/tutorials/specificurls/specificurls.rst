.. _specific_urls:

Capturing or simulating specific URLs
=====================================

We can use the ``hoverctl destination`` command to specify which URLs to capture or simulate. The destination 
setting can be tested using the ``--dry-run`` flag. This makes it easy to check whether the destination setting 
will filter out URLs as required.

.. literalinclude:: filter-dest-dry-run.sh
    :language: sh

This tells us that setting the destination to ``ip`` will allow the URL ``http://ip.jsontest.com`` to be captured or simulated, while the URL
``http://time.jsontest.com`` will be ignored.


Now we have checked the destination setting, we can apply it to filter out the URLs we don't want to capture.

.. literalinclude:: filter-dest.sh
    :language: sh

If we examine the logs and the ``simulation.json`` file, we can see that `only` a request response pair to the ``http://ip.jsontest.com``
URL has been captured.

The destination setting can be either a string or a regular expression. 

.. literalinclude:: filter-dest-dry-run-regex.sh
    :language: sh

Here, we can see that setting the destination to ``^.*api.*com`` will allow the ``https://api.github.com`` and ``https://api.slack.com`` 
URLs to be captured or simulated, while the ``https://github.com`` URL will be ignored.