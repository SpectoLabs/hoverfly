.. _simulation_schema:

Simulation schema
=================

This is the JSON schema for v2 Hoverfly simulations.

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
      "field-matchers": {
        "properties": {
          "exactMatch": {
            "type": "string"
          },
          "globMatch": {
            "type": "string"
          },
          "jsonMatch": {
            "type": "string"
          },
          "regexMatch": {
            "type": "string"
          },
          "xpathMatch": {
            "type": "string"
          }
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
        "required": [
          "schemaVersion"
        ],
        "type": "object"
      },
      "request": {
        "properties": {
          "body": {
            "$ref": "#/definitions/field-matchers"
          },
          "destination": {
            "$ref": "#/definitions/field-matchers"
          },
          "headers": {
            "$ref": "#/definitions/headers"
          },
          "path": {
            "$ref": "#/definitions/field-matchers"
          },
          "query": {
            "$ref": "#/definitions/field-matchers"
          },
          "scheme": {
            "$ref": "#/definitions/field-matchers"
          }
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
        "required": [
          "request",
          "response"
        ],
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
          "status": {
            "type": "integer"
          },
          "templated": {
            "type": "boolean"
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
    "required": [
      "data",
      "meta"
    ],
    "type": "object"
  }