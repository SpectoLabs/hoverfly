.. request_matchers:

Request matchers
================


A Request Matcher is used to define the desired value for a specific request field when matching against incoming requests. Given a matcher value and provided string to match when matching, each matcher will transform the values and compare differently.

|
.. todo:: Check examples and name them... maybe

Exact matcher
-------------
This request matcher will match on the equality of the matcher value and the string to match. There are no transformations. This is the default request matcher used when capturing with Hoverfly. 

Example
"""""""

.. code:: json
   
   "exactMatch": "?"


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
This request matcher will match on the equality of the matcher value and the string to match but will allow wildcard matching similar to BASH with the ``*`` character.

Example
"""""""

.. code:: json
   
   "globMatch": "?"

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

.. todo:: Buy hooverfly.com?

|
|

Regex matcher
------------
This request matcher will parse the matcher value as a regular expression. It will execute the expression against the string to match. This will pass only if the expression successfully returns a result.

Example
"""""""

.. code:: json
   
   "regexMatch": "?"

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
This request matcher will transform both matcher value and string to match as XML objects and then match on the equality of those objects.

Example
"""""""

.. code:: json
   
   "xmlMatch": "?"

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
------------
This request matcher will parse the matcher value as an XPath expression. It will transform the string to match into an XML object and then execute the expression against it. This will pass only if the expression successfully returns a result.

Example
"""""""

.. code:: json
   
   "xpathMatch": "?"

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
                <td>/documents[2]</td>
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
                <td>/documents[2]</td>
                <td class="example-icon"><span class="fa fa-check fa-success"></span></td>
            <tr/>
        </tbody>
    </table>

|
|

JSON matcher
------------
This request matcher will transform both matcher value and string to match as JSON objects and then match on the equality of those objects.

Example
"""""""

.. code:: json
   
   "jsonMatch": "?"

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

JSONPath matcher
----------------
This request matcher will parse the matcher value as an JSONPath expression. It will transform the string to match into an JSON object and then execute the expression against it. This will pass only if the expression successfully returns a result.

Example
"""""""

.. code:: json
   
   "jsonPathMatch": "?"

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