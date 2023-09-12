.. _post_serve_action:

Post Serve Action
=================

Overview
========

- PostServeAction allows you to execute custom code after a response has been served in simulate or spy mode.

- It is custom script that can be written in any language. Hoverfly has the ability to invoke a script or binary file on a host operating system. Custom code is execute after a provided delay(in ms) once simulated response is served.

- We can register multiple post serve actions.

- In order to register post serve action, it takes mainly four parameters - binary to invoke script, script content/location, delay(in ms) post which it will be executed and name of that action.

Ways to register a Post Serve Action
==================================

- At time of startup by passing single/multiple -post-serve-action flag(s) as mentioned in the `hoverfly command page <https://docs.hoverfly.io/en/latest/pages/reference/hoverfly/hoverflycommands.html>`_.

- Via PUT API to register new post serve action as mentioned in the `API page <https://docs.hoverfly.io/en/latest/pages/reference/api/api.html>`_.

- Using command line hoverctl post-serve-action set command as mentioned in the `hoverctl command page <https://docs.hoverfly.io/en/latest/pages/reference/hoverctl/hoverctlcommands.html>`_.


- Once post serve action is registered, we can trigger particular post serve action by putting it in response part of request-response pair in simulation JSON.

**Example Simulation JSON**
::
    {
        ...
        "response": {
            ...
            "postServeAction": "<name of post serve action we want to invoke>"
            ...
        }
    }



