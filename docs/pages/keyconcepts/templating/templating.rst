.. _templating:


Templating
----------

Hoverfly can build responses dynamically through templating. This is particularly useful when combined with loose matching, as it allows a single
matcher to represent an unlimited combination of responses.


Enabling Templating
~~~~~~~~~~~~~~~~~~~

By default templating is disabled. In order to enable it, set the flag to true in the response of a simulation.


Available Data
~~~~~~~~~~~~~~

Currently, the following data is available through templating:

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
| State                        | {{ State.basket }}                           | State Store = {"basket":"eggs"}              | eggs   |
+------------------------------+----------------------------------------------+----------------------------------------------+--------+

Helper Methods
~~~~~~~~~~~~~~

Additional data can come from helper methods. Current we only have some for the current data, but this list is likely to expand:

+---------------------------------------------------------+-----------------------------------------------------------+-----------------------------+
| Description                                             | Example                                                   |  Result                     |
+=========================================================+===========================================================+=============================+
| The current date, formatted in iso8601                  | {{ iso8601DateTime }}                                     |  2006-01-02T15:04:05Z07:00  |
+---------------------------------------------------------+-----------------------------------------------------------+-----------------------------+
| The current date, formatted in iso8601, with days added | {{ iso8601DateTimePlusDays Request.QueryParam.plusDays }} |  2006-02-02T15:04:05Z07:00  |
+---------------------------------------------------------+-----------------------------------------------------------+-----------------------------+

Conditional Templating, Looping and More
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly uses the https://github.com/aymerick/raymond library for templating, which is based on http://handlebarsjs.com/

To learn about more advanced templating functionality, such as looping and conditionals, read the documentation for these projects.
