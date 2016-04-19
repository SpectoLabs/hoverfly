# Synthesizing responses

If you can't (or don't want to) [capture traffic](#), and if you don't want to [import Hoverfly JSON](#), you can use Hoverfly [middleware](#) to "synthesize" responses to requests on the fly.
 
Since synthesize mode requires middleware, you are advised to read the [middleware section](#) if you haven't already.

# Set synthesize mode

There are three ways to put Hoverfly into synthesize mode.

You can (re)start Hoverfly in synthesize mode (while also specifying which middleware to use):

    ./hoverfly -synthesize -middleware "path/to/middleware/script"
 
You can select "Synthesize" in the [AdminUI](#) at `http://<HOVERFLY_HOST>:8888` (you will need to have started Hoverfly with middleware specified).
  
Or you can make an API call (again, you will need to have started Hoverfly with middleware specified):

    curl -H "Content-Type application/json" -X POST -d '{"mode":"synthesize"}' http://localhost:8888/api/state
  
      