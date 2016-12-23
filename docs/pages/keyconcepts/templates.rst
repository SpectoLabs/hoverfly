.. _templates:

Templates
~~~~~~~~~

When it comes to matching URLs, sometimes simple matching is not enough: we want to be able to match a whole swath of URLs, e.g:

::

    http://example.com/api/v1/endpoint
    http://example.com/api/v2/endpoint
    http://example.com/api/v3/endpoint
    http://example.com/api/v4/endpoint

We may want to return the same response for all these URLs. This is done using templates. Templates use globbing, i.e they enable you to replace parts of a URL with a ``*``. So in this case, if we'd want to match all versions of the endpoint, we could have a templated URL looking like so:

::

    http://example.com/api/v*/endpoint

.. seealso::

    Templating is better understood with a practical example, so please refer to :ref:`addingtemplates` to get hands on experience with templating.
