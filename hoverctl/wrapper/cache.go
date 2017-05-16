package wrapper

import (
	"errors"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func FlushCache(target configuration.Target) error {
	response, err := doRequest(target, "DELETE", v2ApiCache, "", nil)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("Cache was not set on Hoverfly")
	}

	return nil
}
