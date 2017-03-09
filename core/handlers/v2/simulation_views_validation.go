package v2

var SimulationViewV2JsonSchema = `{
	"description": "Hoverfly simulation schema",
	"type": "object",
	"required": [
		"data", "meta"
	],
	"additionalProperties": false,
	"properties": {
		"data": {
			"type": "object",
			"properties": {
				"pairs": {
					"type": "array",
					"items": {
						"$ref": "#/definitions/request-response-pair"
					}
				},
				"globalActions": {
					"type": "object",
					"properties": {
						"delays": {
							"type": "array",
							"items": {
								"$ref": "#/definitions/delay"
							}
						}
					}
				}
			}
		},
		"meta": {
			"type": "object",
			"required": [
				"schemaVersion"
			],
			"properties": {
				"schemaVersion": {
					"type": "string"
				},
				"hoverflyVersion": {
					"type": "string"
				},
				"timeExported": {
					"type": "string"
				}
			}
		}
	},
	"definitions": {
		"request-response-pair": {
			"type": "object",
			"required": [
				"request",
				"response"
			],
			"properties": {
				"request": {
					"$ref": "#/definitions/request"
				},
				"response": {
					"$ref": "#/definitions/response"
				}
			}
		},
		"request": {
			"type": "object",
			"properties": {
				"scheme": {
					"$ref": "#/definitions/field-matchers"
				},
				"destination": {
					"$ref": "#/definitions/field-matchers"
				},
				"path": {
					"$ref": "#/definitions/field-matchers"
				},
				"query": {
					"$ref": "#/definitions/field-matchers"
				},
				"body": {
					"$ref": "#/definitions/field-matchers"
				},
				"headers": {
					"$ref": "#/definitions/headers"
				}
			}
		},
		"response": {
			"type": "object",
			"properties": {
				"body": {
					"type": "string"
				},
				"encodedBody": {
					"type": "boolean"
				},
				"headers": {
					"$ref": "#/definitions/headers"
				},
				"status": {
					"type": "integer"
				}
			}
		},
		"field-matchers": {
			"type": "object",
			"properties": {
				"exactMatch": {
					"type": "string"
				},
				"globMatch": {
					"type": "string"
				},
				"regexMatch": {
					"type": "string"
				},
				"xpathMatch": {
					"type": "string"
				},
				"jsonMatch": {
					"type": "string"
				}
			}
		},
		"headers": {
			"type": "object",
			"additionalProperties": {
				"type": "array",
				"items": {
					"type": "string"
				}
			}
		},
		"delay": {
			"type": "object",
			"properties": {
				"urlPattern": {
					"type": "string"
				},
				"httpMethod": {
					"type": "string"
				},
				"delay": {
					"type": "integer"
				}
			}
		}
	}
}`
