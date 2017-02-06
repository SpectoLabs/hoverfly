.. _destination_filtering:

Destination filtering
=====================

By default, Hoverfly will process every request it receives. However, you may wish to control which URLs Hoverfly processes.

This is done by `filtering` the `destination` URLs using either a string or a regular expression. The `destination` string or regular expression will be compared against the host and the path of a URL.

For example, specifying ``hoverfly.io`` as the destination value will tell Hoverfly to process only URLs on the ``hoverfly.io`` host.

.. code:: bash

    hoverctl destination "hoverfly.io"

Specifying ``api`` as the `destination` value during :ref:`capture_mode` will tell Hoverfly to capture only URLs that contain the string ``api``. This would include both ``api.hoverfly.io/endpoint`` and ``hoverfly.io/api/endpoint``.

.. seealso::

  This functionality is best understood via a practical example: see :ref:`specific_urls` in the :ref:`tutorials` section.

.. note::

    The destination setting applies to all Hoverfly modes. If a destination value is set while Hoverfly is running in :ref:`simulate_mode`, requests that are excluded by the destination setting will be passed through to the real URLs. This makes it possible to return both real and simulated responses.
