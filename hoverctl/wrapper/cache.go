package wrapper

import (
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func FlushCache(target configuration.Target) error {
	response, err := doRequest(target, "DELETE", v2ApiCache, "", nil)
	if err != nil {
		return err
	}

	err = handleResponseError(response, "Could not flush cache")
	if err != nil {
		return err
	}

	return nil
}
