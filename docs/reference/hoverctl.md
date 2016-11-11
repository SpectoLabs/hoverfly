# Hoverctl

Hoverctl is a command line tool bundled with Hoverfly. The purpose of hoverctl is to help in the managing of one or many instances of Hoverfly. Hoverctl does not support all the functionality of Hoverfly yet, but its feature set is growing.

##.hoverfly directory
Hoverctl stores its state in a `.hoverfly` directory. Hoverctl will create this directory in your home folder the first time it needs to save state. This directory is used for the configuration for Hoverfly, the process identifiers and the log files. Hoverctl will always check the working directory before your home directory when looking for the `.hoverfly` directory. This allows for multiple configurations on a per project basis if you require different configurations for Hoverfly.

```
.hoverfly/config.yaml
.hoverfly/hoverfly.8888.8500.pid
.hoverfly/hoverfly.8888.8500.log
```

### Configuration
Currently, there are six configuration keys needed to use hoverctl. 

```hoverfly.host``` is used to determine the host address of the Hoverfly you are trying to control.

```hoverfly.admin.port``` is used to determine the port that you would like to access the admin interface from.

```hoverfly.proxy.port``` is used to determine the port that you would the proxy to run on.

```hoverfly.username``` is the username you use to access your authenticated Hoverfly.

```hoverfly.password``` is the password you use to access your authenticated Hoverfly.

## Pid and log files
For each Hoverfly process created with hoverctl, a file is created to store the process identifier and another for the STDOUT and STDERR of Hoverfly. These files will be named after the hoverfly process with both the admin and proxy ports in the file name.

## Hoverctl commands

### start
Hoverctl will let you start a Hoverfly process. For this to work, you need to have the Hoverfly binary on your $PATH. It will start up Hoverfly on the admin and proxy ports as specified in the config.yaml. There is no limit to the number of Hoverfly processes you can start. The only requirement is that each Hoverfly process has its own unique admin and proxy ports.

    hoverctl start
    
By default, hoverctl will start Hoverfly as a proxy. If you wish to start Hoverfly as a webserver instead:

    hoverctl start webserver
    
### stop
You can also stop Hoverfly processes using Hoverctl.

    hoverctl stop
    
### mode
Using hoverctl, you can find out which mode Hoverfly is running in.
    
    hoverctl mode
    
You can also change the mode by specifying the name of mode you want Hoverfly to be in.
    
    hoverctl mode capture
    
### delete
Hoverfly stores internal state while its running. This state is used for testing your application. Using the delete command, you can specify what you want to delete from Hoverfly.

    hoverctl delete simulations
    hoverctl delete delays
    hoverctl delete all
    
### export
Instead of having to save the response from the records API endpoint, you can use the export function to save your simulation to disk. The export function will save this simulation to the disk.

    hoverctl export simulation.json
    
### import
Once you have simulations saved, you can import them into Hoverfly using the import function.

    hoverctl import simulations.json

If your simulation file is hosted over HTTP, you can use hoverctl to import it.

    hoverctl import http://example.org/simulation.json

If you have older, v1 simulations, you may still import them using the v1 flag.

    hoverctl import --v1 old-simulations.json

### delays
If you want to apply delays to individual hosts in a simulation (to simulate netwrok latency, for example), you can use the `delays` function to supply a JSON file containing the delay configuration or to view delays which have been applied (See **Simulating service latency** in the **Usage** section).

Set delays by supplying JSON file:

    hoverctl delays path/to/my_delays.json
 
Show delays that have been set:

    hoverctl delays
### templates
As well importing request/response data using import, you can also import request templates for partial matching to a response using the `templates` function. This function works with a JSON file containing the JSON schema for request templates and responses. (See **Matching requests** in the **Usage** section).

Set templates by supplying JSON file:

    hoverctl templates path/to/my_request_templates.json
 
Show templates that have been set:

    hoverctl templates
    
### middleware
This function is used for getting and setting the middleware being executed by Hoverfly. 

To get the middleware currently being used by Hoverfly

    hoverctl middleware
    
To set the middleware Hoverfly to use

    hoverctl middleware "middleware.sh"
    
The value given to the middleware function should be a string that contains either a file path, a command a file path or a URL.

### logs
Used to get the logs from the instance of Hoverfly started with the hoverctl. This command will return all the logs from when the process was started

    hoverctl logs
    
If you are trying to debug what is happening and you need to watch the Hoverfly logs, you can use the `--follow` flag to tail the logs and watch them in real time.

    hoverctl logs --follow

## Hoverctl flags

### --host
This is a global flag that can be used to override the hoverfly.host configuration value from the config.yaml file.

### --admin-port
This is a global flag that can be used to override the hoverfly.admin.port configuration value from the config.yaml file.

### --proxy-port
This is a global flag that can be used to override the hoverfly.proxy.port configuration value from the config.yaml file.

### --verbose
This is a global flag that can be used to get the verbose logs from hoverctl.

### --version
This is a global flag that can be used to get the version of hoverctl.