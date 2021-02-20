.. _templating:


Templating
==========

Hoverfly can build responses dynamically through templating. This is particularly useful when combined with loose matching, as it allows a single
matcher to represent an unlimited combination of responses.


Enabling Templating
-------------------

By default templating is disabled. In order to enable it, set the ``templated`` field to true in the response of a simulation.


Getting data from the request
-----------------------------

Currently, you can get the following data from request to the response via templating:

+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| Field                        | Example                                      | Request                                      | Result |
+==============================+==============================================+==============================================+========+
| Request scheme               | {{ Request.Scheme }}                         | http://www.foo.com                           | http   |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| Query parameter value        | {{ Request.QueryParam.myParam }}             | http://www.foo.com?myParam=bar               | bar    |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| Query parameter value (list) | {{ Request.QueryParam.NameOfParameter.[1] }} | http://www.foo.com?myParam=bar1&myParam=bar2 | bar2   |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| Path parameter value         | {{ Request.Path.[1] }}                       | http://www.foo.com/zero/one/two              | one    |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| Method                       | {{ Request.Method }}                         | http://www.foo.com/zero/one/two              | GET    |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| jsonpath on body             | {{ Request.Body 'jsonpath' '$.id' }}         | { "id": 123, "username": "hoverfly" }        | 123    |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| xpath on body                | {{ Request.Body 'xpath' '/root/id' }}        | <root><id>123</id></root>                    | 123    |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| Header value                 | {{ Request.Header.X-Header-Id }}             | { "X-Header-Id": ["bar"] }                   | bar    |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| Header value (list)          | {{ Request.Header.X-Header-Id.[1] }}         | { "X-Header-Id": ["bar1", "bar2"] }          | bar2   |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| State                        | {{ State.basket }}                           | State Store = {"basket":"eggs"}              | eggs   |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+

Helper Methods
--------------

Additional data can come from helper methods. These are the ones Hoverfly currently support:

+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Description                                               | Example                                                   |  Result                                 |
+===========================================================+===========================================================+=========================================+
| The current date time with offset, in the given format.   |                                                           |                                         |
|                                                           |                                                           |                                         |
| For example:                                              |                                                           |                                         |
|                                                           |                                                           |                                         |
| - The current date time plus 1 day in unix timestamp      | - {{ now '1d' 'unix' }}                                   |  - 1136300645                           |
| - The current date time in ISO 8601 format                | - {{ now '' '' }}                                         |  - 2006-01-02T15:04:05Z                 |
| - The current date time minus 1 day in custom format      | - {{ now '-1d' '2006-Jan-02' }}                           |  - 2006-Jan-01                          |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random string                                           | {{ randomString }}                                        |  hGfclKjnmwcCds                         |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random string with a specified length                   | {{ randomStringLength 2 }}                                |  KC                                     |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random boolean                                          | {{ randomBoolean }}                                       |  true                                   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random integer                                          | {{ randomInteger }}                                       |  42                                     |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random integer within a range                           | {{ randomIntegerRange 1 10 }}                             |  7                                      |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random float                                            | {{ randomFloat }}                                         |  42                                     |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random float within a range                             | {{ randomFloatRange 1.0 10.0 }}                           |  7.4563213423                           |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random email address                                    | {{ randomEmail }}                                         |  LoriStewart@Photolist.com              |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random IPv4  address                                    | {{ randomIPv4 }}                                          |  224.36.27.8                            |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random IPv6  address                                    | {{ randomIPv6 }}                                          |  41d7:daa0:6e97:6fce:411e:681:f86f:e557 |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random UUID                                             | {{ randomUuid }}                                          |  7b791f3d-d7f4-4635-8ea1-99568d821562   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Replace all occurrences of the old value with the new     | {{ replace Request.Body 'be' 'mock' }}                    |                                         |
|                                                           |                                                           |                                         |
| value in the target string                                | (where Request.Body has the value of "to be or not to be" |  to mock or not to mock                 |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+

.. note::

    The following helper methods will be deprecated, please use ``now`` helper for date time:
        - ``iso8601DateTime``
        - ``iso8601DateTimePlusDays``
        - ``currentDateTime``
        - ``currentDateTimeAdd``
        - ``currentDateTimeSubtract``


Time offset
~~~~~~~~~~~
When using template helper method ``now``, time offset must be formatted using the following syntax.

+-----------+-------------+
| Shorthand | Type        |
+===========+=============+
| ns        | Nanosecond  |
+-----------+-------------+
| us/µs     | Microsecond |
+-----------+-------------+
| ms        | Millisecond |
+-----------+-------------+
| s         | Second      |
+-----------+-------------+
| m         | Minute      |
+-----------+-------------+
| h         | Hour        |
+-----------+-------------+
| d         | Day         |
+-----------+-------------+
| y         | Year        |
+-----------+-------------+

Prefix an offset with ``-`` to subtract the duration from the current date time.

Example time offset
~~~~~~~~~~~~~~~~~~~

+-----------+-------------------+
| 5m        | 5 minutes         |
+-----------+-------------------+
| 1h30m     | 1 hour 5 minutes  |
+-----------+-------------------+
| 1y10d     | 1 year 10 days    |
+-----------+-------------------+

Date time formats
~~~~~~~~~~~~~~~~~
When using template helper method ``now``, date time formats must follow the Golang syntax.
More can be found out here https://golang.org/pkg/time/#Parse

Example date time formats
~~~~~~~~~~~~~~~~~~~~~~~~~

+-------------------------------+
| 2006-01-02T15:04:05Z07:00     |
+-------------------------------+
| Mon, 02 Jan 2006 15:04:05 MST |
+-------------------------------+
| Jan _2 15:04:05               |
+-------------------------------+

.. note::

    If you leave the format string empty, the default format to be used is ISO 8601 (2006-01-02T15:04:05Z07:00).

    You can also get an UNIX timestamp by setting the format to:

    - ``unix``: UNIX timestamp in seconds
    - ``epoch``: UNIX timestamp in milliseconds


Conditional Templating, Looping and More
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly uses the https://github.com/aymerick/raymond library for templating, which is based on http://handlebarsjs.com/

To learn about more advanced templating functionality, such as looping and conditionals, read the documentation for these projects.
