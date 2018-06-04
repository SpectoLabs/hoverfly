.. _matching:


Matching strategies
-------------------

Hoverfly has two matching strategies. Each has advantages and trade-offs.

.. note::
   
   In order to fully understand Hoverfly's matching strategies, it is recommended that you read the :ref:`simulations` section first.

Strongest Match
~~~~~~~~~~~~~~~

This is the default matching strategy for Hoverfly. If Hoverfly finds multiple Request Response Pairs that match
an incoming request, it will return the Response from the pair which has the highest **matching score**.

To set "strongest" as the matching strategy, simply run:

.. code:: bash

    hoverctl mode simulate

Or to be explicit run:

.. code:: bash

    hoverctl mode simulate --matching-strategy=strongest


Matching scores
===============

This example shows how matching scores are calculated.

Let's assume Hoverfly is running in simulate mode, and the simulation data contains four :ref:`pairs`. Each 
Request Response Pair contains one or more :ref:`matchers`.

Hoverfly then receives a ``GET`` request to the destination ``www.destination.com``. The incoming request contains the following 
fields.

**Request**

+-------------+---------------------+
| Field       | Value               |
+=============+=====================+
| method      | GET                 |  
+-------------+---------------------+
| destination | www.destination.com |
+-------------+---------------------+


**Request Response Pair 1**

+-------------+--------------+-------------------------+-----------+-------------+----------+
| Field       | Matcher Type | Value                   | Score     | Total Score | Matched? |
+=============+==============+=========================+===========+=============+==========+
| method      | exact        | DELETE                  | +0        | 1           | false    |
+-------------+--------------+-------------------------+-----------+             +          |
| destination | exact        | www.destination.com     | +1        |             |          |
+-------------+--------------+-------------------------+-----------+-------------+----------+

This pair contains two Request Matchers. The **method** value in the incoming request (``GET``) does not match
the value for the **method** matcher (``DELETE``). However the **destination** value does match.

This gives the Request Response Pair a total score of 1, but since one match failed, it
is treated as unmatched (**Matched?** = ``false``).


**Request Response Pair 2**

+-------------+--------------+-------------------------+-----------+-------------+----------+
| Field       | Matcher Type | Value                   | Score     | Total Score | Matched? |
+=============+==============+=========================+===========+=============+==========+
| method      | exact        | GET                     | +1        | 1           | true     |
+-------------+--------------+-------------------------+-----------+-------------+----------+

This pair contains one Request Matcher. The **method** value in the incoming request (``GET``) matches
the value for the **method** matcher. This gives the pair a total score of 1, and since no matches
failed, it is treated as matched.


**Request Response Pair 3**

+-------------+--------------+-------------------------+-----------+-------------+----------+
| Field       | Matcher Type | Value                   | Score     | Total Score | Matched? |
+=============+==============+=========================+===========+=============+==========+
| method      | exact        | GET                     | +1        | 2           | true     |
+-------------+--------------+-------------------------+-----------+             +          |
| destination | exact        | www.destination.com     | +1        |             |          |
+-------------+--------------+-------------------------+-----------+-------------+----------+

In this pair, the **method** and **destination** values in the incoming request both match the 
corresponding Request Matcher values. This gives the pair a total score of 2, and it treated as matched.



**Request Response Pair 4**

+-------------+--------------+-------------------------+-----------+-------------+----------+
| Field       | Matcher Type | Value                   | Score     | Total Score | Matched? |
+=============+==============+=========================+===========+=============+==========+
| method      | exact        | GET                     | +1        | 1           | false    |
+-------------+--------------+-------------------------+-----------+             +          |
| destination | exact        | www.miss.com            | +0        |             |          |
+-------------+--------------+-------------------------+-----------+-------------+----------+

This pair is treated as unmatched because the **destination** matcher failed. 


Request Response Pair 3 has the highest score, and is therefore the **strongest match**. 

This means that Hoverfly will return the Response contained within Request Response Pair 3.

.. note::
   
   When there are multiple matches all with the same score, Hoverfly will pick the last one in the simulation.

    
The strongest match strategy makes it much easier to identify why Hoverfly has not returned a Response to an incoming Request. 
If Hoverfly is not able to match an incoming Request to a Request Response Pair, it will return the closest match. For more 
information see :ref:`troubleshooting`.

However, the additional logic required to calculate matching scores does affect Hoverfly's performance. 


First Match
~~~~~~~~~~~

**First match** is the alternative (legacy) mechanism of matching. There is no scoring, and Hoverfly simply returns the first
match it finds in the simulation data.

To set first match as the matching strategy, run:

.. code:: bash

    hoverctl mode simulate --matching-strategy=first

The main advantage of this strategy is performance - although it makes debugging matching errors harder.