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
| jsonpath on body             | {{ Request.Body "jsonpath" "$.test" }}       | { "id": 123, "username": "hoverfly" }        | 123    |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| xpath on body                | {{ Request.Body "xpath" "/root/id" }}        | <root><id>123</id></root>                    | 123    |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+
| State                        | {{ State.basket }}                           | State Store = {"basket":"eggs"}              | eggs   |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+

Helper Methods
--------------

Additional data can come from helper methods. Current we only have some for the current data, but this list is likely to expand:

+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Description                                               | Example                                                   |  Result                                 |
+===========================================================+===========================================================+=========================================+
| The current UTC date time, formatted in iso8601           | {{ iso8601DateTime }}                                     |  2006-01-02T15:04:05Z                   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| The current UTC date time, formatted in iso8601,          |                                                           |                                         |
| with days added                                           | {{ iso8601DateTimePlusDays Request.QueryParam.plusDays }} |  2006-02-02T15:04:05Z                   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| The current UTC date time, in the format specified        | {{ currentDateTime "2006-Jan-02" }}                       |  2018-Jul-05                            |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| The current UTC date time, in the format specified,       |                                                           |                                         |
| with duration added                                       | {{ currentDateTimeAdd "1d" "2006-Jan-02" }}               |  2018-Jul-06                            |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| The current UTC date time, in the format specified,       |                                                           |                                         |
| with duration subtracted                                  | {{ currentDateTimeSubtract "1d" "2006-Jan-02" }}          |  2018-Jul-04                            |
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

Durations
~~~~~~~~~
When using template helper methods such as ``currentDateTimeAdd`` and ``currentDateTimeSubtract``, durations must be formatted following the following syntax for durations. 

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

Example Durations
~~~~~~~~~~~~~~~~~

+-----------+-------------------+
| 5m        | 5 minutes         |
+-----------+-------------------+
| 1h30m     | 1 hour 5 minutes  |
+-----------+-------------------+
| 1y10d     | 1 year 10 days    |
+-----------+-------------------+

Date time formats
~~~~~~~~~~~~~~~~~
When using template helper methods such as ``currentDateTime``, ``currentDateTimeAdd`` and ``currentDateTimeSubtract``, date time formats must follow
the Golang syntax. More can be found out here https://golang.org/pkg/time/#Parse

Example date time formats
~~~~~~~~~~~~~~~~~~~~~~~~~

+-------------------------------+
| 2006-01-02T15:04:05Z07:00     |
+-------------------------------+
| Mon, 02 Jan 2006 15:04:05 MST |
+-------------------------------+
| Jan _2 15:04:05               |
+-------------------------------+


Conditional Templating, Looping and More
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly uses the https://github.com/aymerick/raymond library for templating, which is based on http://handlebarsjs.com/

To learn about more advanced templating functionality, such as looping and conditionals, read the documentation for these projects.
