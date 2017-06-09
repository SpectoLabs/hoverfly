.. _caching:

Caching
-------

In :ref:`simulate_mode`, Hoverfly uses caching in order to retain strong performance characteristics, even in the
event of complex matching. The cache is a key-value store of request hashes (hash of all the request excluding headers) to responses.

Caching matches
~~~~~~~~~~~~~~~

When Hoverfly receives a request in simulate mode, it will first hash it and then look for it in the cache. If a cache entry
is found, it will send the cached response back to Hoverfly. If it is not found, it will look for a match in the list of
matchers. Whenever a new match is found, it will be added to the cache.

Caches misses
~~~~~~~~~~~~~

Hoverfly also caches misses. This means that repeating a request which was not matched will return a cached miss, avoiding
the need to perform matching. The closest miss is also cached, so Hoverfly will not lose any useful information about which
matchers came closest.

Header caching
~~~~~~~~~~~~~~

Currently, headers are not included in the hash for a request. This is because headers tend to change across time and clients,
impacting ordinary and eager caching respectively.

Eager caching
~~~~~~~~~~~~~

The cache is automatically pre-populated whenever switching to simulate mode. This only works on certain matchers (such as matchers
where every field is an “exactMatch”), but it means the initial cache population does not happen **during** simulate mode.

Cache invalidation
~~~~~~~~~~~~~~~~~~

Cache invalidation is a straightforward process in Hoverfly. It only occurs when a simulation is modified.