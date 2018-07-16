package testdata

var TemplatingRequest = `{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Body_jsonpath"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Body \"jsonpath\" \"$.test\" }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Body_xpath"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Body \"xpath\" \"/root/text\" }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Method"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Method }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.Scheme"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Scheme }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/one/two/three/Request.Path"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Path }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/one/two/three/Request.Path0"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.Path.[0] }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.QueryParam"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.QueryParam }}",
					"encodedBody": false,
					"templated": true
				}
			},
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/Request.QueryParam.query"
						}
					]
				},
				"response": {
					"status": 200,
					"body": "{{ Request.QueryParam.query }}",
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
