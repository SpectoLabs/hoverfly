package testdata

var TemplatingHelpers = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomString"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomString }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomStringLength10"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomStringLength 10 }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomBoolean"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomBoolean }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomInteger"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomInteger }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomIntegerRange1-10"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomIntegerRange 1 10 }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomFloat"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomFloat }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomFloatRange1-10"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomFloatRange 1.0 10.0 }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomEmail"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomEmail }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomIPv4"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomIPv4 }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomIPv6"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomIPv6 }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/randomuuid"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ randomUuid }}",
					"encodedBody": false,
					"templated": true
				}
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v5"
	}
}`
