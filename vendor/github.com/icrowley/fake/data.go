package fake

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/data/en/adjectives": {
		local:   "data/en/adjectives",
		size:    119,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/xSLwQrCMBBE7/NXhUj1ZLDiPaZjWUi7sm7w97u5zPCYeZMYZjZM3RVZ/zQ8yxa4yrsR
14Q5LzGYl8FJfhWXxuomFZn2UdvLUYlFVuIVobgd3+64dx+VxOKNB3f14W/ipeEMAAD//9gYba13AAAA
`,
	},

	"/data/en/characters": {
		local:   "data/en/characters",
		size:    72,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/wTAwwHEAAAAsH+mOWOc2nanbwKhSCyRyuQKpUqt0er0BqPJbLHa7A6nm7uHp5e3j6+f
vysAAP//cadMlkgAAAA=
`,
	},

	"/data/en/cities": {
		local:   "data/en/cities",
		size:    4837,
		modtime: 1464022103,
		compressed: `
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
`,
	},

	"/data/en/colors": {
		local:   "data/en/colors",
		size:    128,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/wzKQarCUAwF0Pnd1f+CRUQsUgSHoe9SQ9PEpka3b0dncm5s+LciuiQdD5rFF9cUn4i+
8mXEXcP4xsmbToGBYujVZxxrG1XQhTV6RsNF6oO/tWSRVCeGyrVCN+KQumzh+8jY6Wskzk+Z9RcAAP//
/hYIpoAAAAA=
`,
	},

	"/data/en/companies": {
		local:   "data/en/companies",
		size:    3157,
		modtime: 1464022103,
		compressed: `
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
2PwEcuakK25E/w8AAP//Fi4aPFUMAAA=
`,
	},

	"/data/en/continents": {
		local:   "data/en/continents",
		size:    68,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/3Iszkzk8ssvKslQcMxNLcpMTuQKzi9F4jmWFpcUJeYAVTmmgQVcS4vyC1K5HPNKEouS
S0AigAAAAP///00YyUQAAAA=
`,
	},

	"/data/en/countries": {
		local:   "data/en/countries",
		size:    2774,
		modtime: 1464022103,
		compressed: `
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
`,
	},

	"/data/en/currencies": {
		local:   "data/en/currencies",
		size:    1800,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/3RVS3brNgydcxUcxsvw9yWxnbiW256XGSyhEmuKdCEyfsrqeymLSnpOOrEJEAQuLj6a
/1U35EwXyOl8VnN7xj/pHV8YQs0CYWUcCe5aSCXph9+dCVzpIlDgbqZX3trhXmp2Abb6wJ2HHLsgZJOD
bJE0kIuyMdYaV3dqQQ211E0mkIWMyzEX5GpLFXeNPtGFIMuZKv/Fnm1tYqs3Qq5MorSx+oy4EPowVh+Z
kNwi2ppkyO6d1JIcfbFcbuZ6sVzPX7OrQbGeLycZmHlMDWek+TOCuiO71rizUcf9Qj98czFTS299ezaZ
l+XmMPn0oF8fE6vJyDFU4inAdhsdMPZXiZ0++OgqXH1w2cDvNZ6tKfXWS3S9WiEOyUVvBe8FYgxd2bB+
+AEuyPUzvcd1p1YecBDJfXq441nFALcPLxwaFkuINNM/orEVg5Y1dYHF6SWIO58ZrzNh67q/hgxt3QWf
+gYg8K/WUeB4Y/42k/nGuOR7AHMhUgMDnIkYseoJ/R3zD28r/RphiLMw48FKqGxaYvWYHMPh/wB/9K4G
RfjJCB6jQ/l7vfFiXFBPJQ/vE29R1NN+o4srl4ZsinFDd6IwdQNDVyGzY7wyUEBIZborDDXQSKIUz4bj
P7lzn+TufiToqRNiq1/4posGw5WsA9le74ywesYMpCbIUJ/pCqc/2alnLxVN47Bl12N8pulBCzAGsvAx
NDP9p3dqG29kQrbfMcbZuwxiF39xe/ZR6sz7niz1QzbwV5swaNCSQGVwjGKCQQOOye/5lyn92DZ7L76E
tDKCEe7UlzJ8ViHl+8Y0MJGTe/Fyoz7362s70Xegy30djeEOhAdVmu6xBQ4sUb9EfscCKLxNGlBhrleD
koywDhbj47688dZge7xZH3oIEmKNCq+7MmKNqN8okCB8n+IffTssvwR6x2aSh3Ps7kU/p7AFxcroudB5
oG54XRj7zpLDFuCTrijPlHZh/TsSpDy3Sebk/uTv96mI2MTDjj2CsFFzr3GqbXFNmxFpciDYf9esnX5A
I89UIUbvyF2mvi1iaqOxD4obV+yGCniBZMLHvXK5LU5kbjBPRGT4p4bMYIINHdQJI2QqqnTSnPyZav9p
GfEl+fxsnKJceOhzUuOXIxGn1y00+IBMDTRebpFM5duM9T8fmynGH+z4I2LC9AL1fSfhb3V6E1lCujMc
HLV472r1RsM23t6wSUj9GwAA//+8vc6ZCAcAAA==
`,
	},

	"/data/en/currency_codes": {
		local:   "data/en/currency_codes",
		size:    519,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/ySRS3qAIAyE93MrMKDIQxqwRe9/kA64mC9CYpL5cbfiboLdVmxGYG7BWR+EoiivYDs8
XqOYdcYXmJQgvJ//GG2r3vQGy7M9KOmwltF52MyoCXYvq/e4PIbx2BJnlefTVTEq73TDoRHbw9wbIZFy
GcJ8STvGJnB7hXNx7eJPgQ8ZXv2qG+bGrrJqjyg4bo/QIobo8hJmVOpHEFxFSA2hJ5z583tegugaov4h
/gkSeST2yI8isy5zZh4F2XwzJptyRVxZUaPSg6C6gnpMPx018bs7/JCdXoVK0Nui8TzMjrYLGj220Bff
Ode1uvZN7NeENfQ6+Xfu0w+L3hkLpQ+M+95svsOvs5THL3Nv/vhMzos32S528050cZrzV+Sec2/m/gMA
AP//90+SpwcCAAA=
`,
	},

	"/data/en/domain_zones": {
		local:   "data/en/domain_zones",
		size:    753,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/wzLQbalIAwA0XntMiBGhCCPgLauvv/0niqJyIYkZEcUyUhFDGnIhfyQgTgykYU8yD/k
IwghEDZCIuwEJRyETDgJlWCERrgIg+CESbgJD+ElfEQhRuJG3IlKPIiZWIiVaMRGvIiDuIg38R/xJX5s
ie1kK2zGdrF9pEhKJCUdpEFy0iQt9sx+shd2Y7/YBypoQDc0oTuq6IFmtKKGNrSjP/SvdHSiC33Ql6Nw
GEfjGByTY5E3ciJXspEb+SL/yIPs5MmZOI3z4uyURFHKQckUozTKHw7KQ3kpH1WogRqpmVqog+rUSV3U
m/pigkVswxKm2IEVrGKGNezCOvbDBubYxBZ2Yw/2D/vbP5rQIi3RdprSMq3SLlqnDdqifVxGF3qi73Sl
H/RCr/Q/b/RBd/qkP/SXnzAS42I4YzEeXPCAR3zDE674gWf8xAteccMbfuEDn/jCb/zFP2ZkbsydqcyD
eTILszKN2ZgXszMHczJv5sP8WMJSVmEZy1kv6+MW7siduJU7czfuxbPzOG/inbyLT/iM7+F/AAAA//+I
zWV08QIAAA==
`,
	},

	"/data/en/female_first_names": {
		local:   "data/en/female_first_names",
		size:    686,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/zxSXa+bMAx9978q7d3uKLeqSrU9G/AaSyHZnHCl7NfvGKa9mdjH58N8sDW6czWdlWnQ
tDB1bBMb01vUPzxJDdRLSvpTjDCOsXErnLx+sUmlS7ZcQwO6MN04zY2uaCTqpNZG7xJRj5wW7LzklJjO
bDnSY8PqMaDGMp2DxCg08IaxEfyBrrpOYrHRRabsD72UojPaQS1Ko3NLNUDQKb0kMn1I1AINHcjh47Q2
dED3kElmwL6rvTRh/so1RIGqO687kK0GdhpzGDt6rPIrcFI51LZE52BaqibZYxDqOcH9GbvE/PWLwbsU
56Q+t1noohihU1TU/RaBeZd9nJ5iUjwOrKSvMXusb5/iNL2wc4m1iFjiYrLsgg+SPqN7KsHt99uiiPCR
yy7GWa4IsdFN5xydcmn/VR+uGz3DQd3J579skwL/5BVxfUNyhzEaoIke7GfZ/bq0E24oxyLD55A3h/qx
3LQD598bjoDyxx5ih3Pr4R1n2Cb8JBmGny7nHlrEveiWbcUXbziEB7bvAmhQtGH2bQU3TE6a/gYAAP//
/8KRCK4CAAA=
`,
	},

	"/data/en/female_last_names": {
		local:   "data/en/female_last_names",
		size:    1764,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0RV3XLiOhO811sFskAlkFCG+qjvcrAnWIUtcUYSjnn60y2yOTebtf6mp7unOYw+9+4t
9iHF4E5+GLyMCQtBk1tYnIJ7lbtPboctNZ7gwV2Mpu4o8xDNvYROjavHPo6Cy9Je62u9z+o2Ysb7Ytk/
j9y4uRZrvfws68M18ewrhuUgdsVnZ/5SsLHVCde3qu4kwxUQNjIM7gVogvt/LOHiNmpBgOHh3j0+T7jY
Z7eNN6xsANsd2pizW5viykvH/hbCl9YxPGTAqQ+tXS0BBss7n9tecW+v9kSmlpM7FgvY3fck6ZZweryd
6zEAxsafu4Tk/nSTWIfdiFP4PmTFQnYHCXj0AeYqHU28gDPXqHY4Gq9cv0hwCz64K3brZ4D0g86u8Xc1
4akbqjS+7VmAcOO320SWAzX454inodpe81OOtQluy+jZxZuM2DtJ5g6Ejdfk3lFsJjLq5/bmWwWAEBRs
nWLsgMDogyamBJb/yozWdATWNw1XtgiabHb7OBH7NkKCveQfDJty6fHCCj6p5VMPiTIhlEw/Hfw4Rryx
iin/pwitN0vIkFm/KzzXlJT4/tr81xeM9OoF6sqMs7uZ8FcRFGzQ7ZCfvfcyukOBCtCF3hkE7bGtKg7+
qynD6NZh+3PSQDXmEAfsr3zqqSgupx8D801/rt5vuxhkgG5WHjRw6unITxj5gQZGilyM3K/gOBJ1AnA0
r+cz+632hy3uLHksLb2zj9V5mxLqH48BwvMyfdWmNIDfRZw7VHtOn1WS3qlVN1N9U1LyjU0IHmsrZBGW
TcBkHUv2MsEvA33QUOofY1eVwPUCBF0x7cFjGiAhDlod++Q+YLqIv2AVB99DHbCV2qXwNuyhaAiRAa6m
aonXEgJ9UT82pdr1cNPQUmN4lYMECm49OdjLjKt7r1b9RzPtYKBep+RejIqgw0t4pg8VaYQ95/r459/J
NQwc67Yg/BDmriYFBm8JKPAc3bAw6ThRW0E9pJaxQlM8jWQcrhUmqqlD92JjykYvL7lVZTmp3ImfOaJ4
ZDI2VF3CfFn22H5QYcAuiWn3HMU6Z0o7SLgiFHiVoFcgGmZbl4zWa9LMxN72o+9ybcn9T9I/DMGl4L1h
iMhUVc4Nit1orc+B6QBHhHyB9Wz+DYjfMEelmu916oAAAwNzmFvA5PEGO7dxZo5ojdW72G9K0RnINFhr
9RuxUO8h5OtOIn3GFNNmuXpaAvv8uJQZd9aKSFP+GsQzDek75PPoVqX+kGxn5KF7VeFcDTQinzZGT8NG
IqcGJ7ZScwTcPTxBF2RtRRX5ywSnIJYiEAGxhogpmUjLTjsfBHpOqPVvAAAA//9sQKZ55AYAAA==
`,
	},

	"/data/en/female_name_prefixes": {
		local:   "data/en/female_name_prefixes",
		size:    14,
		modtime: 1464388191,
		compressed: `
H4sIAAAJbogA//ItKtbj8gXhzOJiLkAAAAD//2lyagoOAAAA
`,
	},

	"/data/en/female_name_suffixes": {
		local:   "data/en/female_name_suffixes",
		size:    29,
		modtime: 1464388191,
		compressed: `
H4sIAAAJbogA//Lk8gQhIA7jCuPydeFycQnmCsgA0mG+XIAAAAD//6G38R0dAAAA
`,
	},

	"/data/en/female_patronymics": {
		local:   "data/en/female_patronymics",
		size:    666,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/zSSzY4qOwyE934r/gbUw1whQHfWhvZ0ckjbyJ2A8vanktHZIJF2VeorZ+BZFhosKJ3t
Jp7pK94DS6LvmFLkmbb8iiOd26mPtMFv6opFnoGuwWZecOpxyfYM4pjXCPmJS6Iv9gdtTTmNtBfzSehT
VCUHumR5idJufDfbtUduCfrkSnMwrRh9RaWBF1M45RzkTXv2Stc4Ww61h6AjO44G+flxqfThrA+63C1n
2nm8t3sQS2E6OvRnrrMp0rhM5t0iFIa8mWyRLS70zSkD5MQZBg86Sft3YDdk21qZEogPohBsUAatPIfi
dK6dYMLsYEJD4ZYd+lXqvQ6Ay6EdliWD69qv3Is35E+JrROey7/mBVETGj7y20XvQv9hA9auPlultegf
nuGy9oJva1CPaGk1YmGHXsiHywiSqvgKv/pbOB2tgBC4MldagQmZoa10sL6IXZkEkgZmuKosi6REa7vd
Kv0f79m8LbXl37nKkukUWtonXW0c4Yt52jjHCdhgvQR+K20S/zJcBGdNAUF/M/3lYdW7VuQQ5xZKs2k0
+hsAAP//YYz4ZpoCAAA=
`,
	},

	"/data/en/genders": {
		local:   "data/en/genders",
		size:    11,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA//JNzEnlckvNBVKAAAAA//+SJ0gyCwAAAA==
`,
	},

	"/data/en/industries": {
		local:   "data/en/industries",
		size:    4922,
		modtime: 1464022103,
		compressed: `
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
AA==
`,
	},

	"/data/en/jobs": {
		local:   "data/en/jobs",
		size:    2246,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/2RVXXLjNgx+1yk40wP0DKnjtbezidXI0zwzFCRxI5Ea/njj6fTuC4ikSMkvHvMDIAIf
PoBPQmivHDtobVqpuNOmStjxC4R38gYr8gazAQvKcYLZH/+9/p9sHM3lWaqePVkrbWlxnpt79dROUqHB
hK/svbbWS9dJAZiU4qPu2UEa4aVjz2Blr9jEFe+j9W4dq43uDZ8mgtbPviSnjEArOatHrtQGx/AOrNUL
poXkDkrMOz1hVlqxZga0jhgWkv5LakukWCcRVxH0bQ/uz5KQmGY0W6nwy6xBCCaLNd1g1POE/CbP6jDA
JAUf2VH16I3ZHgYJXao/o/ImS6cROaSwnGh10NPklXR3dvHOABfDzjqDsqG89XYEvQOzprgp4IA1s9z+
6pk7vhESAR/cYotzT7UJwc+kopGd770EMlX7ivD86fTMGj/P2jh2BTGohd0qEgXpU9KAoA/rjjV8BFsd
W0n3H0fEzZa+9CeEHtVNGq2mkExBx9ZAd1frOLAGhAFHUv6GlSqKwRJvkkRSIJHEb1q3bOkjHk7fm1QJ
uRgxSIdJVidAJSKStHoCjXLfph6xRP/J8HmQIkoB7WfgoxuwA9TaxeUM40z2T3ZBtpaeZKgg9OxxkHC6
rfZGgN2P5N6ckvyuOm3SRESFZBvqRlGNPjTjb9Qe9igvlR/QkzldVf2QH4YbSueFm0/YbpACSze8YP48
qjxevlIVBjx6hiJevRiBG1brXwitnq/eWAi/OOpc0ABrYvMihJ+X2kgAA9I3UwPDPtoztNJbc+wiFYb/
7mGS25/euuVvORv1wJE6QZ+sh7tdygi33Iuqcfm0Xricbt5GedWFFIrzZkoR/4kKW2n7x6PCcQkctHJG
b0SfTOttbyh03Lb5X8n/G5AWwUDLAokoECQYtbcjp8TDUs1QiGx4B3hvVmQIXIa5iArn7SNUNbDoar+M
IryuigTE1RBPj9MaDQ9URPwhJb1Ev2tUJzrpzv3iBohc68eQSMK2m2eFr2Dd3jYDMsFq7oY48VXjeNex
/Vsb0EbgCnXRa/cGNc6gfrxZK5TFlBTGjKVN/7iyr/yrZPhKTwhG5HX2bhat/Fuz/ObRaZ1dOiRJBxnR
XBC6UyVBYZW/w8e64UIaAdm8AL8DAAD//7ExBFjGCAAA
`,
	},

	"/data/en/jobs_suffixes": {
		local:      "data/en/jobs_suffixes",
		size:       12,
		compressed: "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\xe4\xf2\x04\x21\x4f\x2e\xcf\x30\x2e\x40\x00\x00\x00\xff\xff\xa5\x82\x0b\x05\x0c\x00\x00\x00",
	},

	"/data/en/languages": {
		local:   "data/en/languages",
		size:    821,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/yxSUXYaOwz91ypmK0ASyAPyaOGcnuZPMKotsGWObIfCXxfUVbQLq+zJx4yuZUnWvdLs
h/IFUTLMwhGFUWAWPSqfYKZ47CbS5M8ZI2WC2T2iIswepAxzCqg1t4A5icNgLs4BI8I85Z44r8FZxY60
V1hgwWDnhWfpZ01YWsDiQScPT9ZIbsbqP3n+IG+2lnbzSOIuHuFZXGgxz7mk/sgLnycT+MqSDEgv8qIk
lviNchlelHujS3RYH7CkpK6fSWMzSnSBZTVy8ue3gbOhwrBCbs0NC6UUCFZ0VLrBimW0O1YetqlUWFX5
ZPl6Mk3s8gSvMiYj2J32th+WSKH5jX5z/odX7AKsUQRHNJt9tJIGHnjxsPaRFNZJyaLXVcdGaX1Xd3/A
BpN95aMV2nDxdRrfpv6keEzVqFnsFk80TgptTXKH+d7B59++2FBpTWz//kr29LaxNsW3KYzpoyWa5qn3
a6iQkFMWeBvpSCbHG13b0N/oNrxTm+o47NnJsEFxFZ0FJC2eVIZ9Kj6144267P8r3xF2eGXbKynJYPbN
kHbJdqmPeGf51dXW4a7K2ZYSvlRbk4rwNcWJ9N6ANTG9sL9O+7O/oefmvbUp7m/U1TvgmS/2jxzgQMEq
myk1wsEjwyFdhp1ticAhG200c0NB+M5jz36vocK/AAAA//9x6eo3NQMAAA==
`,
	},

	"/data/en/male_first_names": {
		local:   "data/en/male_first_names",
		size:    665,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/zSSzY4qOwyE934r/gbUw1whQHfWhvZ0ckjbyJ2A8vanktHZIJF2VeorZ+BZFhosKJ3t
Jp7pK94DS6LvmFLkmbb8iiOd26mPtMFv6opFnoGuwWZecOpxyfYM4pjXCPmJS6Iv9gdtTTmNtBfzSehT
VCUHumR5idJufDfbtUduCfrkSnMwrRh9RaWBF1M45RzkTXv2Stc4Ww61h6AjO44G+flxqfThrA+63C1n
2nm8t3sQS2E6OvRnrrMp0rhM5t0iFIa8mWyRLS70zSkD5MQZBg86Sft3YDdk21qZEogPohBsUAatPIfi
dK6dYMLsYEJD4ZYd+lXqvQ6Ay6EdliWD69qv3Is35E+JrROey7/mBVETGj7y20XvQv9hA9auPlultegf
nuGy9oJva1CPaGk1YmGHXsiHywiSqvgKv/pbOB2tgBC4MldagQmZoa10sL6IXZkEkgZmuKosi6REa7vd
Kv0f79m8LbXl37nKkukUWtonXW0c4Yt52jjHCdhgvQR+K20S/zJcBGdNAUF/M/3lYdW7VuQQ5xZKs2m0
vwEAAP//9tvd8ZkCAAA=
`,
	},

	"/data/en/male_last_names": {
		local:   "data/en/male_last_names",
		size:    1764,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0RV3XLiOhO811sFskAlkFCG+qjvcrAnWIUtcUYSjnn60y2yOTebtf6mp7unOYw+9+4t
9iHF4E5+GLyMCQtBk1tYnIJ7lbtPboctNZ7gwV2Mpu4o8xDNvYROjavHPo6Cy9Je62u9z+o2Ysb7Ytk/
j9y4uRZrvfws68M18ewrhuUgdsVnZ/5SsLHVCde3qu4kwxUQNjIM7gVogvt/LOHiNmpBgOHh3j0+T7jY
Z7eNN6xsANsd2pizW5viykvH/hbCl9YxPGTAqQ+tXS0BBss7n9tecW+v9kSmlpM7FgvY3fck6ZZweryd
6zEAxsafu4Tk/nSTWIfdiFP4PmTFQnYHCXj0AeYqHU28gDPXqHY4Gq9cv0hwCz64K3brZ4D0g86u8Xc1
4akbqjS+7VmAcOO320SWAzX454inodpe81OOtQluy+jZxZuM2DtJ5g6Ejdfk3lFsJjLq5/bmWwWAEBRs
nWLsgMDogyamBJb/yozWdATWNw1XtgiabHb7OBH7NkKCveQfDJty6fHCCj6p5VMPiTIhlEw/Hfw4Rryx
iin/pwitN0vIkFm/KzzXlJT4/tr81xeM9OoF6sqMs7uZ8FcRFGzQ7ZCfvfcyukOBCtCF3hkE7bGtKg7+
qynD6NZh+3PSQDXmEAfsr3zqqSgupx8D801/rt5vuxhkgG5WHjRw6unITxj5gQZGilyM3K/gOBJ1AnA0
r+cz+632hy3uLHksLb2zj9V5mxLqH48BwvMyfdWmNIDfRZw7VHtOn1WS3qlVN1N9U1LyjU0IHmsrZBGW
TcBkHUv2MsEvA33QUOofY1eVwPUCBF0x7cFjGiAhDlod++Q+YLqIv2AVB99DHbCV2qXwNuyhaAiRAa6m
aonXEgJ9UT82pdr1cNPQUmN4lYMECm49OdjLjKt7r1b9RzPtYKBep+RejIqgw0t4pg8VaYQ95/r459/J
NQwc67Yg/BDmriYFBm8JKPAc3bAw6ThRW0E9pJaxQlM8jWQcrhUmqqlD92JjykYvL7lVZTmp3ImfOaJ4
ZDI2VF3CfFn22H5QYcAuiWn3HMU6Z0o7SLgiFHiVoFcgGmZbl4zWa9LMxN72o+9ybcn9T9I/DMGl4L1h
iMhUVc4Nit1orc+B6QBHhHyB9Wz+DYjfMEelmu916oAAAwNzmFvA5PEGO7dxZo5ojdW72G9K0RnINFhr
9RuxUO8h5OtOIn3GFNNmuXpaAvv8uJQZd9aKSFP+GsQzDek75PPoVqX+kGxn5KF7VeFcDTQinzZGT8NG
IqcGJ7ZScwTcPTxBF2RtRRX5ywSnIJYiEAGxhogpmUjLTjsfBHpOqPVvAAAA//9sQKZ55AYAAA==
`,
	},

	"/data/en/male_name_prefixes": {
		local:   "data/en/male_name_prefixes",
		size:    8,
		modtime: 1464388191,
		compressed: `
H4sIAAAJbogA//It0uNyAWJAAAAA///WJopzCAAAAA==
`,
	},

	"/data/en/male_name_suffixes": {
		local:   "data/en/male_name_suffixes",
		size:    37,
		modtime: 1464388191,
		compressed: `
H4sIAAAJbogA//Iq0uMKBmJPLk8QAuIwrjAuXxcuF5dgroAMIB3mywUIAAD//29NCl0lAAAA
`,
	},

	"/data/en/male_patronymics": {
		local:   "data/en/male_patronymics",
		size:    666,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/zSSzY4qOwyE934r/gbUw1whQHfWhvZ0ckjbyJ2A8vanktHZIJF2VeorZ+BZFhosKJ3t
Jp7pK94DS6LvmFLkmbb8iiOd26mPtMFv6opFnoGuwWZecOpxyfYM4pjXCPmJS6Iv9gdtTTmNtBfzSehT
VCUHumR5idJufDfbtUduCfrkSnMwrRh9RaWBF1M45RzkTXv2Stc4Ww61h6AjO44G+flxqfThrA+63C1n
2nm8t3sQS2E6OvRnrrMp0rhM5t0iFIa8mWyRLS70zSkD5MQZBg86Sft3YDdk21qZEogPohBsUAatPIfi
dK6dYMLsYEJD4ZYd+lXqvQ6Ay6EdliWD69qv3Is35E+JrROey7/mBVETGj7y20XvQv9hA9auPlultegf
nuGy9oJva1CPaGk1YmGHXsiHywiSqvgKv/pbOB2tgBC4MldagQmZoa10sL6IXZkEkgZmuKosi6REa7vd
Kv0f79m8LbXl37nKkukUWtonXW0c4Yt52jjHCdhgvQR+K20S/zJcBGdNAUF/M/3lYdW7VuQQ5xZKs2k0
+hsAAP//YYz4ZpoCAAA=
`,
	},

	"/data/en/months": {
		local:   "data/en/months",
		size:    86,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA//JKzCtNLKrkcktNKgIzfBOLkjO4HAuKMnOA7Eour9K8VCCRU8nlWJpeWlzCFZxaUJKa
m5RaxOWfXJIPov3yyyACLqnJEAYgAAD//6SVRvJWAAAA
`,
	},

	"/data/en/months_short": {
		local:   "data/en/months_short",
		size:    48,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA//JKzONyS03i8k0s4nIsKALSlVxepXlAnMPlWJrOFZxawOWfXMLll1/G5ZKazAUIAAD/
/8Xt5DkwAAAA
`,
	},

	"/data/en/nouns": {
		local:   "data/en/nouns",
		size:    128,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/yyNwQrCQAwF7/mrWvAiXrR4X9qnBHeTJUkV/96seMpAhjdHrgGjWVs3uKvR9eOBRjfG
Ox9n3SVoLg6attKHO7Ve+c5JB+PtgTxlfSLoxEGLFfHGMcQLVvArYdnl15AwrfWfU0EufwMAAP//aKza
G4AAAAA=
`,
	},

	"/data/en/phones_format": {
		local:   "data/en/phones_format",
		size:    26,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/1LWVVaGYSDiQrAAAQAA//+15odeGgAAAA==
`,
	},

	"/data/en/state_abbrevs": {
		local:   "data/en/state_abbrevs",
		size:    149,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/wzGwQ0CIAxA0fvfqgGUCi0GEMT9B5HkHZ5UpCA/pBOE0AiTmHhUnkJWNKIVdVQog3Ko
giUsYjeKOTawhk084QvP+As3/OABj7RMK7TOW+jKCIzIdOaXz2Rdwr4WW9nnHwAA///Jy0DGlQAAAA==
`,
	},

	"/data/en/states": {
		local:   "data/en/states",
		size:    471,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0yQb2/CQAjG3/Mp+lUWzWbn1GV/NHuJLeuRXsEcV7X79OOqyZY05cdBHnh4iHjEAcGj
9R4S/6iU2KMYGiww8rcmYYSFRk3YqoMINZmbMcOSIl4wETx6kVuEJ9LUeffKn5mhbjEo1DGyKBvU0jK6
fq0XhPVtxJokj00/wYuObHN5gyzk/zRFlNbBDJswGuVssOEmcIfi4HuYZu9ns/KdTjyzjslBJRexLR3T
bG5LZ2xLuFQrHE4W2Pcu2TMlo2nGDV250Rm/NPWw1ZRDtcCk7gDv6RL7MnUXWGHXR3foF9wl6lTglURs
imcsJ3sL2lJV2+ziXcf/Urf0LvVBxYoROV39Jp8ZA+wpDe4B9uwXLXIH9J2lyz7mQJarvwJbo2Ls75MO
3vIbAAD//57fcibXAQAA
`,
	},

	"/data/en/street_suffixes": {
		local:   "data/en/street_suffixes",
		size:    132,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/xzJwQrCQAyE4fu8VakHEQ9FC57DOkgwZCGbraxPb+vlG4Z/MuPAtNE7MdOTgVmj2P5q
j8QctTX1F06hG3FWM1y6l9TquIoTi8T7z0fGvq1hMSk8/AqWqp64VXninkEmVkYcfQ1Rw0PGLwAA//+w
pDANhAAAAA==
`,
	},

	"/data/en/streets": {
		local:   "data/en/streets",
		size:    4162,
		modtime: 1464022103,
		compressed: `
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
`,
	},

	"/data/en/top_level_domains": {
		local:   "data/en/top_level_domains",
		size:    37,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0rKrOJKzs/lysxLy+fKS8xN5cpLLeHKL0rnSs8v40pNKeXKzcwBBAAA//8EGszrJQAA
AA==
`,
	},

	"/data/en/weekdays": {
		local:   "data/en/weekdays",
		size:    57,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/wouzUtJrOTyzQdTIaWpxSA6PDUlD8IKySgtAjPcijJBVHBiSWkRiAEIAAD//38mK3c5
AAAA
`,
	},

	"/data/en/weekdays_short": {
		local:   "data/en/weekdays_short",
		size:    28,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/wouzePyzc/jCilN5QpPTeEKySjlcivK5ApOLOECBAAA///cuZyeHAAAAA==
`,
	},

	"/data/en/words": {
		local:   "data/en/words",
		size:    1685,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0xV23LrMAh85y+xRBtmdDtCyvTzz4KcjF86sc1ld1koF2Wj1JvJv81rT+K9aMj8kSkt
q5HponcveyxeUolT2sZt6a6Ue+lT6r8txEhRriTsTzqMCfXw/iItpZO2t7SFaHojEKVQWZbHmBLP9NIl
aXW6hBfS3up/s6bFZLstkr9RNPHViQ2t2hdqz8DnP372L36ZZBRV/lBCLgIr/za94aJvN4+hCRS9yZOd
pyg1saTetYmzuWk6K5CO6vEqpOEKGtEsLQlQWYciH0TLE1DrjdqGWK6SbqMK6ISmo0+GPkmztwSTApoQ
CnVP54CPLMwqsl1YYH/iRpY0RUimqvjB6N40YiF067YmoMufzBTigzaSSkFE6hMIPEpqP1V8fPVZPgjc
z8cRtkFx6A1Wuyc0xSyBJHBqRjvYAfVrMH16LFBhaCgNlUiATTeIThlTXvCdeyRmpM1rfKHcooqZUGjR
9KUFUhaxpRyiQVsYCS07DIS5m2nFxPNOWveZ/FUY3l4KFGOymBw/F975tvbqC9V96p/e7ncp0GfBsMt9
AeUmPilKev2Pu46rv5jw5i/JAG2FcphwT4klYQiU9tCswcvdMWZ/awaWMM1ZQxuaNLZlQxXqFXxIzTMA
CbXnDCNWddHlLAtES7sMDq7958dLALzJ9K8ABg3xhjFrKIxB2T1H98jtOvzERrE/v3ge32dxSeITJzS0
yJTYTnEiEByUW4LByB1b9JLZb5NDMGQa5Fy+DdcnveivXxrqw9PSjnty5qo16rrAhAoo5c45K2THILA7
ZFBfsg7z859WoVE44ZA4TJwxpLoHEHbk+xqZDQuN3hzVzsez2TCiFAw7I+mg12t/LIu2eMgcYvnpeVrb
/YyLdwT7ng45cmKSweWLoPkO0mDczrgdy89jjwvltePWBvAYosOTS9dZnnsUTRLWIbbaIRoDOomvv8Q1
kcehMCe2s8L92JVwo6zH9rgLp8SB9wAU9VlIzP90e2mCIO1cOjQbCr9KbEZyP+51QwPc8//j0zzQVdZY
ktDlSOrn/Lw8okOD/wEAAP//Xc5jD5UGAAA=
`,
	},

	"/data/en/zips_format": {
		local:   "data/en/zips_format",
		size:    8,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/1KGAC5AAAAA//+H3Sc9CAAAAA==
`,
	},

	"/data/ru/characters": {
		local:   "data/ru/characters",
		size:    119,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/wTAVRXDAADF0P+oGYO7McOZkzG3tfDiqDdnciFXciN38sATeZIXeZMP+ZIf+ZOClKTC
EY5xglOc4RwXuMQVrnGDW9zhHg94pEGTFm06dOnRZ8CQOgAA//+3jCKOdwAAAA==
`,
	},

	"/data/ru/cities": {
		local:   "data/ru/cities",
		size:    2271,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/2xWaXLyRhD9zym4AHdE4NhO4QrxVnbkBQOV8q/EYhGMhZCvMHOjvPd6RojU94veppfX
i/B/+q0/+DIMeyBB+CpkvvBHvw1D3/g1xUcIfnwBErI9mNoXIRMbfpPxBg8P4QZPqyR2eNIM/DtkIv3a
V+YjCyP4KexpuOn5W7zeIYdLviZzgMNG8bcUGLk277fwPsU7MXdwcUAFhf8mw5cOYZz/7vs3xew4+p9+
qUT+8BtpYFSGKckDvG1hssZLWOBt1RXX3hGsO8iVpt+wGPNPEQRMSKwiH1kd2DDxK4kYECXsww2B7Pl7
FeHClUjkdIUwR9OU6gm6AYAcE7lP9RCvlgWWLFaiB6U4Su17CGPElbtHQwNxy4Tmo/9ORMXYAoLRV2Fs
qT6qcrB4XBirME8wZ+Hx/RPD453rIv4MIceDuTVkd6fIz5gk4WHsXyCvlcMQQ8OYPZ/j8T4NSW4oKcTx
BLrEKTljAC1txGr6akMiR/BaBcZ8cuqQwMRINAop1ObICYeDGrpOgpbMTM4qLf3cciPBSW1NG+TSWY08
jkhtDcltxGRucbUYGbVh2tfasJ7juY7dLziEHVGYJiRzm7pwjW4CojCiaJy0L4Ypx906AduSYi5pXMEX
mNtAuXBJdF6FvWIwY01EnLZX2O1OY/AqfcW+g9EM1eqgaSVArSQvOS9JPkGyLvzunTHXWlLA+YafnXr+
q40+0y4QZ8Oe9aKhbUamEYjpJQVFK7bUxO8YAyMguaWHd61+oUTUyplKjl3vsyN9PKuI3UyDv2fFs3hI
EkRo/ICboqXJhIzjqTK7mnNN0taOhV1EgSrcy0c7YzPzQoIT0JzGioIqDd0sdlXk+ARlh4mJlP14ECbU
TtmCKJDXj/YArvvnx/NMNyAYA5wdDt0YKmyjwPjQxXOpyR9azLjN8/gVGdJwriNRhdGgg0Z7gObx0EQ3
if1BmEPLWipk1/Jgd7gtZc5p56rIztZKWRjujuo+YdTHb5s0P2kRKKi1u+UpjyYCUqTpmcft5pfOieUp
ByFEtjpOiGqDyoLYsY0J6oSMMSuLumB1SCIjTAvVNjSjRibtBCwigyayjSaoWx0i7exDRmaUvjwLYhgJ
XYtSQfFRrFUoIi2Ravw4LS2C3XNzvQwXfP63kh+FyUCn4thduU94Wmk+MxUzMRHPWmc4Pm1x4h8FXZZ/
FFdr2PP/ps/WQMMe67ZZ/AIro6+z80Ss/wsAAP//TJJvpd8IAAA=
`,
	},

	"/data/ru/colors": {
		local:   "data/ru/colors",
		size:    216,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0SOvQ7CMAyE9zx1KAMDCwgJUJEQSEisUP4JzTN8fiPOtDSDZd+XnH3UFjnZiNamvAJ7
nrQqjQsyycac1SVX3Eiq/uPSpVWd2PoS+e6CTYcOWpLdYJV6DzeCqci1LFehCzmwE3+UtxlvGi39KF30
UIGa7KNNZBvuHHUh/kPNf6F8/AYAAP//Xjb0JtgAAAA=
`,
	},

	"/data/ru/continents": {
		local:   "data/ru/continents",
		size:    133,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/3TMwQ2FIBCE4TtVvFJ5xDtcLMAGDCGiBEVb+LcjN3Lw5Gmz802GwEYRb5jIJLJYGlH8
T+V43kIlGmbWDwkk+YvTILL3sSDD66O65eTq5SZOZyxVb2HR8A4AAP//f/Vwi4UAAAA=
`,
	},

	"/data/ru/countries": {
		local:   "data/ru/countries",
		size:    3010,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/2xWWXYbtxL95164Rw6Wn86RX2QNObZJMSIlJdaP4maTMMGhm1tAbyEryb23Ct2t43wB
dQHUPCB9TmUzbibNKBXpmGJzPehBTu5SALHGjX3app9YK8JHQVW+dcRJbEbcVmmTagCFEdtUNyNKENlM
UkybZioS6CYFB6sMnQg53+YDbpgYKGWyb7AQPDVXRlAgxEFSaC4FucIUPTbgAlCABXofcJ1KTcF0bJIM
jGnn2+YTmLYnFfXjlnZtxN/Par0qWxL8qo6YgCylIWy7kZt3uO6uvqEKOINeJKY4P5gjhukFQsapdlzX
tika6X645UJPQsyUpFlwgE5rKTjpAnQrIzZZbyMDzJ82/7dQ3VJ3PKmb/4HrFVW6hRMCoIrOHqQ7LGsw
r2xL/tG2+2wiiRMuSchdh5ZUEjJPyjOHolS4fkcM4TaEGtenxIPSobXijtLTVp4qGNk7RbVSqAsnm49+
WW7bGXHfMblnpsJDU9P/HjxPEk/XFb8AcMaYK5/+zpxIZzplkL70bf0C/KScW8OYMEhf007hjkwlkoof
r9q2Vtx3GUI20oSs41e9PeRNZZujLmQ2UOsX4NxnMemK+pt8H8hoZlEcMgckNBi2U5G0RTZz62qr+QwE
y0UjW7/PFN6C5T/rindGXxkWVekx22tApL8sCLO0SEsuKCfwb906UyxYqij1JRMWHKc4PXbRomV1Oigl
wAzpG7qIkfE9Wf/n45nnPOpt2ccmw7T9Z/QNcSlVX7RBLwvblCxZ5sAcHlCLmUu4t4G5ZNUsJhK0M+TK
m6tbyMXzrnHMVU0l+XN7oWKpmkutalrz5jfYOJa1axXAZpAewKiUJ9k39j2gX/sPFib1LXPMyMADuFkm
tveYUSUjkol9jpkB+QSdYUv1GRcHKFIE5CMe7ENONpd6i56bHwTZ4h5/AL7LxcSkf1BX3ZoqJrtWsjnL
uh0xdgrZ14BOZLYQo5gTaMFOwpzlNqg+jtxGDh76od3m+1FjpC026rswTVWKProWtMVSbKHiZUr4wPgD
c+6Ny8ni/CjNewPs0XxLpR69jKQ7ibPxf8zilCx7AsHMeFSNfKJPRQBF+9dAdH8sldU+15bkpyZRGIGp
mQO+8kiWPipHBrGcCttWQw8o22HdYk8qyyAaLLcKbFtzKzG0OcdmsCXEgrm0waqb0Z5rDJnyTrL4bilP
yV6xylqcfWPCcmsRlEW2dOX/FCeikkQONcNijvHKMtGj0gNC55jX9JlLrZOTMs6TdyXAa2El6xXUlYoy
5mG50uzUv2fIhmANfwDXFeqo8X1KPNmosFAZWakoXCN6fMMIPHn2I0BPLi9KBb5Sa2JMp0bo26Stvha9
j5WgPCmf/Y+lHHmG2DULtK/es/XX/Et71hRoc/NFfSTkDgyYH4jAg2jG2vao4XlWcK8MqqDsdTfDXmzc
ZcX+9NKyxg7gL4u7/VfhYUalaD54C3w3q7+bT75Lr8rclzOA4EXeRgvlq3f1j93vTlDW5ZXoUB2/sr71
hqVXOk7qC0V9sobS5s2asv4nIvXVytF4wzc3gjwb+TcayJLLz/wT+6GU25vcH7iXn/4bAAD///NH2GXC
CwAA
`,
	},

	"/data/ru/currencies": {
		local:   "data/ru/currencies",
		size:    4636,
		modtime: 1464022103,
		compressed: `
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
`,
	},

	"/data/ru/female_first_names": {
		local:   "data/ru/female_first_names",
		size:    506,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0SQXU7DQAyE3zl1k+cKUCECtZS/A6AlP7BKWrjC+EadsXcbVars8fjzZHGPwTbI+Ocv
44x0Q2nBhNkaJAojxy5yZo21+ke2u5DKwpn6H5tC2HFTDMlXaQrQjsosu84K8+DmYnskfmNb6R16jKTO
YetYJO5NHtfNXRypZcYvUT0draQnnuDe8+rf05Sc6sS9J2jWjAcdlxRNrsaD3eKbSXrbRjPixOki20vs
BKA02jny0sh0P/xP3ioVlvi2Y914vb6zNxxjUPnmJD31TEgf0uT5hzjwzrJE/SgfvVTQJ4tTELyxVmdj
9qWXIuASAAD//6+InKH6AQAA
`,
	},

	"/data/ru/female_last_names": {
		local:   "data/ru/female_last_names",
		size:    1726,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/2xUW27bQAz876kdJQ0KJED8SmO4blynBpp+VZIlWNHzCtwbdch9kHb7YXjFJbnDmdml
DZWU00Aj/3+iJeXuhhrq3CPVPrSn2mVuFjMO1FODz1Syw+dn9PigLoaObkEVjbZopBa/lPFHPksE0ikj
TXF3ju4Vquu4O5fEFtD+E8pDdij+hnCBXxUzV9hr6MR4dKqDy7Cc7OwHBHogT4E9dkvFvMDihLam5Dsa
+7ny2Hgt/GRyZgg9+8OVYy5o4+fW3dIZLWt3f9kWXfA/RP4Eh3sEQMl5w16PtndKySvOSGjn2GGk0FOn
PrK0+E3KNaZG2Um7vKAEeireHT4yLFPGUmCBDOTkLLQZBGf1ETUKRSABGnKwOCvIjRgtpK8ESX5ZsBMk
rTLDi0YDu8DsOQZ+YNFbjV6l4WQaotqM98LeBfd1VLkGqKdkHvaAiNkZd4KPUU94BtUGMhtw5BzlfS1d
HpSrlbsVM6fdwZ6whfY3gmOR7o8AbiNXb3LpxnAXkwAznNFdu+jOmm8rwCsh2gBk7RDE3KWlggsZSij+
iqzW2FjaFurJ36JgIQyPeq+Zzlb9louXzAvzLpapXIa5uOUlnb8YAbCaHgvR7+LqNv8+OIVcmjiglEza
oVaidsJ95+3nuc2SmIBZ6rs2s4/YBtdmuOLRu0+B7INuGEqzFpCpVTsVMkwaeCWU3cojljgcYuPINQdb
fw4EYmN8pPsCRA9KzlY84kcMoXe0hwCQOjzgM+vnozeWd9qSHze1gAiEgnu5RE+cXAv1mdoux200FPx0
X6wj1igYrKvEkiwxKxrevCuPDRENmgN6pOpvAAAA///AVLtIvgYAAA==
`,
	},

	"/data/ru/female_patronymics": {
		local:   "data/ru/female_patronymics",
		size:    237,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/1SOQQ7DIAwE7301WO2NQ5/iViBVUJIvrH+UFQmRERc8DOvFGw0F1SIUHdkCNnx50wef
TlB4JrJgL/ywe2sGeA/dhE53Ho3KsVlCXr4TfFCYLN4OtJX7k9/+HC3V4oIHkit76Soc/3fqEQAA///8
O6Kc7QAAAA==
`,
	},

	"/data/ru/genders": {
		local:   "data/ru/genders",
		size:    30,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/7ow52LzhW0XGy/surDvwk6uC9MubL2wF8zdAeQCAgAA//9kZ2SfHgAAAA==
`,
	},

	"/data/ru/languages": {
		local:   "data/ru/languages",
		size:    419,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/1yQTW6DMBCF95zagGCF1J8VolRUtF3jNqCYgMkV3twob5CQTVZ4vpl5fDYaOMlgMUuK
GxzmBDWPdyIf0Cs8Lli0CPCXH48Rjn0xXPgLvQ6bGMkk55rFIlVovWHa057wl+SShnKgwhYrDPg/O33A
SsHfe9gEnwxdMUl5dBsVIIx8f3ZJLyWtrgH3nFpUVPIELUM5RfNC79VTysSxPRNWHk/q71JxJlMTrAF/
c3XU0Hp/YV15iS/QKjw/fcccw9Jx8YCPAAAA//9qph2bowEAAA==
`,
	},

	"/data/ru/male_first_names": {
		local:   "data/ru/male_first_names",
		size:    430,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0xQy0rEQBC8z1ebBW8L6rqCiK8FzxI2GZjNZOIvVP+RVd076CWpx/SjGveoyFhswIiG
2W7SPynjItow2g4b5dIFPgz3ARthsSHhwBasFJVDyhalVx34W9RGI0gIxz+vkswkKw36j2zeNNfFSwga
NF2pD5litoSjF/gaR2+y684TznzorugzKZMKeLntBavt7S7hhSGbDb5a47fo4SvJ4klWkWK3FBgr4U2R
IkknmyfhnHc/4pTwQX4mrIJZeyV88t0aa5z8YFMc88RWGT/hfBHOca7va5jfAAAA//90NUW2rgEAAA==
`,
	},

	"/data/ru/male_last_names": {
		local:   "data/ru/male_last_names",
		size:    1531,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/2RU224aSxB891fjxQcdyZbCzTEi2AQHKc5TloUVeGF3f6Hnj1JdPT0DyQNaZqYv1VU1
IwuppJRWOqnuZCpluJeTnMOT1LqxljoUYWCnG7nICYsYvMLiP+R+ytk2tmEiB+lyeCcNfvH0NxcVlrFy
J72djFHxgKzaTsYMaQDkn40yRjLtGzZ3+B0saob9k+y1v6PfhAJ/+jzfBssLUMblGieV45vgs0exFPyK
coa/tHJzzl+wDzeerZ1zp6GNLZZhKEcUqsPouhiy8W2NG/YNT4CD83fsX1DswYd+Q92IbIxdRQVtfLKt
ioRf7xxiMiTsPfsFwdDGsa3wt8Co8XSqIHRcnJcqWYKM+hdDiBRSTlg8x+fokBa0CANn7Fxeh67YufHJ
VfWTL1eRsSOXd/Id30sm/Y11+lQHaWmKFzUbCK1NrhoYvkTtVUjqck6WwsCdV30Ggwmd+qbTcydzzuxH
J2IWhvRePGlz1SUEvGffSTQ4oTVGxDvvQxcvSeR0gLrnWws8ZM8sCfFA/hIcFQJbmK3Kw2qKNmfaV0Q0
yXUstnMj/aIcOzLX+UVToho3SkkjpCv+QcUPocAEWuqaqp/aE8hS9oRqXN2n0983fkdn2yAM7j2zdiJW
ZPSsQI21IgoDSJU/JYP8cizg6/aGIfOMt11HFQDdIyagvXEr7Ag5jjQjHUO+HJGd1ssZg7rVWG0QruJ+
RksDwaOPvqTKNgg3PlAUpIaRvYyD7L6tmUI9MtXXxGUk3QgdeZMtX7PO+XgF6mMe80f4P6s6R2ibHUEj
qVSqjT0yt/5orT9KAqZR8ScAAP//DObsTPsFAAA=
`,
	},

	"/data/ru/male_patronymics": {
		local:   "data/ru/male_patronymics",
		size:    1165,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/2yTz07zMBDE799Tt/1AvVQCSitVUP4UcQVKEkshTtJXWL8Rs0NwdhvUQ62fx+vx7Eau
JUqQJs3lKJ1UaSa9FFKn5T+zFfDLsJNjWkAV08pjPeyV1HWmIhQNLqrc0TSDLkg76m6wLKWiJQtxvIad
jNZQqEe9qLbaYQOVzUWAOFx742uARn3alwMCHc+VUa1D3wI49bCB6roqzFa6Skv4mGzcqmuNkgWDxeqk
dFCfDdT752wYD54zVt2wqotoQ7cLlshHt1iUzG1lYelT2KYLPjaDHRauJbvhSBiBRjYse00bNbL8jrmp
42ghRiTNmfekj/eaB19ZGYgrTJR7HGzY2NZB1CrsPO7xdwntl739Qduvw4yNrHzkUJWj6onNixYEZmpA
HD4L4/4Zy09O4cLD1sd44PSHM8QxMP4PPxH9NWYHfkAnX/UFhmYceYMYEnpi6r5iWfnOK4r8MOYW1oQn
/DJ8w7f8nxeZBr37Qf1gwgVD/UV62+Qd3wEAAP//W2hQm40EAAA=
`,
	},

	"/data/ru/months": {
		local:   "data/ru/months",
		size:    148,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/0yMMQ4CIRQFe069rq1RCgsLNTGxFzEUYtArzLuR/0Oi2zHD+8ONRiZp0iZwpZA1kXg5
HrufA5GP2fKzPAMHbWnO/ugf0UIPrbXyk4vNm2btuI/2mbrEE+8/7G1cLTvwGwAA///p67TZlAAAAA==
`,
	},

	"/data/ru/phones_format": {
		local:   "data/ru/phones_format",
		size:    26,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/1LWVVaGYSDiQrAAAQAA//+15odeGgAAAA==
`,
	},

	"/data/ru/states": {
		local:   "data/ru/states",
		size:    2756,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/4xWS1IbQQzdcwovkwV3NEMIVEHFQEiFMmBiUxSrBGN7YPB8uIL6CjlJ3nvdNqSYkVlA
eVqv9XlSS7JjW4RDm1seBlv4KENmU3veshObhgNbWhP6OGrCjgQ1BGEXRwN8FrxyhlPeXsvPrbY5MAc4
TJghJA/Q17eFFRA321RvJdTzcI0prYIvy/UBZNOwD0FuM166g8o+5LD2FpPjYvpsrLJiyy6jXnvu2W8r
+d3IeAM9Qo4B2LNp75Pd2xLOytHPPKcpWqmBGPRsBFMpjt7f/lmPFDHOqGYCUEZbr+FPEMFMmidhl7+g
9QamK9GWCLkFlrGnMBjXvtX6GS+BvRiRKHibGdxZxsiWlDFVPyO99kydwB21YEhvRS5D1iaVmR2mBq4q
Hg+D7LchrkVc1SZRNioV03vpWJHN8NnYC/7aI7hdFxGy2CIfQVrSQ5JMwiLjOyIamWxe5eEIcsj2lMk5
CyiZbEfyRz9WDLLUhTtRIc5VaqCxG8fn4+g5VYUtRFjhRHBKqmgPz40cdFuMyEZYFyXPa8Tx2Ik6Z3mq
UGbdGL5SvqpORHzt7AWF2oXvv9B4GN1eDeFzpfLz/BqSzw2IJmW7SYXcimJlgU2vGoZu9V3A34/FfgHU
CyL72iG/VP9lvXjesAFS5sV+qagqV88VvHmE75ur/ErN9qM4eldA4NX7yMnISGOgtoeYGRdXuixcS88T
/3djNnL5K01MF4Mu8OSyPWbTXrfjDoQ6k28pTbv44hxUnHyLDQyN2chx4DE0kecPrp5JtOfIvTfI6ep1
7Un4xo7g+Hij+4MNvexO02HA2nR03bPKIS01yLpQ31evAftB+ALFCxGumS038Fx6qIlaa0++2jJWjxeN
N9ZEje6abWveZtKQyn61v0xT7aVhymE400GtBalude6HJnCeNgwhjm307i7Wm5x3MXeRgTldupKrsUVp
JuverUo7C4fb6lG1Xvib9QWguFH9iT0w7j5cOLP/QCBXTankBthq6l8AAAD//511iVTECgAA
`,
	},

	"/data/ru/streets": {
		local:   "data/ru/streets",
		size:    694,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/1ySy27yQAyF9zx1Avr1VwIJ9aJWCqUpm7YrFC4DCSnhFY7fqMfOTD3qIkG2D/bn4+AD
AVeZSoEGvSxwRSPLCV4xoOdzZPmUkl+4aNpVG4Z7BP69ZMlSK4Zlqq+tHqUNm/2Ga4Y7Sou8/wt/mJD/
TtEqmHVlrtW8ygLZQhLVhhpS55XMcPBBb7gQb4mdFLJwzIqKb4Z8jyt56ZFoByOJ82rj3Ju2QefKjU3V
YvJDpkzdcnt6BtlGlW1kdBNsffWawwaTFt7/PhqUrbOla62t27nJyjs4RWUb/eP7nFIPrJeEG5zkPceq
ZSZ3HJs2NsZeL6tGONFztF4b68cxmPDkPfWqOuTv9yJzdOOpeHUcZa6XOefeVfoBGkCE+GTixieMZj1R
lnh/AgAA//8jPgMXtgIAAA==
`,
	},

	"/data/ru/weekdays": {
		local:   "data/ru/weekdays",
		size:    117,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/zSNOwoDMQxEe53aSSBVIOALGFdpsx/B4sU+w9ONdlwsxtI8zQhRGHScXf+Mj/RBM3I8
GZFurJIz8zd+uLxVNbEZJb5CxeI93RovFr2hoTBLPGhzW93puuB2BQAA///LVlOtdQAAAA==
`,
	},

	"/data/ru/weekdays_short": {
		local:   "data/ru/weekdays_short",
		size:    35,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/7ow/8JerguTLjZxXVh4sYHrwnIQaz6Ye2EjSKKRCxAAAP//vnIh2CMAAAA=
`,
	},

	"/data/ru/zips_format": {
		local:   "data/ru/zips_format",
		size:    8,
		modtime: 1464022103,
		compressed: `
H4sIAAAJbogA/1KGAC5AAAAA//+H3Sc9CAAAAA==
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},

	"/data": {
		isDir: true,
		local: "/data",
	},

	"/data/en": {
		isDir: true,
		local: "/data/en",
	},

	"/data/ru": {
		isDir: true,
		local: "/data/ru",
	},
}
