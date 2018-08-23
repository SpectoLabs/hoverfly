package wrapper

import (
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func SetPACFile(target configuration.Target) error {
	response, err := doRequest(target, "PUT", v2ApiPac, target.PACFile, nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not set PAC file")
	if err != nil {
		return err
	}

	return nil
}
