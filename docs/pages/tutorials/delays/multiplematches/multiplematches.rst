Delays and multiple matches
...........................

You can easily get into a situation where your request URL has multiple matches within your delays.json file. In this case, the first successful match that is met in the json file wins.

.. literalinclude:: delays.json
   :language: json

.. literalinclude:: delays.sh
   :language: json

Alternatively, the configuration can be applied using the Hoverfly API directly:

To apply delays:

::

    curl -H "Content-Type application/json" -X PUT -d '{"data":[{"urlPattern":"1\\.myhost\\.com","delay":1000},{"urlPattern":"2\\.myhost\\.com","delay":2000}]}' http://${HOVERFLY_HOST}:8888/api/delays

To view the delays which have been applied

::

    curl http://${HOVERFLY_HOST}:8888/api/delays