package matchers

import (
    "strings"

    "github.com/SpectoLabs/hoverfly/core/util"
)

// JWTJsonPath is the matcher key for JSONPath matching over a JWT's payload/header
var JWTJsonPath = "jwtjsonpath"

// JwtJsonPathMatch evaluates a JSONPath against the decoded JWT (header+payload).
// Returns true when the path resolves to at least one value.
func JwtJsonPathMatch(match interface{}, toMatch string) bool {
    path, ok := match.(string)
    if !ok || path == "" {
        return false
    }

    composite, err := util.ParseJWTComposite(toMatch)
    if err != nil {
        return false
    }

    norm := normalizeJWTJsonPath(path)
    norm = util.PrepareJsonPathQuery(norm)

    out, err := util.JsonPathExecution(norm, composite)
    if err != nil || out == norm {
        return false
    }
    return true
}

// JwtJsonPathMatchValueGenerator extracts the JSONPath result as a string to feed into chained matchers.
func JwtJsonPathMatchValueGenerator(match interface{}, toMatch string) string {
    composite, err := util.ParseJWTComposite(toMatch)
    if err != nil {
        return ""
    }
    // Normalize path to default to payload, then reuse existing generator
    if p, ok := match.(string); ok {
        norm := normalizeJWTJsonPath(p)
        return JsonPathMatcherValueGenerator(norm, composite)
    }
    return ""
}

// normalizeJWTJsonPath allows shorthand paths like "$.user_name" to address payload
// by default. If the path already targets $.payload or $.header, it is left as-is.
func normalizeJWTJsonPath(path string) string {
    p := strings.TrimSpace(path)
    lower := strings.ToLower(p)
    if strings.HasPrefix(lower, "$.payload.") || strings.HasPrefix(lower, "$.header.") {
        return p
    }
    if strings.HasPrefix(p, "$.") {
        return "$.payload" + p[1:]
    }
    // if someone passed without leading '$', make it refer to payload
    if strings.HasPrefix(p, ".") {
        return "$.payload" + p
    }
    return p
}
