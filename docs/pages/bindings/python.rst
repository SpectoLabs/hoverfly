.. _hoverpy:

HoverPy
=======

To get started:

.. code::

    sudo pip install hoverpy
    python

And in Python you can simply get started with:

.. code::

    import hoverpy
    import requests

    # capture mode
    with hoverpy.HoverPy(capture=True) as hp:
        data = requests.get("http://time.jsontest.com/").json()

    # simulation mode
    with hoverpy.HoverPy() as hp:
        simData = requests.get("http://time.jsontest.com/").json()
        print(simData)

    assert(data["milliseconds_since_epoch"] == simData["milliseconds_since_epoch"])

For more information, read the `HoverPy documentation <https://hoverpy.readthedocs.io/en/latest/>`_.
