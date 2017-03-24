.. request_matchers:

Request Matchers
================


A Request Matcher is used to define the desired value for a specific request field when matching against incoming requests. Given a matcher value and provided string to match when matching, each matcher will transform the values and compare differently.

|
.. todo:: Check examples and name them... maybe

Exact Matcher
-------------
This request matcher will match on the equality of the matcher value and the string to match. There are no transformations. This is the default request matcher used when capturing with Hoverfly. 

Example
"""""""

.. code:: json
   
   "exactMatch": "docs.hoverfly.io"


.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>specto.io</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

Glob Matcher
------------
This request matcher will match on the equality of the matcher value and the string to match but will allow wildcard matching similar to BASH with the ``*`` character.

Example 1
"""""""""

.. code:: json
   
   "globMatch": "*.hoverfly.io"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>docs.specto.io</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

Example 2
"""""""""

.. code:: json
   
   "globMatch": "h*verfly.*"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>hooverfly.com</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
        </tbody>
    </table>

.. todo:: Buy hooverfly.com?

|
|

Regex Matcher
------------
This request matcher will parse the matcher value as a regular expression. It will execute the expression against the string to match. This will pass only if the expression successfully returns a result.

Example 1
"""""""""

.. code:: json
   
   "regexMatch": "(\\Ad)"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

Example 2
"""""""""

.. code:: json
   
   "regexMatch": "(.*).(.*).(io|com|biz)"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td>docs.hoverfly.io</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td>buy.stuff.biz</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

XML Matcher
-----------
This request matcher will transform both matcher value and string to match as XML objects and then match on the equality of those objects.

Example
"""""""

.. code:: json
   
   "xmlMatch": "<?xml version="1.0" encoding="UTF-8"?><document type="book">Hoverfly Documentation</document>"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
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
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

XPath Matcher
------------
This request matcher will parse the matcher value as an XPath expression. It will transform the string to match into an XML object and then execute the expression against it. This will pass only if the expression successfully returns a result.

Example 1
"""""""""

.. code:: json
   
   "xpathMatch": "/documents"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-odd">
                <td style="white-space:pre;">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;documents&gt;
        &lt;document&gt;
            Hoverfly Documentation
        &lt;/document&gt;
    &lt;/documents&gt;</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-even">
                <td style="white-space:pre;">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;document&gt;
        Hoverfly Documentation
    &lt;/document&gt;</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

Example 2
"""""""""

.. code:: json
   
   "xpathMatch": "/documents[2]"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-odd">
                <td style="white-space:pre;">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;documents&gt;
        &lt;document type=&quot;book&quot;&gt;
            Hoverfly Documentation
        &lt;/document&gt;
    &lt;/documents&gt;</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
            <tr class="row-odd">
                <td style="white-space:pre;">&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;
    &lt;documents type=&quot;book&quot;&gt;
        &lt;document&gt;
            Someone Else's Documentation
        &lt;/document&gt;
        &lt;document&gt;
            Hoverfly Documentation
        &lt;/document&gt;
    &lt;/documents&gt;</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

JSON Matcher
------------
This request matcher will transform both matcher value and string to match as JSON objects and then match on the equality of those objects.

Example
"""""""

.. code:: json
   
   "jsonMatch": "{\"objects\": [{\"name\": \"Object 1\", \"set\": true},{\"name\": \"Object 2\", \"set\": false, \"age\": 400}]}"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td style="white-space:pre;">{
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
                <td style="white-space:pre;">{
    "objects": [
        {
            "name": "Object 1", 
            "set": true
        }]
    }</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

JSONPath Matcher
------------
This request matcher will parse the matcher value as an JSONPath expression. It will transform the string to match into an JSON object and then execute the expression against it. This will pass only if the expression successfully returns a result.

Example 1
"""""""""

.. code:: json
   
   "jsonPathMatch": "$.objects"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td style="white-space:pre;">{
    "objects": [
        {
            "name": "Object 1", 
            "set": true
        }]
    }</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
            <tr class="row-odd">
                <td style="white-space:pre;">{
    "name": "Object 1", 
    "set": true
    }</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
            <tr/>
        </tbody>
    </table>

Example 2
"""""""""

.. code:: json
   
   "jsonPathMatch": "$.objects[1].name"

.. raw:: html
    
    <table border="1" class="docutils matcher-examples">
        <thead>
            <tr class="row-odd">
                <th class="head">Example</th>
                <th class="head">Success</th>
            </tr>
        </thead>
        <tbody>
            <tr class="row-even">
                <td style="white-space:pre;">{
    "objects": [
        {
            "name": "Object 1", 
            "set": true
        }]
    }</td>
                <td class="example-icon"><span class="fa fa-times fa-failure"></span></td>
                
            <tr/>
            <tr class="row-odd">
                <td style="white-space:pre;">{
    "objects": [
        {
            "name": "Object 1", 
            "set": true
        }, {
            "name": "Object 2", 
            "set": false
        }]
    }</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>    
            <tr/>
        </tbody>
    </table>