# Admin UI

The Hoverfly Admin UI is available at `http://${HOVERFLY_HOST}:8888` by default. The port can be changed with flags or environment variables (see **Flags and environment variables** in the **Reference** section).

If Hoverfly authentication is disabled (see the **Authentication** section), you can log in to the Admin UI using **any username and password combination**.

The Admin UI is intended to provide basic functionality and information. Currently, it has the following features:

* Change Hoverfly mode (Simulate, Capture, Synthesize, Modify)
* Delete all captured traffic ("Wipe records")
* View the number of request/responses currently stored in Hoverfly
* View basic metrics on the number of operations executed in each mode

