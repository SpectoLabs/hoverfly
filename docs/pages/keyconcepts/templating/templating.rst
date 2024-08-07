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

+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| Field                        | Example                                         | Request                                      | Result         |
+==============================+=================================================+==============================================+================+
| Request scheme               | ``{{ Request.Scheme }}``                        | http://www.foo.com                           | http           |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| Query parameter value        | ``{{ Request.QueryParam.myParam }}``            | http://www.foo.com?myParam=bar               | bar            |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| Query parameter value (list) | ``{{ Request.QueryParam.NameOfParameter.[1] }}``| http://www.foo.com?myParam=bar1&myParam=bar2 | bar2           |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| Path parameter value         | ``{{ Request.Path.[1] }}``                      | http://www.foo.com/zero/one/two              | one            |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| Method                       | ``{{ Request.Method }}``                        | http://www.foo.com/zero/one/two              | GET            |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| Host                         | ``{{ Request.Host }}``                          | http://www.foo.com/zero/one/two              | www.foo.com    |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| jsonpath on body             | ``{{ Request.Body 'jsonpath' '$.id' }}``        | { "id": 123, "username": "hoverfly" }        | 123            |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| xpath on body                | ``{{ Request.Body 'xpath' '/root/id' }}``       | <root><id>123</id></root>                    | 123            |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| From data                    | ``{{ Request.FormData.email }}``                | email=foo@bar.com                            | foo@bar.com    |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| Header value                 | ``{{ Request.Header.X-Header-Id }}``            | { "X-Header-Id": ["bar"] }                   | bar            |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| Header value (list)          | ``{{ Request.Header.X-Header-Id.[1] }}``        | { "X-Header-Id": ["bar1", "bar2"] }          | bar2           |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+
| State                        | ``{{ State.basket }}``                          | State Store = {"basket":"eggs"}              | eggs           |
+------------------------------+-------------------------------------------------+----------------------------------------------+----------------+

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
| - The current date time plus 1 day in unix timestamp      | - ``{{ now '1d' 'unix' }}``                               |  - 1136300645                           |
| - The current date time in ISO 8601 format                | - ``{{ now '' '' }}``                                     |  - 2006-01-02T15:04:05Z                 |
| - The current date time minus 1 day in custom format      | - ``{{ now '-1d' '2006-Jan-02' }}``                       |  - 2006-Jan-01                          |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random string                                           | ``{{ randomString }}``                                    |  hGfclKjnmwcCds                         |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random string with a specified length                   | ``{{ randomStringLength 2 }}``                            |  KC                                     |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random boolean                                          | ``{{ randomBoolean }}``                                   |  true                                   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random integer                                          | ``{{ randomInteger }}``                                   |  42                                     |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random integer within a range                           | ``{{ randomIntegerRange 1 10 }}``                         |  7                                      |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random float                                            | ``{{ randomFloat }}``                                     |  42                                     |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random float within a range                             | ``{{ randomFloatRange 1.0 10.0 }}``                       |  7.4563213423                           |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random email address                                    | ``{{ randomEmail }}``                                     |  LoriStewart@Photolist.com              |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random IPv4  address                                    | ``{{ randomIPv4 }}``                                      |  224.36.27.8                            |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random IPv6  address                                    | ``{{ randomIPv6 }}``                                      |  41d7:daa0:6e97:6fce:411e:681:f86f:e557 |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| A random UUID                                             | ``{{ randomUuid }}``                                      |  7b791f3d-d7f4-4635-8ea1-99568d821562   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Generate random data using go-fakeit                      | ``{{ faker 'Name' }}``                                    |  John Smith                             |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Query CSV data source where ID = 3 and return its name    | ``{{csv 'test-csv' 'id' '3' 'name'}}``                    |  John Smith                             |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Query Journal index where index value = 1 and return Name | ``{{journal "Request.QueryParam.id" "1"``                 |                                         |
|  from associated Response body in journal entry.          |    ``"response" "jsonpath" "$.name"}}``                   |  John Smith                             |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+

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

Faker
~~~~~

Support for `go-fakeit <https://github.com/brianvoe/gofakeit>`_ was added in order to extend the
templating capabilities of Hoverfly. Faker covers many different test data requirements and it can be used within
Hoverfly templated responses by using the ``faker`` helper followed by the faker type (e.g. ``Name``, ``Email``)
For example, you can generate a random name using the following expression:

.. code:: json

    {
        "body": "{\"name\": \"{{faker 'Name'}}\"}"
    }

Fakers that require arguments are currently not supported.

CSV Data Source
~~~~~~~~~~~~~~~

You can query data from a CSV data source.

.. code:: json

    {
        "body": "{\"name\": \"{{csv '(data-source-name)' '(column-name)' '(query-value)' '(selected-column)' }}\"}"
    }

.. note::

    The data source name is case sensitive whereas other parameters in this function are case insensitive.
    You can use hoverctl or call the Admin API to upload CSV data source to a running Hoverfly instance.


Example: Start Hoverfly with a CSV data source (student-marks.csv) provided below.

.. code:: bash

    hoverfly -templating-data-source "student-marks <path to below CSV file>"


+-----------+-------------------+
| ID        | Name    |  Marks  |
+-----------+-------------------+
| 1         |  Test1  |    55   |
+-----------+-------------------+
| 2         |  Test2  |    65   |
+-----------+-------------------+
| 3         |  Test3  |    98   |
+-----------+-------------------+
| 4         |  Test4  |    23   |
+-----------+-------------------+
| 5         |  Test5  |    15   |
+-----------+-------------------+
| *         |  NA     |    0    |
+-----------+-------------------+

+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Description                                               | Example                                                   |  Result                                 |
+===========================================================+===========================================================+=========================================+
| Search where ID = 3 and return name                       | csv 'student-marks' 'Id' '3' 'Name'                       |  Test3                                  |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Search where ID = 4 and return its marks                  | csv 'student-marks' 'Id' '4' 'Marks'                      |  Test23                                 |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Search where Name = Test1 and return marks                | csv 'student-marks' 'Name' 'Test1' 'Marks'                |  55                                     |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Search where Id is not match and return marks             | csv 'student-marks' 'Id' 'Test100' 'Marks'                |  0                                      |
| (in this scenario, it matches wildcard * and returns)     |                                                           |                                         |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Search where Id = first path param and return marks       | csv 'student-marks' 'Id' 'Request.Path.[0]' 'Marks'       |  15                                     |
| URL looks like - http://test.com/students/5/marks         |                                                           |                                         |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+


Journal Entry Data
~~~~~~~~~~~~~~~~~~

Journal Entry can be queried using its index and its extracted value.

Syntax

.. code:: bash

    {{ journal "index name" "extracted value" "request/response" "xpath/jsonpath" "lookup query" }}


``index name`` should be the same key expression you have specified when you enable the journal index.
``extracted value`` is for doing a key lookup for the journal entry from that index.
``request/response`` specifies if you want to get data from the request or response.
``xpath/jsonpath`` specifies whether you want to extract it using xpath or json path expression.
``lookup query`` is either jsonpath or xpath expressions to parse the request/response data.

Example:

.. code:: json

    {
        "body": "{\"name\": \"{{ journal 'Request.QueryParam.id' '1' 'response' 'jsonpath' '$.name' }}\"}"
    }

In the above example, we are querying the name from JSON response in the journal entry where index ``Request.QueryParam.id`` has a key value of 1.

If you only need to check if a journal contains a particular key, you can do so using the following function:

.. code:: bash

    {{ hasJournalKey "index name" "key name" }}


Key Value Data Store
~~~~~~~~~~~~~~~~~~~~

Sometimes you may need to store a temporary variable and retrieve it later in other part of the templated response.
In this case, you can use the internal key value data store. The following helper methods are available:

+----------------------------+--------------------------------------------+-----------------------+
| Description                | Example                                    |  Result               |
+============================+============================================+=======================+
| Put an entry               | ``{{ putValue 'id' 123 true }}``           |  123                  |
+----------------------------+--------------------------------------------+-----------------------+
| Get an entry               | ``{{ getValue 'id' }}``                    |  123                  |
+----------------------------+--------------------------------------------+-----------------------+
| Add a value to an arra     | ``{{ addToArray 'names' 'John' true }}``   |  John                 |
+----------------------------+--------------------------------------------+-----------------------+
| Get an array               | ``{{ getArray 'names' }}``                 |  []string{"John"      |
+----------------------------+--------------------------------------------+-----------------------+

``addToArray`` will create a new array if one doesn't exist. The boolean argument in ``putValue`` and ``addToArray``
is used to control whether the set value is returned.

.. note::

    Each templating session has its own key value store, which means all the data you set will be cleared after the current response is rendered.


Maths Operations
~~~~~~~~~~~~~~~~

The basic maths operations are currently supported: add, subtract, multiply and divide. These functions
take three parameters: two values it operates on and the precision. The precision is given in a string
format such as ``'0.00'``. For example ``{{ add 3 2.5 '0.00' }}`` should give you ``5.50``.
If no format is given, the exact value will be printed with up to 6 decimal places.

+------------+---------------------------------+---------------+
| Description| Example                         |  Result       |
+============+=================================+===============+
| Add        | ``{{ add 10 3 '0.00' }}``       |  13.33        |
+------------+---------------------------------+---------------+
| Subtract   | ``{{ subtract 10 3 '' }}``      |  7            |
+------------+---------------------------------+---------------+
| Multiply   | ``{{ multiply 10 3 '' }}``      |  30           |
+------------+---------------------------------+---------------+
| Divide     | ``{{ divide 10 3 '' }}``        |  3.333333     |
+------------+---------------------------------+---------------+

A math functions for summing an array of numbers is also supported; it's usually used in conjunction
with the ``#each`` block helper. For example:

With the request payload of

.. code:: json

    {
        "lineitems": {
            "lineitem": [
                {
                    "upc": "1001",
                    "quantity": "1",
                    "price": "3.50"
                },
                {
                    "upc": "1002",
                    "quantity": "2",
                    "price": "4.50"
                }
            ]
        }
    }

We can get the total price of all the line items using this templating function:


``{{#each (Request.Body 'jsonpath' '$.lineitems.lineitem') }}``
``{{ addToArray 'subtotal' (multiply (this.price) (this.quantity) '') false }} {{/each}}``
``total: {{ sum (getArray 'subtotal') '0.00' }}``

String Operations
~~~~~~~~~~~~~~~~~

You can use the following helper methods to join, split or replace string values.

+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Description                                               | Example                                                   |  Result                                 |
+===========================================================+===========================================================+=========================================+
| String concatenate                                        | ``{{ concat 'bee' 'hive' }}``                             |  beehive                                |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| String splitting                                          | ``{{ split 'bee,hive' ',' }}``                            |  []string{"bee", "hive"}                |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Replace all occurrences of the old value with the new     | ``{{ replace (Request.Body 'jsonpath' '$.text')``         |                                         |
|                                                           |    ``'be' 'mock' }}``                                     |                                         |
| value in the target string                                | (where Request.Body has the value of                      |                                         |
|                                                           |                                                           |                                         |
|                                                           | ``{"text":"to be or not to be"}``                         |  to mock or not to mock                 |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Return a substring of a string                            | ``{{substring 'thisisalongstring' 7 11}}``                |  long                                   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Return the length of a string                             | ``{{length 'thisisaverylongstring'}}``                    |  21                                     |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Return the rightmost characters of a string               | ``{{rightmostCharacters 'thisisalongstring' 3}}``         |  ing                                    |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+

Validation Operations
~~~~~~~~~~~~~~~~~~~~~

You can use the following helper methods to validate various types, compare value, and perform regular expression matching on strings.

+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Description                                               | Example                                                   |  Result                                 |
+===========================================================+===========================================================+=========================================+
| Is the value numeric                                      | ``{{isNumeric '12.3'}}``                                  |  true                                   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Is the value alphanumeric                                 | ``{{isAlphanumeric 'abc!@123'}}``                         |  false                                  |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Is the value a boolean                                    | ``{{isBool (Request.Body 'jsonpath' '$.paidInFull')}}``   |  true                                   |
|                                                           |  Where the payload is {"paidInFull":"false"}              |                                         |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Is one value greater than another                         | ``{{isGreater (Request.Body 'jsonpath' '$.age') 25}``     |  false                                  |
|                                                           |  Where the payload is {"age":"19"}                        |                                         |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Is one value less than another                            | ``{{isLess (Request.Body 'jsonpath' '$.age') 25}``        |  true                                   |
|                                                           |  Where the payload is {"age":"19"}                        |                                         |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Is a value between two values                             | ``{{isBetween (Request.Body 'jsonpath' '$.age') 25 35}``  |  false                                  |
|                                                           |  Where the payload is {"age":"19"}                        |                                         |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+
| Does a string match a regular expression                  | ``{{matchesRegex '2022-09-27' '^\d{4}-\d{2}-\d{2}$'}}``   |  true                                   |
+-----------------------------------------------------------+-----------------------------------------------------------+-----------------------------------------+

Conditional Templating, Looping and More
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly uses the https://github.com/aymerick/raymond library for templating, which is based on http://handlebarsjs.com/

To learn about more advanced templating functionality, such as looping and conditionals, read the documentation for these projects.

Global Literals and Variables
-----------------------------
You can define global literals and variables for templated response. This comes in handy when you
have a lot of templated responses that share the same constant values or helper methods.

Literals
~~~~~~~~

Literals are constant values. You can declare literals as follows and then reference it in templated response as ``{{ Literals.<literal name> }}``.

::

    {
      "data": {
      ...
      "literals": [
            {
                "name":"literal1",
                "value":"value1"
            },
            {
                "name":"literal2",
                "value":["value1", "value2", "value3"]
            },
            {
                "name":"literal3",
                "value": {
                    "key": "value"
                }
            }
        ]
    }


Variables
~~~~~~~~~

Variable lets you define a helper method that can be shared among templated responses.
You can associate the helper method with a name and then reference it in templated response as ``{{ Vars.<variable name> }}``.

::

    {
      "data": {
      ...
      "variables": [
            {
                "name":"<variable name>",
                "function":"<helper method name>",
                "arguments":["arg1", "arg2"]

            }
        ]
    }

    {
      "data": {
      ...
      "variables": [
            {
                "name":"varOne",
                "function":"faker",
                "arguments":["Name"]

            },
            {
                "name":"idFromJSONRequestBody",
                "function":"requestBody",
                "arguments":["jsonpath", "$.id"]
            },
            {
                "name":"idFromXMLRequestBody",
                "function":"requestBody",
                "arguments":["xpath", "/root/id"]
            }
        ]
    }

