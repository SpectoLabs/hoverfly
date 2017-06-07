.. _simulate_mode:

Simulate mode
=============

In this mode, Hoverfly uses its simulation data in order to simulate external APIs. Each time Hoverfly receives a request,
rather than forwarding it on to the real API, it will respond instead. No network traffic will ever reach the real external API.

.. figure:: simulate.mermaid.png

The simulation can be produced automatically via by running Hoverfly in capture mode, or created manually. See :ref:`simulations` for information.

Matching
--------

In order for Hoverfly to determine which response it should give for a request, it makes use of a process called “matching”.
Essentially, matching involves comparing the request to each matcher in the simulation until it finds one which matches.
Once a matching matcher has been found, Hoverfly will return the response associated with that matcher. A more detailed
explanation of how a match takes place can be found in :ref:`matchers`.

Strongest Match
~~~~~~~~~~~~~~~

Strongest match is the default matching strategy for Hoverfly. This means that if there are multiple matching matchers
for a request, the one with the highest matching score will be used. A matching score is calculated by adding together
the total amount of matches in a matcher.

To set "strongest" as the matching strategy, simply run:

.. code:: bash

    hoverctl mode simulate

Or to be explicit run

.. code:: bash

    hoverctl mode simulate --matching-strategy=strongest


As an example, let's run a request against some matchers and see what happens:

**Request:**

+-------------+---------------------+
| Method      | Destination         |
+=============+=====================+
| GET         | www.destination.com |
+-------------+---------------------+

**Matchers:**

+-------------+--------------+-------------------------+-----------+-------------+----------+
| Field       | Matcher Type | Value                   | Score     | Total Score | Matched? |
+=============+==============+=========================+===========+=============+==========+
| method      | exactMatch   | DELETE                  | +0        | 1           | false    |
+-------------+--------------+-------------------------+-----------+             +          |
| destination | exactMatch   | www.destination.com     | +1        |             |          |
+-------------+--------------+-------------------------+-----------+-------------+----------+

+-------------+--------------+-------------------------+-----------+-------------+----------+
| Field       | Matcher Type | Value                   | Score     | Total Score | Matched? |
+=============+==============+=========================+===========+=============+==========+
| method      | exactMatch   | GET                     | +1        | 1           | true     |
+-------------+--------------+-------------------------+-----------+-------------+----------+

+-------------+--------------+-------------------------+-----------+-------------+----------+
| Field       | Matcher Type | Value                   | Score     | Total Score | Matched? |
+=============+==============+=========================+===========+=============+==========+
| method      | exactMatch   | GET                     | +1        | 2           | true     |
+-------------+--------------+-------------------------+-----------+             +          |
| destination | exactMatch   | www.destination.com     | +1        |             |          |
+-------------+--------------+-------------------------+-----------+-------------+----------+

+-------------+--------------+-------------------------+-----------+-------------+----------+
| Field       | Matcher Type | Value                   | Score     | Total Score | Matched? |
+=============+==============+=========================+===========+=============+==========+
| method      | exactMatch   | GET                     | +1        | 1           | false    |
+-------------+--------------+-------------------------+-----------+             +          |
| destination | exactMatch   | www.miss.com            | +0        |             |          |
+-------------+--------------+-------------------------+-----------+-------------+----------+

The first matcher and last matchers are missed, so are not candidates, but the remaining two are. Hoverfly picks the
third matcher as it has the highest score of two. In the case that there is more than one match with the highest score,
the Hoverfly will pick the first one.

Another advantage of this strategy is knowing which miss is the closest when there are no matches. This is because all
matchers are scored, even if they do not match. This is what allows Hoverfly to return the closest matcher if there is
a miss. For more information see :ref:`_troubleshooting`.

The main disadvantage of this strategy is performance, as Hoverfly has to iterate through all matchers each time.
However, much of this is likely to be negated by ordinary and eager caching. See caching for more information.

First Match
~~~~~~~~~~~

First match is the alternative (legacy) mechanism of matching. In this case, there is no scoring of matchers - once a match
is found it is returned, even if there would have been more matches later on in the simulation.

To set first match as the matching strategy, simply run:

.. code:: bash

    hoverctl mode simulate --matching-strategy=first

The main advantage of this is performance, as there is no logic in calculating match scores, and Hoverfly does not iterate
through each matcher in the simulation unless there is a miss. It does make debugging harder, as Hoverfly cannot tell you
what the closest matcher was in the event of a miss. It also means that the order of the pairs in the simulation matters -
ones at the end will not be matched is there is an earlier match.

Caching
-------

During simulate, mode Hoverfly makes use of caching in order to retain strong performance characteristics, even in the
event of complex matching. The cache is a key-value store of request hashes (hash of all the request excluding headers) to responses.

Caching Matches
~~~~~~~~~~~~~~~

When Hoverfly receives a request in simulate mode, it will first hash it and then look for it in the cache. If a cache entry
is found, it will send the cached response back to Hoverfly. If it is not found, it will look for a match in the list of
matchers. Whenever a new match is found, it will be added to the cache. What this means is that whenever a request is repeated
Hoverflies cache will be used, bring up performance to O(1).

Caches Misses
~~~~~~~~~~~~~

Hoverfly also caches misses. This means that repeating a request which was not matched will return a cached miss, avoiding
the need to perform matching. The closest miss is also cached, so Hoverfly will not lose any useful information about which
matchers came closest.

Header caching
~~~~~~~~~~~~~~

Right now, headers are not included in the hash for a request. This is because headers tend to change across time and clients,
impacting ordinary and eager caching respectively. It means if a matcher takes place on headers, it cannot only be partially cached.

Eager caching
~~~~~~~~~~~~~

The cache is automatically pre-populated whenever switching to simulate mode. This only works on certain matchers (such a matchers
where every field is an “exactMatch”), but means the initial cache population does not happen during a simulation.

Cache Invalidation
~~~~~~~~~~~~~~~~~~

Cache invalidation is a straightforward process in Hoverfly. It only occurs when a simulation is modified.