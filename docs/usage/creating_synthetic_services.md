# Creating synthetic services

In *synthesize mode*, Hoverfly does not make use of its cache. *Synthesize mode* is dependent on middleware. For each request that Hoverfly receives, it executes middleware, which must generate a response. Hoverfly then returns the generated response to the client.

This mode allows you use Hoverfly as a dynamic stub server. For example, you could use it with a set of template responses stored on a file system or in a database, which could be populated with data and returned based on any characteristic in the request.

TODO