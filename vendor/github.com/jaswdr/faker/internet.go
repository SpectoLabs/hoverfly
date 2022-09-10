package faker

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	freeEmailDomain = []string{"gmail.com", "yahoo.com", "hotmail.com"}

	tld = []string{"com", "com", "com", "com", "com", "com", "biz", "info", "net", "org"}

	userFormats = []string{"{{lastName}}.{{firstName}}",
		"{{firstName}}.{{lastName}}",
		"{{firstName}}",
		"{{lastName}}"}

	emailFormats = []string{"{{user}}@{{domain}}", "{{user}}@{{freeEmailDomain}}"}

	urlFormats = []string{"http://www.{{domain}}/",
		"http://{{domain}}/",
		"http://www.{{domain}}/{{slug}}",
		"http://www.{{domain}}/{{slug}}",
		"https://www.{{domain}}/{{slug}}",
		"http://www.{{domain}}/{{slug}}.html",
		"http://{{domain}}/{{slug}}",
		"http://{{domain}}/{{slug}}",
		"http://{{domain}}/{{slug}}.html",
		"https://{{domain}}/{{slug}}.html",
	}

	statusCodes        = []string{"100", "101", "102", "200", "201", "202", "203", "204", "205", "206", "207", "208", "226", "300", "301", "302", "303", "304", "305", "306", "307", "308", "400", "401", "402", "403", "404", "405", "406", "407", "408", "409", "410", "411", "412", "413", "414", "415", "416", "417", "418", "420", "422", "423", "424", "425", "426", "428", "429", "431", "444", "449", "450", "451", "499", "500", "501", "502", "503", "504", "505", "506", "507", "508", "509", "510", "511", "598", "599"}
	statusCodeMessages = []string{"Continue", "Switching Protocols", "Processing (WebDAV)", "OK", "Created", "Accepted", "Non-Authoritative Information", "No Content", "Reset Content", "Partial Content", "Multi-Status (WebDAV)", "Already Reported (WebDAV)", "IM Used", "Multiple Choices", "Moved Permanently", "Found", "See Other", "Not Modified", "Use Proxy", "(Unused)", "Temporary Redirect", "Permanent Redirect (experimental)", "Bad Request", "Unauthorized", "Payment Required", "Forbidden", "Not Found", "Method Not Allowed", "Not Acceptable", "Proxy Authentication Required", "Request Timeout", "Conflict", "Gone", "Length Required", "Precondition Failed", "Request Entity Too Large", "Request-URI Too Long", "Unsupported Media Type", "Requested Range Not Satisfiable", "Expectation Failed", "I'm a teapot (RFC 2324)", "Enhance Your Calm (Twitter)", "Unprocessable Entity (WebDAV)", "Locked (WebDAV)", "Failed Dependency (WebDAV)", "Reserved for WebDAV", "Upgrade Required", "Precondition Required", "Too Many Requests", "Request Header Fields Too Large", "No Response (Nginx)", "Retry With (Microsoft)", "Blocked by Windows Parental Controls (Microsoft)", "Unavailable For Legal Reasons", "Client Closed Request (Nginx)", "Internal Server Error", "Not Implemented", "Bad Gateway", "Service Unavailable", "Gateway Timeout", "HTTP Version Not Supported", "Variant Also Negotiates (Experimental)", "Insufficient Storage (WebDAV)", "Loop Detected (WebDAV)", "Bandwidth Limit Exceeded (Apache)", "Not Extended", "Network Authentication Required", "Network read timeout error", "Network connect timeout error"}
)

// Internet is a faker struct for Internet
type Internet struct {
	Faker *Faker
}

func (i Internet) transformIntoValidEmailName(name string) string {
	name = strings.ToLower(name)
	onlyValidCharacters := regexp.MustCompile(`[^a-z0-9._%+\-]+`)
	name = onlyValidCharacters.ReplaceAllString(name, "_")
	return name
}

// User returns a fake user for Internet
func (i Internet) User() string {
	user := i.Faker.RandomStringElement(userFormats)

	p := i.Faker.Person()

	// {{firstName}}
	user = strings.Replace(user, "{{firstName}}", i.transformIntoValidEmailName(p.FirstName()), 1)

	// {{lastName}}
	user = strings.Replace(user, "{{lastName}}", i.transformIntoValidEmailName(p.LastName()), 1)

	return user
}

// Password returns a fake password for Internet
func (i Internet) Password() string {
	pattern := strings.Repeat("*", i.Faker.IntBetween(6, 16))
	return i.Faker.Asciify(pattern)
}

// Domain returns a fake domain for Internet
func (i Internet) Domain() string {
	domain := strings.ToLower(i.Faker.Lexify("???"))
	return strings.Join([]string{domain, i.TLD()}, ".")
}

// FreeEmailDomain returns a fake free email domain for Internet
func (i Internet) FreeEmailDomain() string {
	return i.Faker.RandomStringElement(freeEmailDomain)
}

// SafeEmailDomain returns a fake safe email domain for Internet
func (i Internet) SafeEmailDomain() string {
	return "example.org"
}

// Email returns a fake email address for Internet
func (i Internet) Email() string {
	email := i.Faker.RandomStringElement(emailFormats)

	// {{user}}
	email = strings.Replace(email, "{{user}}", i.transformIntoValidEmailName(i.User()), 1)

	// {{domain}}
	email = strings.Replace(email, "{{domain}}", i.Domain(), 1)

	// {{freeEmailDomain}}
	email = strings.Replace(email, "{{freeEmailDomain}}", i.FreeEmailDomain(), 1)

	return email
}

// FreeEmail returns a fake free email address for Internet
func (i Internet) FreeEmail() string {
	domain := i.Faker.RandomStringElement(freeEmailDomain)

	return strings.Join([]string{i.transformIntoValidEmailName(i.User()), domain}, "@")
}

// SafeEmail returns a fake safe email address for Internet
func (i Internet) SafeEmail() string {
	return strings.Join([]string{i.transformIntoValidEmailName(i.User()), i.SafeEmailDomain()}, "@")
}

// CompanyEmail returns a fake company email address for Internet
func (i Internet) CompanyEmail() string {
	c := i.Faker.Company()

	companyName := i.transformIntoValidEmailName(c.Name())

	domain := strings.Join([]string{companyName, i.Domain()}, ".")

	return strings.Join([]string{i.transformIntoValidEmailName(i.User()), domain}, "@")
}

// TLD returns a fake tld for Internet
func (i Internet) TLD() string {
	return i.Faker.RandomStringElement(tld)
}

// Slug returns a fake slug for Internet
func (i Internet) Slug() string {
	slug := strings.Repeat("?", i.Faker.IntBetween(1, 5)) +
		"-" +
		strings.Repeat("?", i.Faker.IntBetween(1, 6))

	return strings.ToLower(i.Faker.Lexify(slug))
}

// URL returns a fake url for Internet
func (i Internet) URL() string {
	url := i.Faker.RandomStringElement(urlFormats)

	// {{domain}}
	url = strings.Replace(url, "{{domain}}", i.Domain(), 1)

	// {{slug}}
	url = strings.Replace(url, "{{slug}}", i.Slug(), 1)

	return url
}

// Ipv4 returns a fake ipv4 for Internet
func (i Internet) Ipv4() string {
	ips := []string{}

	for j := 0; j < 4; j++ {
		ips = append(ips, strconv.Itoa(i.Faker.IntBetween(1, 255)))
	}

	return strings.Join(ips, ".")
}

// LocalIpv4 returns a fake local ipv4 for Internet
func (i Internet) LocalIpv4() string {
	ips := []string{i.Faker.RandomStringElement([]string{"10", "172", "192"})}

	if ips[0] == "10" {
		for j := 0; j < 3; j++ {
			ips = append(ips, strconv.Itoa(i.Faker.IntBetween(1, 255)))
		}
	}

	if ips[0] == "172" {
		ips = append(ips, strconv.Itoa(i.Faker.IntBetween(16, 31)))

		for j := 0; j < 2; j++ {
			ips = append(ips, strconv.Itoa(i.Faker.IntBetween(1, 255)))
		}
	}

	if ips[0] == "192" {
		ips = append(ips, "168")

		for j := 0; j < 2; j++ {
			ips = append(ips, strconv.Itoa(i.Faker.IntBetween(1, 255)))
		}
	}

	return strings.Join(ips, ".")
}

// Ipv6 returns a fake ipv6 for Internet
func (i Internet) Ipv6() string {
	ips := []string{}

	for j := 0; j < 8; j++ {
		block := ""
		for w := 0; w < 4; w++ {
			block = block + strconv.Itoa(i.Faker.RandomDigitNotNull())
		}

		ips = append(ips, block)
	}

	return strings.Join(ips, ":")
}

// MacAddress returns a fake mac address for Internet
func (i Internet) MacAddress() string {
	values := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}

	mac := []string{}
	for j := 0; j < 6; j++ {
		m := i.Faker.RandomStringElement(values)
		m = m + i.Faker.RandomStringElement(values)
		mac = append(mac, m)
	}

	return strings.Join(mac, ":")
}

// HTTPMethod returns a fake http method for Internet
func (i Internet) HTTPMethod() string {
	return i.Faker.RandomStringElement([]string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	})
}

// Query returns a fake query for Internet
func (i Internet) Query() string {
	lorem := i.Faker.Lorem()
	query := "?" + lorem.Word() + "=" + lorem.Word()
	for j := 0; j < i.Faker.IntBetween(1, 3); j++ {
		if i.Faker.Boolean().Bool() {
			query += "&" + lorem.Word() + "=" + lorem.Word()
		} else {
			query += "&" + lorem.Word() + "=" + strconv.Itoa(i.Faker.RandomDigitNotNull())
		}
	}

	return query
}

// StatusCode returns a fake status code for Internet
func (i Internet) StatusCode() int {
	statusCode, _ := strconv.Atoi(i.Faker.RandomStringElement(statusCodes))
	return statusCode
}

// StatusCodeMessage returns a fake status code message for Internet
func (i Internet) StatusCodeMessage() string {
	return i.Faker.RandomStringElement(statusCodeMessages)
}

// StatusCodeWithMessage returns a fake status code with message for Internet
func (i Internet) StatusCodeWithMessage() string {
	index := i.Faker.IntBetween(0, len(statusCodes)-1)
	return statusCodes[index] + " " + statusCodeMessages[index]
}
