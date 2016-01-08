# Hoverfly admin UI

These files serve admin interface. Since Hoverfly is API first - admin interface is supposed to be minimalistic since
all of the actions such as mode change, record wiping, import/export actions should be automated and done through the API.

## Building admin interface 

Admin interface is created using ReactJS. To build it - use webpack. 

    webpack
    
During development, useful command to build upon any change:

    webpack --watch
    
To build minimized version for production, add "-p" flag:

    webpack -p
    
Hoverfly uses [statik](https://github.com/rakyll/statik) to embed whole UI into the binary:

    go get github.com/rakyll/statik
    $GOPATH/bin/statik -src=./static/dist
    
    
Then _/statik/statik.go_ should then be updated. Rebuild Hoverfly so it is then included into binary.
Commit it in order to see changes. 
    

