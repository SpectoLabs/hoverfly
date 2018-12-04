.. _simulation_schema:

Simulation schema
=================

This is the JSON schema for v5 Hoverfly simulations.

.. code:: json

  {
    "additionalProperties": false,
    "definitions": {
      "delay": {
        "properties": {
          "delay": {
            "type": "integer"
          },
          "httpMethod": {
            "type": "string"
          },
          "urlPattern": {
            "type": "string"
          }
        },
        "type": "object"
      },
      "delay-log-normal": {
        "properties": {
          "min": {
            "type": "integer"
          },
          "max": {
            "type": "integer"
          },
          "mean": {
            "type": "integer"
          },
          "median": {
            "type": "integer"
          },
          "httpMethod": {
            "type": "string"
          },
          "urlPattern": {
            "type": "string"
          }
        },
        "type": "object"
      },
      "field-matchers": {
        "properties": {
          "matcher": {
            "type": "string"
          },
          "value": {}
        },
        "type": "object"
      },
      "headers": {
        "additionalProperties": {
          "items": {
            "type": "string"
          },
          "type": "array"
        },
        "type": "object"
      },
      "meta": {
        "properties": {
          "hoverflyVersion": {
            "type": "string"
          },
          "schemaVersion": {
            "type": "string"
          },
          "timeExported": {
            "type": "string"
          }
        },
        "required": ["schemaVersion"],
        "type": "object"
      },
      "request": {
        "properties": {
          "body": {
            "items": {
              "$ref": "#/definitions/field-matchers"
            },
            "type": "array"
          },
          "destination": {
            "items": {
              "$ref": "#/definitions/field-matchers"
            },
            "type": "array"
          },
          "headers": {
            "$ref": "#/definitions/request-headers"
          },
          "path": {
            "items": {
              "$ref": "#/definitions/field-matchers"
            },
            "type": "array"
          },
          "query": {
            "$ref": "#/definitions/request-queries"
          },
          "requiresState": {
            "patternProperties": {
              ".{1,}": {
                "type": "string"
              }
            },
            "type": "object"
          },
          "scheme": {
            "items": {
              "$ref": "#/definitions/field-matchers"
            },
            "type": "array"
          }
        },
        "type": "object"
      },
      "request-headers": {
        "additionalProperties": {
          "items": {
            "$ref": "#/definitions/field-matchers"
          },
          "type": "array"
        },
        "type": "object"
      },
      "request-queries": {
        "additionalProperties": {
          "items": {
            "$ref": "#/definitions/field-matchers"
          },
          "type": "array"
        },
        "type": "object"
      },
      "request-response-pair": {
        "properties": {
          "request": {
            "$ref": "#/definitions/request"
          },
          "response": {
            "$ref": "#/definitions/response"
          }
        },
        "required": ["request", "response"],
        "type": "object"
      },
      "response": {
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
          "removesState": {
            "type": "array"
          },
          "status": {
            "type": "integer"
          },
          "templated": {
            "type": "boolean"
          },
          "transitionsState": {
            "patternProperties": {
              ".{1,}": {
                "type": "string"
              }
            },
            "type": "object"
          }
        },
        "type": "object"
      }
    },
    "description": "Hoverfly simulation schema",
    "properties": {
      "data": {
        "properties": {
          "globalActions": {
            "properties": {
              "delays": {
                "items": {
                  "$ref": "#/definitions/delay"
                },
                "type": "array"
              },
              "delaysLogNormal": {
                "items": {
                  "$ref": "#/definitions/delay-log-normal"
                },
                "type": "array"
              }
            },
            "type": "object"
          },
          "pairs": {
            "items": {
              "$ref": "#/definitions/request-response-pair"
            },
            "type": "array"
          }
        },
        "type": "object"
      },
      "meta": {
        "$ref": "#/definitions/meta"
      }
    },
    "required": ["data", "meta"],
    "type": "object"
  }
