# Hoverfly admin UI

These files serve admin interface. Since Hoverfly is API first - admin interface is supposed to be minimalistic since
all of the actions such as mode change, record wiping, import/export actions should be automated and done through the API.

## Building admin interface

Admin interface is created using this excellent starter kit: https://github.com/davezuko/react-redux-starter-kit
During development, useful command to compile and save to disk (/dist):

    npm run compile

You can open redux dev tools by pressing ctrl+h (shows store state).

To build minimized version for production, use:

    npm run deploy

Hoverfly uses [statik](https://github.com/rakyll/statik) to embed whole UI into the binary:

    go get github.com/rakyll/statik
    $GOPATH/bin/statik -src=./static/admin/dist


Then _/statik/statik.go_ should then be updated. Rebuild Hoverfly so it is then included into binary.
Commit it in order to see changes.


