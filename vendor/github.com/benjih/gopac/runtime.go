// Copyright 2014 Jack Wakefield
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gopac

import "github.com/robertkrimen/otto"

type runtime struct {
	vm *otto.Otto
}

func newRuntime() (*runtime, error) {
	rt := &runtime{
		vm: otto.New(),
	}

	rt.vm.Set("isPlainHostName", rt.isPlainHostName)
	rt.vm.Set("dnsDomainIs", rt.dnsDomainIs)
	rt.vm.Set("localHostOrDomainIs", rt.localHostOrDomainIs)
	rt.vm.Set("isResolvable", rt.isResolvable)
	rt.vm.Set("isInNet", rt.isInNet)
	rt.vm.Set("dnsResolve", rt.dnsResolve)
	rt.vm.Set("myIpAddress", rt.myIpAddress)
	rt.vm.Set("dnsDomainLevels", rt.dnsDomainLevels)
	rt.vm.Set("shExpMatch", rt.shExpMatch)

	if _, err := rt.vm.Run(javascriptUtils); err != nil {
		return nil, err
	}

	return rt, nil
}

func (rt *runtime) run(content string) error {
	if _, err := rt.vm.Run(content); err != nil {
		return err
	}

	return nil
}

func (rt *runtime) findProxyForURL(url, host string) (string, error) {
	value, err := rt.vm.Call("FindProxyForURL", nil, url, host)

	if err != nil {
		return "", err
	}

	var proxy string

	if proxy, err = otto.Value.ToString(value); err != nil {
		return "", err
	}

	return proxy, nil
}

func (rt *runtime) isPlainHostName(call otto.FunctionCall) otto.Value {
	if host, err := call.Argument(0).ToString(); err == nil {
		if value, err := rt.vm.ToValue(isPlainHostName(host)); err == nil {
			return value
		}
	}

	return otto.Value{}
}

func (rt *runtime) dnsDomainIs(call otto.FunctionCall) otto.Value {
	if host, err := call.Argument(0).ToString(); err == nil {
		if domain, err := call.Argument(1).ToString(); err == nil {
			if value, err := rt.vm.ToValue(dnsDomainIs(host, domain)); err == nil {
				return value
			}
		}
	}

	return otto.Value{}
}

func (rt *runtime) localHostOrDomainIs(call otto.FunctionCall) otto.Value {
	if host, err := call.Argument(0).ToString(); err == nil {
		if hostdom, err := call.Argument(1).ToString(); err == nil {
			if value, err := rt.vm.ToValue(localHostOrDomainIs(host, hostdom)); err == nil {
				return value
			}
		}
	}

	return otto.Value{}
}

func (rt *runtime) isResolvable(call otto.FunctionCall) otto.Value {
	if host, err := call.Argument(0).ToString(); err == nil {
		if value, err := rt.vm.ToValue(isResolvable(host)); err == nil {
			return value
		}
	}

	return otto.Value{}
}

func (rt *runtime) isInNet(call otto.FunctionCall) otto.Value {
	if host, err := call.Argument(0).ToString(); err == nil {
		if pattern, err := call.Argument(1).ToString(); err == nil {
			if mask, err := call.Argument(2).ToString(); err == nil {
				if value, err := rt.vm.ToValue(isInNet(host, pattern, mask)); err == nil {
					return value
				}
			}
		}
	}

	return otto.Value{}
}

func (rt *runtime) dnsResolve(call otto.FunctionCall) otto.Value {
	if host, err := call.Argument(0).ToString(); err == nil {
		if value, err := rt.vm.ToValue(dnsResolve(host)); err == nil {
			return value
		}
	}

	return otto.Value{}
}

func (rt *runtime) myIpAddress(call otto.FunctionCall) otto.Value {
	if value, err := rt.vm.ToValue(myIpAddress()); err == nil {
		return value
	}

	return otto.Value{}
}

func (rt *runtime) dnsDomainLevels(call otto.FunctionCall) otto.Value {
	if host, err := call.Argument(0).ToString(); err == nil {
		if value, err := rt.vm.ToValue(dnsDomainLevels(host)); err == nil {
			return value
		}
	}

	return otto.Value{}
}

func (rt *runtime) shExpMatch(call otto.FunctionCall) otto.Value {
	if str, err := call.Argument(0).ToString(); err == nil {
		if shexp, err := call.Argument(1).ToString(); err == nil {
			if value, err := rt.vm.ToValue(shExpMatch(str, shexp)); err == nil {
				return value
			}
		}
	}

	return otto.Value{}
}
