package middleware

import (
	"os"
	"path"
	"strings"

	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/errors"
	"github.com/SpectoLabs/hoverfly/core/models"
)

type Middleware struct {
	Binary string
	Script *os.File
	Remote string
}

func ConvertToNewMiddleware(middleware string) (*Middleware, error) {
	newMiddleware := &Middleware{}
	if strings.HasPrefix(middleware, "http") {

		err := newMiddleware.SetRemote(middleware)
		if err != nil {
			return nil, err
		}

		return newMiddleware, nil
	} else if strings.Contains(middleware, " ") {
		splitMiddleware := strings.Split(middleware, " ")
		fileContents, _ := ioutil.ReadFile(splitMiddleware[1])

		newMiddleware.SetBinary(splitMiddleware[0])
		newMiddleware.SetScript(string(fileContents))

		return newMiddleware, nil

	} else {
		err := newMiddleware.SetBinary(middleware)
		if err != nil {
			return nil, err
		}
		return newMiddleware, nil
	}
}

func (this *Middleware) SetScript(scriptContent string) error {
	tempDir := path.Join(os.TempDir(), "hoverfly")

	//We ignore the error it outputs as this directory may already exist
	os.Mkdir(tempDir, 0777)

	newScript, err := ioutil.TempFile(tempDir, "hoverfly_")
	if err != nil {
		return err
	}

	_, err = newScript.Write([]byte(scriptContent))
	if err != nil {
		return err
	}

	// remove existing script once the new one is created
	this.DeleteScript()
	this.Script = newScript
	return nil
}

func (this Middleware) GetScript() (string, error) {
	if this.Script == nil {
		return "", nil
	}
	contents, err := ioutil.ReadFile(this.Script.Name())
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func (this *Middleware) DeleteScript() error {
	if this.Script == nil {
		return nil
	}
	err := os.Remove(this.Script.Name())
	if err != nil {
		return err
	}
	this.Script = nil

	return nil
}

func (this *Middleware) SetBinary(binary string) error {
	if binary == "" {
		this.Binary = ""
		return nil
	}
	this.Binary = binary
	return nil
}

func (this *Middleware) SetRemote(remoteUrl string) error {
	this.Remote = remoteUrl
	return nil
}

func (this *Middleware) Execute(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if !this.IsSet() {
		return pair, errors.MiddlewareNotSetError()
	}

	if this.Remote == "" {
		return this.executeMiddlewareLocally(pair)
	} else {
		return this.executeMiddlewareRemotely(pair)
	}
}

func (this Middleware) IsSet() bool {
	return this.Binary != "" || this.Remote != ""
}

func (this Middleware) toString() string {
	if this.Remote != "" {
		return this.Remote
	} else {
		if this.Script != nil {
			return this.Binary + " " + this.Script.Name()
		}
		return this.Binary
	}
}
