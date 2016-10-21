# API

## GET /api/v2/simulation
Gets all simulation data being used in the running instance of Hoverfly.

## PUT /api/v2/simulation
Puts the simulation into Hoverfly and replaces the previous set of data.

## GET /api/v2/hoverfly
Gets all configuration for the running instance of Hoverfly.
```
{
    destination: ".",
    middleware: "",
    mode: "simulate",
    usage: {
        counters: {
            capture: 0,
            modify: 0,
            simulate: 0,
            synthesize: 0
        }
    }
}
```

## GET /api/v2/hoverfly/destination
Gets the current destination for the running instance of Hoverfly.
```
{
    destination: "."
}
```

## PUT /api/v2/hoverfly/destination	
Puts the new destination and overwrites current destination for the running instance of Hoverfly. This requires a JSON body on the request.

```
{
    destination: "new-destination"
}
```


## GET /api/v2/hoverfly/middleware
Gets the current middleware for the running instance of Hoverfly.
```
{
	"middleware": "python ~/middleware.py"
}
```

## PUT /api/v2/hoverfly/middleware
Puts the new middleware and overwrites current middleware for the running instance of Hoverfly. This requires a JSON body on the request.
{
	"middleware": "python ~/new-middleware.py"
}
## GET /api/v2/hoverfly/mode
Gets the current mode for the running instance of Hoverfly.
```
{
    mode: "simulate"
}
```

## PUT /api/v2/hoverfly/mode
Puts the new mode and overwrites current mode of the running instace of Hoverfly. This requires a JSON body on the request.
```
{
    mode: "simulate"
}
```
                                       
## GET /api/v2/hoverfly/usage
Gets the metrics for the running instance of Hoverfly.
```
{
	"metrics": {
		"counters": {
			"capture": 0,
			"modify": 0,
			"simulate": 0,
			"synthesize": 0
		}
	}
}
```