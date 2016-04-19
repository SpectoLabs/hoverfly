# "Virtualizing" services

Once you have [captured some traffic](#), you can use swap the real service for a virtualized Hoverfly version.

## Set virtualize mode

There are three ways to put Hoverfly into virtualize mode.

You can (re)start Hoverfly without specifying a mode (as virtualize is the default mode):

    ./hoverfly
    
Or you can select "Virtualize" in the [AdminUI] at `http://<HOVERFLY_HOST>:8888`
    
Or you can make an [API](#) call:
    
    curl -H "Content-Type application/json" -X POST -d '{"mode":"virtualize"}' http://localhost:8888/api/state
    
Now that Hoverfly is in "virtualize" mode - provided you have either [captured traffic](#) or [imported Hoverfly JSON](#) - you can use Hoverfly as an alternative to a real external service.
    
A [more detailed step-by-step guide to capturing and virtualizing traffic is available here](https://specto.io/blog/speeding-up-your-slow-dependencies.html).       