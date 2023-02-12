.. _request_matchers:

Request matchers
================

A Request Matcher is used to define the desired value for a specific request field when matching against incoming requests.
Given a **matcher value** and **string to match**, each matcher will transform and compare the values in a different way.


Exact matcher
-------------
Evaluates the equality of the matcher value and the string to match. There are no transformations.
This is the default Request Matcher type which is set by Hoverfly when requests and responses are captured.

Example
"""""""

.. code:: json

   "matcher": "exact"
   "value": "?"

.. raw:: html

    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">String to match</th>
                <th class="head">Matcher value</th>
                <th class="head">Match</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td>docs.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>specto.io</td>
                <td>docs.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

Glob matcher
------------

Allows wildcard matching (similar to BASH) using the ``*`` character.

Example
"""""""

.. code:: json

   "matcher": "glob"
   "value": "?"

.. raw:: html

    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">String to match</th>
                <th class="head">Matcher value</th>
                <th class="head">Match</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td>*.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>docs.specto.io</td>
                <td>*.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td>h*verfly.*</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>hooverfly.com</td>
                <td>h*verfly.*</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

Regex matcher
-------------
Parses the matcher value as a regular expression which is then executed against the string to match. This will pass only if the regular expression successfully
returns a result.

Example
"""""""

.. code:: json

   "matcher": "regex"
   "value": "?"

.. raw:: html

    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">String to match</th>
                <th class="head">Matcher value</th>
                <th class="head">Match</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td>(\\Ad)</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>hoverfly.io</td>
                <td>(\\Ad)</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td>(.*).(.*).(io|com|biz)</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>buy.stuff.biz</td>
                <td>(.*).(.*).(io|com|biz)</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

XML matcher
-----------
Transforms both the matcher value and string to match into XML objects and then evaluates their equality.

Example
"""""""

.. code:: json

   "matcher": "xml"
   "value": "?"

.. raw:: html

    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">String to match</th>
                <th class="head">Matcher value</th>
                <th class="head">Match</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td style="white-space:pre;">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;document type=&quot;book&quot;&gt;
        Hoverfly Documentation
    &lt;/document&gt;</td>
                <td style="white-space:pre;">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;document type=&quot;book&quot;&gt;
        Hoverfly Documentation
    &lt;/document&gt;</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td style="white-space:pre;">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;documents type=&quot;book&quot;&gt;
        &lt;document type=&quot;book&quot;&gt;
            Hoverfly Documentation
        &lt;/document&gt;
    &lt;/document&gt;</td>
                <td style="white-space:pre;">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;document type=&quot;book&quot;&gt;
        Hoverfly Documentation
    &lt;/document&gt;</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

XPath matcher
-------------
Parses the matcher value as an XPath expression, transforms the string to match into an XML object and then executes the expression against it. This will pass only if the expression successfully
returns a result.

Example
"""""""

.. code:: json

   "matcher": "xpath"
   "value": "?"

.. raw:: html

    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">String to match</th>
                <th class="head">Matcher value</th>
                <th class="head">Match</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-odd">
                <td class="example">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;documents&gt;
        &lt;document&gt;
            Hoverfly Documentation
        &lt;/document&gt;
    &lt;/documents&gt;</td>
                <td>/documents</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-even">
                <td class="example">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;document&gt;
        Hoverfly Documentation
    &lt;/document&gt;</td>
                <td>/documents</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
            <tr class="row-odd">
                <td class="example">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;documents&gt;
        &lt;document type=&quot;book&quot;&gt;
            Hoverfly Documentation
        &lt;/document&gt;
    &lt;/documents&gt;</td>
                <td>/documents/document[2]</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
            <tr class="row-odd">
                <td class="example">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;documents type=&quot;book&quot;&gt;
        &lt;document&gt;
            Someone Else's Documentation
        &lt;/document&gt;
        &lt;document&gt;
            Hoverfly Documentation
        &lt;/document&gt;
    &lt;/documents&gt;</td>
                <td>/documents/document[2]</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

JSON matcher
------------
Transforms both the matcher value and string to match into JSON objects and then evaluates their equality.

Example
"""""""

.. code:: json

   "matcher": "json"
   "value": "?"

.. raw:: html

    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">String to match</th>
                <th class="head">Matcher value</th>
                <th class="head">Match</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        },{
            "name": "Object 2",
            "set": false,
            "age": 400
        }]
    }</td>
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        },{
            "name": "Object 2",
            "set": false,
            "age": 400
        }]
    }</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        }]
    }</td>
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        },{
            "name": "Object 2",
            "set": false,
            "age": 400
        }]
    }</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

JSON partial matcher
--------------------
Unlike a JSON matcher which does the full matching of two JSON documents, this matcher evaluates if the matcher value is a subset of the incoming JSON document. The matcher ignores any absent fields and lets you match only the part of JSON document you care about.

Example
"""""""

.. code:: json

   "matcher": "jsonPartial"
   "value": "?"

.. raw:: html

    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">String to match</th>
                <th class="head">Matcher value</th>
                <th class="head">Match</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
        },{
            "name": "Object 2",
            "set": false,
            "age": 400
        }]
    }</td>
                <td class="example">{
    "objects": [
        {
            "name": "Object 1"
        },{
            "name": "Object 2"
        }]
    }</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
    <tr class="row-odd">
        <td class="example">{
    "objects": [
        {
            "name": "Object 1",
        },{
            "name": "Object 2",
            "set": false,
            "age": 400
        }]
    }</td>
                <td class="example">{
            "name": "Object 2",
            "set": false,
            "age": 400
        }</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-even">
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        }]
    }</td>
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        },{
            "name": "Object 2",
            "set": false,
            "age": 400
        }]
    }</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

JSONPath matcher
----------------
Parses the matcher value as a JSONPath expression, transforms the string to match into a JSON object and then executes
the expression against it. This will pass only if the expression successfully returns a result.


Example
"""""""

.. code:: json

   "matcher": "jsonpath"
   "value": "?"

.. raw:: html

    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">String to match</th>
                <th class="head">Matcher value</th>
                <th class="head">Match</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        }]
    }</td>
                 <td>$.objects</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td class="example">{
    "name": "Object 1",
    "set": true
    }</td>
                <td>$.objects</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
            <tr class="row-even">
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        }]
    }</td>
                <td>$.objects[1].name</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>

            <tr/>
            <tr class="row-odd">
                <td class="example">{
    "objects": [
        {
            "name": "Object 1",
            "set": true
        }, {
            "name": "Object 2",
            "set": false
        }]
    }</td>
                <td>$.objects[1].name</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
        </tbody>
    </table>


Generic Array matcher
-----------------------

Matches the matcher group with value passed in request array based on the configuration passed.

- ignoreOrder - ignore order in which values that are passed.
- ignoreUnknown - ignore unknowns in the values that are passed.
- ignoreOccurrences - ignore number of occurrences of elements that are passed with respect to matcher value.

Example
"""""""

.. code:: json

   "matcher": "array"
   "value": "[?]"
   "configuration": "{}"

   Example:

   {
        "matcher": "array",
        "config": {
            "ignoreUnknown": <true/false>,
            "ignoreOrder": <true/false>,
            "ignoreOccurrences": <true/false>
        },
        "value": [
            "access:vod",
            "order:latest",
            "profile:vd"
    ]

JWT Matcher
-----------

This matcher is primarily used for matching JWT tokens. This matcher converts base64 encoded JWT to JSON document ({"header": {}, "payload": ""}) and does partial match with the matcher value.

Matcher value contains only keys that they want to match in JWT.

Example
"""""""
    .. code:: json

        "matcher": "jwt"
        "value": "{\"header\":{\"alg\":\"HS256\"},\"payload\":{\"sub\":\"1234567890\",\"name\":\"John Doe\"}}"


Matcher Chaining
----------------

- Matcher chaining helps to chain multiple matchers. MatchedValue of parent matcher is feed into child matcher and further matching is done.

- It typically removes the stress of composing and testing complex expressions and make matchers more readable.

- It can be combine any of matchers.

For an example, with matcher chaining, one can use JSONPath to get a JSON node, and use other matcher to compare its value as mentioned below.

Example
"""""""
.. code:: json

    "matcher":<any matcher>,
    "value":?,
    "doMatch": {
        "matcher": <any matcher>
        "value":?
    }

    {
        body : [
            {
                "matcher": "jsonpath",
                "value": "$.user.id",
                "doMatch": {
                    matcher: "exact",
                    value: "1"
                }
            }
        ]

