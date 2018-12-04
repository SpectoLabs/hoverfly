package v2

var requestResponsePairDefinition = map[string]interface{}{
	"type": "object",
	"required": []string{
		"request",
		"response",
	},
	"properties": map[string]interface{}{
		"request": map[string]interface{}{
			"$ref": "#/definitions/request",
		},
		"response": map[string]interface{}{
			"$ref": "#/definitions/response",
		},
	},
}

var requestV1Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"scheme": map[string]interface{}{
			"type": "string",
		},
		"destination": map[string]interface{}{
			"type": "string",
		},
		"path": map[string]interface{}{
			"type": "string",
		},
		"query": map[string]interface{}{
			"type": "string",
		},
		"body": map[string]interface{}{
			"type": "string",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
	},
}

var requestV3Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"scheme": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"destination": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"path": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"query": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"body": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
	},
}

var responseDefinitionV3 = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"body": map[string]interface{}{
			"type": "string",
		},
		"encodedBody": map[string]interface{}{
			"type": "boolean",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
		"status": map[string]interface{}{
			"type": "integer",
		},
		"templated": map[string]interface{}{
			"type": "boolean",
		},
	},
}

var responseDefinitionV1 = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"body": map[string]interface{}{
			"type": "string",
		},
		"encodedBody": map[string]interface{}{
			"type": "boolean",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
		"status": map[string]interface{}{
			"type": "integer",
		},
	},
}

var requestFieldMatchersV3Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"exactMatch": map[string]interface{}{
			"type": "string",
		},
		"globMatch": map[string]interface{}{
			"type": "string",
		},
		"regexMatch": map[string]interface{}{
			"type": "string",
		},
		"xpathMatch": map[string]interface{}{
			"type": "string",
		},
		"jsonMatch": map[string]interface{}{
			"type": "string",
		},
	},
}

var headersDefinition = map[string]interface{}{
	"type": "object",
	"additionalProperties": map[string]interface{}{
		"type": "array",
		"items": map[string]interface{}{
			"type": "string",
		},
	},
}

var delaysDefinition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"urlPattern": map[string]interface{}{
			"type": "string",
		},
		"httpMethod": map[string]interface{}{
			"type": "string",
		},
		"delay": map[string]interface{}{
			"type": "integer",
		},
	},
}

var delaysLogNormalDefinition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"urlPattern": map[string]interface{}{
			"type": "string",
		},
		"httpMethod": map[string]interface{}{
			"type": "string",
		},
		"min": map[string]interface{}{
			"type": "integer",
		},
		"max": map[string]interface{}{
			"type": "integer",
		},
		"mean": map[string]interface{}{
			"type": "integer",
		},
		"median": map[string]interface{}{
			"type": "integer",
		},
	},
}

var metaDefinition = map[string]interface{}{
	"type": "object",
	"required": []string{
		"schemaVersion",
	},
	"properties": map[string]interface{}{
		"schemaVersion": map[string]interface{}{
			"type": "string",
		},
		"hoverflyVersion": map[string]interface{}{
			"type": "string",
		},
		"timeExported": map[string]interface{}{
			"type": "string",
		},
	},
}

var SimulationViewV3Schema = map[string]interface{}{
	"description": "Hoverfly simulation schema",
	"type":        "object",
	"required": []string{
		"data", "meta",
	},
	"additionalProperties": false,
	"properties": map[string]interface{}{
		"data": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pairs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/request-response-pair",
					},
				},
				"globalActions": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"delays": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/delay",
							},
						},
					},
				},
			},
		},
		"meta": map[string]interface{}{
			"$ref": "#/definitions/meta",
		},
	},
	"definitions": map[string]interface{}{
		"request-response-pair": requestResponsePairDefinition,
		"request":               requestV3Definition,
		"response":              responseDefinitionV3,
		"field-matchers":        requestFieldMatchersV3Definition,
		"headers":               headersDefinition,
		"delay":                 delaysDefinition,
		"meta":                  metaDefinition,
	},
}

var SimulationViewV2Schema = map[string]interface{}{
	"description": "Hoverfly simulation schema",
	"type":        "object",
	"required": []string{
		"data", "meta",
	},
	"additionalProperties": false,
	"properties": map[string]interface{}{
		"data": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pairs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/request-response-pair",
					},
				},
				"globalActions": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"delays": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/delay",
							},
						},
					},
				},
			},
		},
		"meta": map[string]interface{}{
			"$ref": "#/definitions/meta",
		},
	},
	"definitions": map[string]interface{}{
		"request-response-pair": requestResponsePairDefinition,
		"request":               requestV3Definition,
		"response":              responseDefinitionV1,
		"field-matchers":        requestFieldMatchersV3Definition,
		"headers":               headersDefinition,
		"delay":                 delaysDefinition,
		"meta":                  metaDefinition,
	},
}

var SimulationViewV1Schema = map[string]interface{}{
	"description": "Hoverfly simulation schema",
	"type":        "object",
	"required": []string{
		"data", "meta",
	},
	"additionalProperties": false,
	"properties": map[string]interface{}{
		"data": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pairs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/request-response-pair",
					},
				},
				"globalActions": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"delays": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/delay",
							},
						},
					},
				},
			},
		},
		"meta": map[string]interface{}{
			"$ref": "#/definitions/meta",
		},
	},
	"definitions": map[string]interface{}{
		"request-response-pair": requestResponsePairDefinition,
		"request":               requestV1Definition,
		"response":              responseDefinitionV1,
		"headers":               headersDefinition,
		"delay":                 delaysDefinition,
		"meta":                  metaDefinition,
	},
}

// V4 Schema

var SimulationViewV4Schema = map[string]interface{}{
	"description": "Hoverfly simulation schema",
	"type":        "object",
	"required": []string{
		"data", "meta",
	},
	"additionalProperties": false,
	"properties": map[string]interface{}{
		"data": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pairs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/request-response-pair",
					},
				},
				"globalActions": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"delays": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/delay",
							},
						},
					},
				},
			},
		},
		"meta": map[string]interface{}{
			"$ref": "#/definitions/meta",
		},
	},
	"definitions": map[string]interface{}{
		"request-response-pair": requestResponsePairDefinition,
		"request":               requestV4Definition,
		"response":              responseDefinitionV4,
		"field-matchers":        requestFieldMatchersV3Definition,
		"headers":               headersDefinition,
		"delay":                 delaysDefinition,
		"meta":                  metaDefinition,
	},
}

var requestV4Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"scheme": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"destination": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"path": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"query": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"body": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
		"requiresState": map[string]interface{}{
			"type": "object",
			"patternProperties": map[string]interface{}{
				".{1,}": map[string]interface{}{"type": "string"},
			},
		},
	},
}

var responseDefinitionV4 = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"body": map[string]interface{}{
			"type": "string",
		},
		"encodedBody": map[string]interface{}{
			"type": "boolean",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
		"status": map[string]interface{}{
			"type": "integer",
		},
		"templated": map[string]interface{}{
			"type": "boolean",
		},
		"removesState": map[string]interface{}{
			"type": "array",
		},
		"transitionsState": map[string]interface{}{
			"type": "object",
			"patternProperties": map[string]interface{}{
				".{1,}": map[string]interface{}{"type": "string"},
			},
		},
	},
}

// V5 Schema

var SimulationViewV5Schema = map[string]interface{}{
	"description": "Hoverfly simulation schema",
	"type":        "object",
	"required": []string{
		"data", "meta",
	},
	"additionalProperties": false,
	"properties": map[string]interface{}{
		"data": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pairs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/request-response-pair",
					},
				},
				"globalActions": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"delays": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/delay",
							},
						},
						"delaysLogNormal": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/delay-log-normal",
							},
						},
					},
				},
			},
		},
		"meta": map[string]interface{}{
			"$ref": "#/definitions/meta",
		},
	},
	"definitions": map[string]interface{}{
		"request-response-pair": requestResponsePairDefinition,
		"request":               requestV5Definition,
		"response":              responseDefinitionV4,
		"field-matchers":        requestFieldMatchersV5Definition,
		"headers":               headersDefinition,
		"request-headers":       v5MatchersMapDefinition,
		"request-queries":       v5MatchersMapDefinition,
		"delay":                 delaysDefinition,
		"delay-log-normal":      delaysLogNormalDefinition,
		"meta":                  metaDefinition,
	},
}

var requestV5Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"scheme": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"$ref": "#/definitions/field-matchers",
			},
		},
		"destination": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"$ref": "#/definitions/field-matchers",
			},
		},
		"path": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"$ref": "#/definitions/field-matchers",
			},
		},
		"query": map[string]interface{}{
			"$ref": "#/definitions/request-queries",
		},
		"body": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"$ref": "#/definitions/field-matchers",
			},
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/request-headers",
		},
		"requiresState": map[string]interface{}{
			"type": "object",
			"patternProperties": map[string]interface{}{
				".{1,}": map[string]interface{}{"type": "string"},
			},
		},
	},
}

var requestFieldMatchersV5Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"matcher": map[string]interface{}{
			"type": "string",
		},
		"value": map[string]interface{}{},
	},
}

var v5MatchersMapDefinition = map[string]interface{}{
	"type": "object",
	"additionalProperties": map[string]interface{}{
		"type": "array",
		"items": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
	},
}
