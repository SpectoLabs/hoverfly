Python
------

HoverPy enables you to quickly and easily get started:

.. code::

    sudo pip install hoverpy
    python

And in python you can simply get started with:

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

You can read hoverpy's full documentation at http://hoverpy.io.