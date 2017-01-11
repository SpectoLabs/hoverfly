.. _templates:

Templates
*********

Sometimes simple one-to-one matching of responses to requests is not enough. For example, you may want to return a single response for all these URLs:

::

    http://example.com/api/v1/endpoint
    http://example.com/api/v2/endpoint
    http://example.com/api/v3/endpoint
    http://example.com/api/v4/endpoint

This is done using templates. Templates use globbing: i.e they allow you to replace parts of a URL with a ``*``. In this case, if we want to match all versions of the endpoint, we could specify a URL like this:

::

    http://example.com/api/v*/endpoint

.. seealso::

    Templating is best understood with a practical example, so please refer to :ref:`addingtemplates` to get hands on experience with templating.
